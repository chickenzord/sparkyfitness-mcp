package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
	"github.com/chickenzord/sparkyfitness-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "sparkyfitness-mcp"
	serverVersion = "0.1.0"
	serverTitle   = "SparkyFitness MCP Server"
)

// Server wraps the MCP server and application configuration
type Server struct {
	mcp      *mcp.Server
	config   *config.Config
	registry *tools.Registry
}

// New creates a new SparkyFitness MCP server
func New(cfg *config.Config) (*Server, error) {
	// Create tool registry first
	registry := tools.NewRegistry(cfg)

	// Create MCP server implementation
	impl := &mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
		Title:   serverTitle,
	}

	mcpServer := mcp.NewServer(impl, nil)

	// Register all tools
	if err := registry.RegisterAll(mcpServer); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return &Server{
		mcp:      mcpServer,
		config:   cfg,
		registry: registry,
	}, nil
}

// createMCPServer creates a new MCP server instance with tools registered
// This is used for HTTP transport where each connection may need a separate server instance
func (s *Server) createMCPServer() (*mcp.Server, error) {
	impl := &mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
		Title:   serverTitle,
	}

	mcpServer := mcp.NewServer(impl, nil)

	// Register all tools
	if err := s.registry.RegisterAll(mcpServer); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return mcpServer, nil
}

// Run starts the MCP server using the configured transport
func (s *Server) Run(ctx context.Context) error {
	switch s.config.Transport {
	case config.TransportStdio:
		return s.runStdio(ctx)
	case config.TransportHTTP:
		return s.runHTTP(ctx)
	default:
		return fmt.Errorf("unsupported transport mode: %s", s.config.Transport)
	}
}

// runStdio starts the MCP server over stdio transport
func (s *Server) runStdio(ctx context.Context) error {
	// Use stdio transport for MCP communication
	transport := &mcp.StdioTransport{}

	// Run the server until the client disconnects or context is cancelled
	if err := s.mcp.Run(ctx, transport); err != nil {
		return fmt.Errorf("stdio server error: %w", err)
	}

	return nil
}

// runHTTP starts the MCP server over HTTP using the SDK's built-in StreamableHTTPHandler
func (s *Server) runHTTP(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", s.config.HTTPHost, s.config.HTTPPort)

	// Create the StreamableHTTPHandler with SDK's built-in support
	// Use a stateless configuration for simplicity
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			// Return the MCP server for each request
			// For stateless mode, we can reuse the same server instance
			return s.mcp
		},
		&mcp.StreamableHTTPOptions{
			Stateless:      true,
			JSONResponse:   false,
			SessionTimeout: 30 * time.Minute,
		},
	)

	// Create HTTP mux with health check and MCP handler
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// MCP endpoint - delegate to the StreamableHTTPHandler with optional basic auth
	mcpHandler := http.StripPrefix("/mcp", handler)
	mux.Handle("/mcp/", s.config.BasicAuthMiddleware(mcpHandler))

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Channel to signal server errors
	errChan := make(chan error, 1)

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting MCP server (HTTP) on http://%s", addr)
		log.Printf("MCP endpoint: http://%s/mcp/", addr)
		if s.config.BasicAuthEnabled() {
			log.Printf("HTTP Basic Authentication: enabled")
		} else {
			log.Printf("HTTP Basic Authentication: disabled")
		}
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("HTTP server shutdown error: %w", err)
		}
		return nil
	case err := <-errChan:
		return err
	}
}
