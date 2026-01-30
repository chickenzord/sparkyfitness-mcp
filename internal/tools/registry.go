package tools

import (
	"fmt"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
	"github.com/chickenzord/sparkyfitness-mcp/internal/sparkyfitness"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry manages all MCP tools
type Registry struct {
	config *config.Config
	client *sparkyfitness.Client
}

// NewRegistry creates a new tool registry
func NewRegistry(cfg *config.Config) *Registry {
	return &Registry{
		config: cfg,
	}
}

// RegisterAll registers all available tools with the MCP server
func (r *Registry) RegisterAll(server *mcp.Server) error {
	// Initialize SparkyFitness API client
	client, err := sparkyfitness.NewClient(r.config)
	if err != nil {
		return fmt.Errorf("failed to create SparkyFitness client: %w", err)
	}
	r.client = client

	// Register search_foods tool (sfmcp-tcr)
	if err := r.RegisterSearchFoods(server, client); err != nil {
		return fmt.Errorf("failed to register search_foods: %w", err)
	}

	// Register add_food_variant tool (sfmcp-248.3)
	if err := r.RegisterAddFoodVariant(server, client); err != nil {
		return fmt.Errorf("failed to register add_food_variant: %w", err)
	}

	// Register create_food_variant tool (sfmcp-248.4, sfmcp-boc)
	if err := r.RegisterCreateFoodVariant(server, client); err != nil {
		return fmt.Errorf("failed to register create_food_variant: %w", err)
	}

	return nil
}
