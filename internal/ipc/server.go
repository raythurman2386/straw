package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
)

// Server handles IPC via Unix Domain Sockets.
type Server struct {
	socketPath string
	listener   net.Listener
	clients    map[net.Conn]bool
	mu         sync.Mutex
	handlers   map[Method]func(json.RawMessage) (interface{}, error)
	shutdown   chan struct{}
}

// NewServer creates a new IPC server.
func NewServer(socketPath string) *Server {
	return &Server{
		socketPath: socketPath,
		clients:    make(map[net.Conn]bool),
		handlers:   make(map[Method]func(json.RawMessage) (interface{}, error)),
		shutdown:   make(chan struct{}),
	}
}

// Start begins listening on the socket.
func (s *Server) Start() error {
	// Cleanup existing socket
	if _, err := os.Stat(s.socketPath); err == nil {
		slog.Debug("Removing existing socket file", "path", s.socketPath)
		if err := os.Remove(s.socketPath); err != nil {
			return err
		}
	}

	l, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}

	// Restrict permissions to owner only
	if err := os.Chmod(s.socketPath, 0600); err != nil {
		l.Close()
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	s.listener = l

	slog.Info("IPC server listening", "socket", s.socketPath)

	go s.acceptLoop()
	return nil
}

func (s *Server) acceptLoop() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				slog.Error("Failed to accept connection", "error", err)
				continue
			}
		}

		slog.Debug("New client connected", "remote", conn.RemoteAddr())

		s.mu.Lock()
		s.clients[conn] = true
		s.mu.Unlock()

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		slog.Debug("Client disconnected", "remote", conn.RemoteAddr())
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()
		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			slog.Error("Failed to unmarshal request", "error", err, "payload", string(line))
			continue
		}

		slog.Debug("Received request", "method", req.Method, "id", req.ID)
		s.processRequest(conn, req)
	}
}

func (s *Server) processRequest(conn net.Conn, req Request) {
	var resp Response
	resp.ID = req.ID

	handler, ok := s.handlers[req.Method]
	if !ok {
		slog.Warn("Unknown method requested", "method", req.Method)
		resp.Error = fmt.Sprintf("method not found: %s", req.Method)
	} else {
		res, err := handler(req.Params)
		if err != nil {
			slog.Error("Handler error", "method", req.Method, "error", err)
			resp.Error = err.Error()
		} else {
			bytes, _ := json.Marshal(res)
			resp.Result = bytes
		}
	}

	respBytes, _ := json.Marshal(resp)
	if _, err := conn.Write(append(respBytes, '\n')); err != nil {
		slog.Error("Failed to write response to client", "remote", conn.RemoteAddr(), "error", err)
	}
}

// Register adds a handler for a specific method.
func (s *Server) Register(method Method, handler func(json.RawMessage) (interface{}, error)) {
	s.handlers[method] = handler
}

// Broadcast sends an event to all connected clients.
func (s *Server) Broadcast(eventType EventType, payload interface{}) {
	event := Event{
		Type: eventType,
	}
	bytes, _ := json.Marshal(payload)
	event.Payload = bytes

	msg, _ := json.Marshal(event)
	msg = append(msg, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		if _, err := conn.Write(msg); err != nil {
			slog.Debug("Failed to write to client", "remote", conn.RemoteAddr(), "error", err)
		}
	}
}

// Stop closes the listener and disconnects clients.
func (s *Server) Stop() {
	slog.Info("Stopping IPC server")
	close(s.shutdown)
	if s.listener != nil {
		s.listener.Close()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.clients {
		conn.Close()
	}
}
