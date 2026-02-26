package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

// Setup initializes the global slog logger.
// If logFile is provided, it will write to that file in addition to any other writers.
func Setup(level slog.Level, logFilePath string, verbose bool) error {
	var writers []io.Writer

	if verbose {
		writers = append(writers, os.Stderr)
	}

	if logFilePath != "" {
		// Ensure directory exists
		dir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writers = append(writers, f)
	}

	if len(writers) == 0 {
		writers = append(writers, io.Discard)
	}

	multiWriter := io.MultiWriter(writers...)

	// Use charmbracelet/log as the handler for slog
	l := log.New(multiWriter)

	// If we're writing to a file, we might want to disable colors for that output
	// if the TUI's Log tab doesn't handle them.
	// However, charmbracelet/log handles it for the whole logger.
	// For now, let's just make it look good but structured.
	l.SetStyles(log.DefaultStyles())
	l.SetReportTimestamp(true)
	l.SetTimeFormat("2006-01-02 15:04:05")

	// Map slog.Level to log.Level
	switch level {
	case slog.LevelDebug:
		l.SetLevel(log.DebugLevel)
	case slog.LevelInfo:
		l.SetLevel(log.InfoLevel)
	case slog.LevelWarn:
		l.SetLevel(log.WarnLevel)
	case slog.LevelError:
		l.SetLevel(log.ErrorLevel)
	default:
		l.SetLevel(log.InfoLevel)
	}

	logger := slog.New(l)
	slog.SetDefault(logger)

	return nil
}
