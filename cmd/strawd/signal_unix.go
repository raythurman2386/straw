//go:build !windows

package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// waitForSignal blocks until a termination signal is received.
// On Unix, SIGHUP triggers a config reload instead of shutdown.
func waitForSignal(reloadConfig func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	for sig := range sigCh {
		if sig == syscall.SIGHUP {
			reloadConfig()
		} else {
			slog.Info("Shutting down", "signal", sig)
			return
		}
	}
}
