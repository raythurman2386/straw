package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

// Client handles IPC connection to daemon.
type Client struct {
	socketPath string
	conn       net.Conn
	mu         sync.Mutex
	pending    map[string]chan Response
	events     chan Event
	done       chan struct{}
}

// NewClient creates a new IPC client.
func NewClient(socketPath string) *Client {
	return &Client{
		socketPath: socketPath,
		pending:    make(map[string]chan Response),
		events:     make(chan Event, 10), // buffer events
		done:       make(chan struct{}),
	}
}

// Connect establishes connection to the socket.
func (c *Client) Connect() error {
	slog.Debug("Connecting to Unix socket", "path", c.socketPath)
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.readLoop()
	return nil
}

// Close closes the connection.
func (c *Client) Close() {
	slog.Debug("Closing IPC client")
	close(c.done)
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) readLoop() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		line := scanner.Bytes()

		// Try to decode as Response first
		var resp Response
		if err := json.Unmarshal(line, &resp); err == nil && resp.ID != "" {
			c.mu.Lock()
			ch, ok := c.pending[resp.ID]
			c.mu.Unlock()
			if ok {
				ch <- resp
			} else {
				slog.Warn("Received response for unknown request ID", "id", resp.ID)
			}
			continue
		}

		// Try to decode as Event
		var event Event
		if err := json.Unmarshal(line, &event); err == nil && event.Type != "" {
			slog.Debug("Received IPC event", "type", event.Type)
			select {
			case c.events <- event:
			default:
				slog.Warn("Event buffer full, dropping event", "type", event.Type)
			}
			continue
		}

		slog.Warn("Received unknown message format", "payload", string(line))
	}

	if err := scanner.Err(); err != nil {
		slog.Error("IPC read loop error", "error", err)
	}
}

// Call sends a request and waits for a response.
func (c *Client) Call(method Method, params interface{}) (json.RawMessage, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())

	req := Request{
		ID:     id,
		Method: method,
	}
	if params != nil {
		p, _ := json.Marshal(params)
		req.Params = p
	}

	bytes, _ := json.Marshal(req)
	bytes = append(bytes, '\n')

	respCh := make(chan Response, 1)

	c.mu.Lock()
	c.pending[id] = respCh
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	slog.Debug("Sending IPC request", "method", method, "id", id)
	if _, err := c.conn.Write(bytes); err != nil {
		return nil, err
	}

	select {
	case resp := <-respCh:
		if resp.Error != "" {
			slog.Debug("IPC request returned error", "method", method, "id", id, "error", resp.Error)
			return nil, fmt.Errorf("%s", resp.Error)
		}
		slog.Debug("IPC request successful", "method", method, "id", id)
		return resp.Result, nil
	case <-time.After(5 * time.Second):
		slog.Error("IPC request timed out", "method", method, "id", id)
		return nil, fmt.Errorf("timeout")
	}
}

// Events returns the channel for asynchronous events.
func (c *Client) Events() <-chan Event {
	return c.events
}
