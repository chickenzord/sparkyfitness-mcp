package tools

import (
	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry manages all MCP tools
type Registry struct {
	config *config.Config
}

// NewRegistry creates a new tool registry
func NewRegistry(cfg *config.Config) *Registry {
	return &Registry{
		config: cfg,
	}
}

// RegisterAll registers all available tools with the MCP server
func (r *Registry) RegisterAll(server *mcp.Server) error {
	// Tool implementations will be added in separate issues
	// This provides a clean registration point for:
	// - create_food_variant (sfmcp-boc)
	// - search_food_variants (sfmcp-tcr)

	return nil
}
