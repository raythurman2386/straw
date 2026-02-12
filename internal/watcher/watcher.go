package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type EventType int

const (
	Create EventType = iota
	Write
	Remove
	Rename
	Chmod
)

type Event struct {
	Path string
	Type EventType
}

type Watcher struct {
	fsWatcher *fsnotify.Watcher
	events    chan Event
	errors    chan error
	done      chan struct{}
	watched   map[string]bool // set of all watched directories
	roots     map[string]bool // root paths -> recursive status
	mu        sync.Mutex
}

func New() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fsWatcher: w,
		events:    make(chan Event, 100),
		errors:    make(chan error, 10),
		done:      make(chan struct{}),
		watched:   make(map[string]bool),
		roots:     make(map[string]bool),
	}, nil
}

// Add adds a path to be watched.
func (w *Watcher) Add(path string, recursive bool) error {
	w.mu.Lock()
	w.roots[path] = recursive
	w.mu.Unlock()

	if !recursive {
		return w.addDir(path)
	}

	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories we can't read
			return nil
		}
		if info.IsDir() {
			return w.addDir(p)
		}
		return nil
	})
}

func (w *Watcher) addDir(path string) error {
	w.mu.Lock()
	if w.watched[path] {
		w.mu.Unlock()
		return nil
	}
	w.watched[path] = true
	w.mu.Unlock()

	return w.fsWatcher.Add(path)
}

func (w *Watcher) Start() {
	go w.loop()
}

func (w *Watcher) loop() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			// Handle new directories
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					// Check if inside a recursive root
					if w.isRecursive(event.Name) {
						if err := w.Add(event.Name, true); err != nil {
							w.errors <- err
						}
					}
				}
			}

			// Forward event
			// We might want to debounce or filter here, but for now just raw events.
			w.events <- translateEvent(event)

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			w.errors <- err
		case <-w.done:
			return
		}
	}
}

func (w *Watcher) isRecursive(path string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if path starts with any root that is recursive
	for root, recursive := range w.roots {
		if recursive {
			if strings.HasPrefix(path, root) {
				return true
			}
		}
	}
	return false
}

func (w *Watcher) Close() {
	close(w.done)
	w.fsWatcher.Close()
}

func (w *Watcher) Events() <-chan Event {
	return w.events
}

func (w *Watcher) Errors() <-chan error {
	return w.errors
}

func translateEvent(e fsnotify.Event) Event {
	var t EventType
	if e.Op&fsnotify.Create == fsnotify.Create {
		t = Create
	} else if e.Op&fsnotify.Write == fsnotify.Write {
		t = Write
	} else if e.Op&fsnotify.Remove == fsnotify.Remove {
		t = Remove
	} else if e.Op&fsnotify.Rename == fsnotify.Rename {
		t = Rename
	} else if e.Op&fsnotify.Chmod == fsnotify.Chmod {
		t = Chmod
	}
	return Event{Path: e.Name, Type: t}
}
