package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opencode-ai/opencode/internal/app"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/daemon"
	"github.com/opencode-ai/opencode/internal/db"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/spf13/cobra"
)

var (
	daemonPort    int
	daemonHost    string
	obsidianVault string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start OpenCode headless daemon and HTTP/JSON-RPC server",
	Long: `Run OpenCode as a headless background service exposing REST and JSON-RPC APIs.
Enables seamless integration with LangGraph, Obsidian Command Centre, scripts, and CI/CD pipelines.`,
	Example: `
  # Start daemon on default port 8080
  opencode daemon

  # Start daemon on port 9000 with Obsidian vault watcher
  opencode daemon --port 9000 --watch-obsidian /path/to/Obsidian-Vault
  `,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := cmd.Flags().GetString("cwd")
		if cwd != "" {
			_ = os.Chdir(cwd)
		} else {
			c, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %v", err)
			}
			cwd = c
		}

		_, err := config.Load(cwd, false)
		if err != nil {
			return err
		}

		conn, err := db.Connect()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		appInstance, err := app.New(ctx, conn)
		if err != nil {
			return fmt.Errorf("failed to initialize app: %w", err)
		}
		defer appInstance.Shutdown()

		// Handle graceful shutdown signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			logging.Info("Shutting down daemon...")
			cancel()
		}()

		// Start Obsidian file watcher if requested
		if obsidianVault != "" {
			watcher := daemon.NewObsidianTaskWatcher(obsidianVault, appInstance)
			go watcher.Start(ctx)
		}

		server := daemon.NewDaemonServer(appInstance, daemonHost, daemonPort)
		return server.Start(ctx)
	},
}

func init() {
	daemonCmd.Flags().IntVarP(&daemonPort, "port", "P", 8080, "Port for daemon HTTP server")
	daemonCmd.Flags().StringVar(&daemonHost, "host", "0.0.0.0", "Host address to bind to")
	daemonCmd.Flags().StringVar(&obsidianVault, "watch-obsidian", "", "Path to Obsidian Vault to watch for tasks")
	rootCmd.AddCommand(daemonCmd)
}
