package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
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

	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
