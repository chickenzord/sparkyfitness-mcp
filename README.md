# SparkyFitness MCP Server

A stateless Go-based MCP (Model Context Protocol) server that acts as an API bridge between Claude Chat and the SparkyFitness backend. This server enables users to create and search food items in SparkyFitness by uploading nutrition label photos to Claude Chat.

## Features

- **Stateless Design**: Pure API translation layer with no local storage
- **Dual Transport Support**:
  - **stdio**: For local Claude Desktop integration
  - **HTTP/SSE**: For remote deployment and claude.ai web integration
- **MCP Tools**:
  - `search_foods`: Search for existing foods in the database
  - `create_food_variant`: Create new food entries with nutrition data
- **Vision Integration**: Upload nutrition labels to Claude Chat, which extracts data and creates food entries

## Architecture

- **Language**: Go
- **Design**: Stateless API translation layer
- **Package Structure**:
  - `/cmd/sparkyfitness-mcp` - Main entry point
  - `/internal/server` - MCP server implementation
  - `/internal/config` - Configuration management
  - `/internal/sparkyfitness` - SparkyFitness API client
  - `/internal/tools` - MCP tool implementations

## Prerequisites

- Go 1.23 or later
- SparkyFitness API access (URL and API key)

## Configuration

### Environment Variables

#### Required

| Variable | Description |
|----------|-------------|
| `SPARKYFITNESS_API_URL` | Base URL of the SparkyFitness API |
| `SPARKYFITNESS_API_KEY` | Authentication credential for the API |

#### Optional (Transport Configuration)

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_TRANSPORT` | `stdio` | Transport mode: `stdio` or `http` |
| `MCP_HTTP_HOST` | `0.0.0.0` | Host to bind to (HTTP mode only) |
| `MCP_HTTP_PORT` | `8080` | Port to listen on (HTTP mode only) |

## Development

### Generate API Client

The SparkyFitness API client is auto-generated from the OpenAPI specification:

```bash
make generate
```

This uses `oapi-codegen` to generate type-safe client code from `swagger.json`. The generated code should not be edited manually.

### Build

```bash
make build
# or
go build -o sparkyfitness-mcp ./cmd/sparkyfitness-mcp
```

### Test

```bash
make test
# or
go test ./...
```

### Run

#### Stdio Transport (Claude Desktop)

Default mode for local Claude Desktop integration:

```bash
# Set environment variables
export SPARKYFITNESS_API_URL=https://api.sparkyfitness.com
export SPARKYFITNESS_API_KEY=your-api-key

# Run the server
make run
# or
./sparkyfitness-mcp
```

Add to Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "sparkyfitness": {
      "command": "/path/to/sparkyfitness-mcp",
      "env": {
        "SPARKYFITNESS_API_URL": "https://api.sparkyfitness.com",
        "SPARKYFITNESS_API_KEY": "your-api-key"
      }
    }
  }
}
```

#### HTTP Transport (Remote / claude.ai Web)

For remote deployment or claude.ai web integration:

```bash
export SPARKYFITNESS_API_URL=https://api.sparkyfitness.com
export SPARKYFITNESS_API_KEY=your-api-key
export MCP_TRANSPORT=http
export MCP_HTTP_HOST=0.0.0.0
export MCP_HTTP_PORT=8080

./sparkyfitness-mcp
```

The server will start on `http://0.0.0.0:8080` with:
- MCP endpoint: `http://localhost:8080/mcp/`
- Health check: `http://localhost:8080/health`

## Docker

### Build Image

```bash
make docker-build
# or
docker build -t sparkyfitness-mcp .
```

### Run Container

#### Stdio Transport

```bash
docker run \
  -e SPARKYFITNESS_API_URL=https://api.sparkyfitness.com \
  -e SPARKYFITNESS_API_KEY=your-api-key \
  sparkyfitness-mcp
```

#### HTTP Transport

```bash
docker run -p 8080:8080 \
  -e SPARKYFITNESS_API_URL=https://api.sparkyfitness.com \
  -e SPARKYFITNESS_API_KEY=your-api-key \
  -e MCP_TRANSPORT=http \
  -e MCP_HTTP_HOST=0.0.0.0 \
  -e MCP_HTTP_PORT=8080 \
  sparkyfitness-mcp
```

## MCP Tools

### `search_foods`

Search for foods in the SparkyFitness database by name and optional filters.

**Parameters:**
- `name` (string, required) - Food name to search for
- `brand` (string, optional) - Brand name to filter results
- `broad_match` (boolean, optional) - If true, performs broad matching (default: true)
- `limit` (integer, optional) - Maximum number of results to return (default: 10)

### `create_food_variant`

Create a new food entry with nutrition data. Upserts the Food entity by name+brand.

**Parameters:**
- `name` (string, required) - Food name
- `brand` (string, optional) - Brand name
- `serving_size` (number, required) - Numeric amount (e.g., 100, 1)
- `serving_unit` (string, required) - Unit (e.g., "g", "ml", "cup", "piece")
- `calories` (number, required) - Calories per serving
- `protein_g` (number, required) - Protein in grams
- `carbs_g` (number, required) - Carbohydrates in grams
- `fat_g` (number, required) - Fat in grams
- `fiber_g` (number, optional) - Fiber in grams
- `sugar_g` (number, optional) - Sugar in grams
- `sodium_mg` (number, optional) - Sodium in milligrams
- `set_default` (boolean, optional) - Set this variant as the food's default variant

## Usage Example

1. **Upload a nutrition label photo to Claude Chat**
2. **Claude extracts the nutrition facts automatically using vision**
3. **Claude calls `create_food_variant` with the structured data**
4. **The food entry is created in SparkyFitness**

Or search for existing foods:

```
User: "Find nutrition info for chicken breast"
Claude: [Calls search_foods tool]
Claude: "I found several chicken breast options..."
```

