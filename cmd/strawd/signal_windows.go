//go:build windows

package main

import (
	"log/slog"
	"os"
	"os/signal"
)

// waitForSignal blocks until a termination signal is received.
// On Windows, only os.Interrupt (Ctrl+C) is supported.
// Config reload is handled via the TRIGGER_RELOAD IPC command.
func waitForSignal(reloadConfig func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	sig := <-sigCh
	slog.Info("Shutting down", "signal", sig)
}
