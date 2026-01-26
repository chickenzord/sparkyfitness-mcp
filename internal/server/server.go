package server

import (
	"context"
	"fmt"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
	"github.com/chickenzord/sparkyfitness-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "sparkyfitness-mcp"
	serverVersion = "0.1.0"
)

// Server wraps the MCP server and application configuration
type Server struct {
	mcp      *mcp.Server
	config   *config.Config
	registry *tools.Registry
}

// New creates a new SparkyFitness MCP server
func New(cfg *config.Config) (*Server, error) {
	// Create MCP server implementation
	impl := &mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}

	mcpServer := mcp.NewServer(impl, nil)

	// Create tool registry
	registry := tools.NewRegistry(cfg)

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

// Run starts the MCP server over stdio transport
func (s *Server) Run(ctx context.Context) error {
	// Use stdio transport for MCP communication
	transport := &mcp.StdioTransport{}

	// Run the server until the client disconnects or context is cancelled
	if err := s.mcp.Run(ctx, transport); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
