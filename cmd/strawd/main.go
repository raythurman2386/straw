package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"straw/internal/actions"
	"straw/internal/config"
	"straw/internal/ipc"
	"straw/internal/logging"
	"straw/internal/rules"
	"straw/internal/watcher"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var configPath string
	var socketPath string
	var logFilePath string
	var verbose bool

	rootCmd := &cobra.Command{
		Use:     "strawd",
		Short:   "Straw daemon",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level := slog.LevelInfo
			if verbose {
				level = slog.LevelDebug
			}

			if logFilePath == "" {
				stateDir, err := config.DefaultStateDir()
				if err == nil {
					logFilePath = filepath.Join(stateDir, "strawd.log")
				}
			}

			return logging.Setup(level, logFilePath, true) // Always log to stderr in daemon for now, or maybe only if verbose?
			// Actually, let's log to stderr if it's attached to a terminal or if verbose is set.
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Starting strawd")

			// Load Config
			if configPath == "" {
				var err error
				configPath, err = config.DefaultConfigPath()
				if err != nil {
					return fmt.Errorf("failed to get default config path: %w", err)
				}
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Determine socket path
			sock := socketPath
			if sock == "" {
				sock = cfg.SocketPath
			}
			if sock == "" {
				sock = config.DefaultSocketPath()
			}

			slog.Info("IPC server address", "socket", sock, "config", configPath)

			// Initialize IPC Server
			server := ipc.NewServer(sock)

			// Start Watcher
			w, err := watcher.New()
			if err != nil {
				return fmt.Errorf("failed to create watcher: %w", err)
			}
			defer w.Close()

			// Helper to setup watches
			setupWatches := func(c *config.Config) {
				for _, watchCfg := range c.Watch {
					slog.Info("Adding watch", "path", watchCfg.Path, "recursive", watchCfg.Recursive)
					if err := w.Add(watchCfg.Path, watchCfg.Recursive); err != nil {
						slog.Error("Failed to add watch", "path", watchCfg.Path, "error", err)
					}
				}
			}
			setupWatches(cfg)
			w.Start()

			// Initialize Rules Engine & Executor
			engine := rules.NewEngine(cfg.Rules)
			executor := actions.NewExecutor()

			// Reload Logic
			reloadConfig := func() {
				slog.Info("Reloading configuration")
				newCfg, err := config.Load(configPath)
				if err != nil {
					slog.Error("Failed to reload config", "error", err)
					return
				}

				// Update Local reference for GET_RULES
				cfg = newCfg

				// Update Rules
				engine.SetRules(newCfg.Rules)

				// Update Watcher
				setupWatches(newCfg)

				slog.Info("Configuration reloaded successfully")
			}

			// Register handlers
			server.Register(ipc.MethodGetStatus, func(params json.RawMessage) (interface{}, error) {
				return map[string]string{"status": "running", "version": version}, nil
			})

			server.Register(ipc.MethodGetRules, func(params json.RawMessage) (interface{}, error) {
				return cfg.Rules, nil
			})

			server.Register(ipc.MethodAddRule, func(params json.RawMessage) (interface{}, error) {
				var newRule config.Rule
				if err := json.Unmarshal(params, &newRule); err != nil {
					return nil, err
				}

				// Persist to config file
				data, err := os.ReadFile(configPath)
				if err != nil {
					return nil, fmt.Errorf("failed to read config for update: %w", err)
				}

				var fullCfg config.Config
				if err := toml.Unmarshal(data, &fullCfg); err != nil {
					return nil, fmt.Errorf("failed to parse config for update: %w", err)
				}

				fullCfg.Rules = append(fullCfg.Rules, newRule)

				newData, err := toml.Marshal(fullCfg)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal updated config: %w", err)
				}

				if err := os.WriteFile(configPath, newData, 0644); err != nil {
					return nil, fmt.Errorf("failed to write updated config: %w", err)
				}

				slog.Info("New rule added and persisted", "name", newRule.Name)
				go reloadConfig()

				return "rule added", nil
			})

			server.Register(ipc.MethodUpdateRule, func(params json.RawMessage) (interface{}, error) {
				var args ipc.UpdateRuleParams
				if err := json.Unmarshal(params, &args); err != nil {
					return nil, err
				}

				// Persist to config file
				data, err := os.ReadFile(configPath)
				if err != nil {
					return nil, fmt.Errorf("failed to read config for update: %w", err)
				}

				var fullCfg config.Config
				if err := toml.Unmarshal(data, &fullCfg); err != nil {
					return nil, fmt.Errorf("failed to parse config for update: %w", err)
				}

				found := false
				for i, r := range fullCfg.Rules {
					if r.Name == args.OriginalName {
						fullCfg.Rules[i] = args.Rule
						found = true
						break
					}
				}

				if !found {
					return nil, fmt.Errorf("rule not found: %s", args.OriginalName)
				}

				newData, err := toml.Marshal(fullCfg)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal updated config: %w", err)
				}

				if err := os.WriteFile(configPath, newData, 0644); err != nil {
					return nil, fmt.Errorf("failed to write updated config: %w", err)
				}

				slog.Info("Rule updated and persisted", "original_name", args.OriginalName, "new_name", args.Rule.Name)
				go reloadConfig()

				return "rule updated", nil
			})

			server.Register(ipc.MethodTriggerReload, func(params json.RawMessage) (interface{}, error) {
				go reloadConfig()
				return "reload triggered", nil
			})

			if err := server.Start(); err != nil {
				return fmt.Errorf("failed to start server: %w", err)
			}
			defer server.Stop()

			// Event Loop
			go func() {
				for {
					select {
					case event, ok := <-w.Events():
						if !ok {
							return
						}

						slog.Debug("FileSystem event", "path", event.Path, "type", event.Type)
						actionsList := engine.Evaluate(event)
						for _, action := range actionsList {
							slog.Info("Executing action", "type", action.Type, "file", event.Path)
							err := executor.Execute(action, event.Path)

							payload := map[string]interface{}{
								"file":   event.Path,
								"action": action.Type,
							}

							if err != nil {
								slog.Error("Action failed", "type", action.Type, "file", event.Path, "error", err)
								payload["status"] = "error"
								payload["error"] = err.Error()
							} else {
								slog.Info("Action completed successfully", "type", action.Type, "file", event.Path)
								payload["status"] = "success"
							}

							server.Broadcast(ipc.EventNotification, payload)
						}

					case err, ok := <-w.Errors():
						if !ok {
							return
						}
						slog.Error("Watcher error", "error", err)
					}
				}
			}()

			// Wait for signal
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

			for sig := range sigCh {
				if sig == syscall.SIGHUP {
					reloadConfig()
				} else {
					slog.Info("Shutting down", "signal", sig)
					return nil
				}
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")
	rootCmd.PersistentFlags().StringVar(&socketPath, "socket", "", "Override socket path")
	rootCmd.PersistentFlags().StringVar(&logFilePath, "log-file", "", "Path to log file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Override Cobra's default version flag to avoid conflict with -v (verbose)
	rootCmd.Flags().Bool("version", false, "Print the version")
	rootCmd.SetVersionTemplate(fmt.Sprintf("strawd version {{.Version}} (commit: %s, built: %s)\n", commit, date))

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load(configPath)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "config OK")
			return nil
		},
	}

	rootCmd.AddCommand(checkCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Execution failed", "error", err)
		os.Exit(1)
	}
}
