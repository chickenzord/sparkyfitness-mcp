package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
	"github.com/chickenzord/sparkyfitness-mcp/internal/logger"
	"github.com/chickenzord/sparkyfitness-mcp/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration from environment variables
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger based on configuration
	logger.InitLogger(cfg)

	slog.Info("Starting SparkyFitness MCP Server",
		"transport", cfg.Transport,
		"log_level", cfg.LogLevel,
		"log_format", cfg.LogFormat,
	)

	// Create the MCP server
	srv, err := server.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Set up context with signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()

	// Run the server
	slog.Info("Server ready")
	if err := srv.Run(ctx); err != nil {
		slog.Error("Server error", "error", err)
		return err
	}

	slog.Info("Server stopped gracefully")
	return nil
}
