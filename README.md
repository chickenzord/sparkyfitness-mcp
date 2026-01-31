# SparkyFitness MCP Server

A stateless Go-based MCP (Model Context Protocol) server that acts as an API bridge between AI agents (like Claude Chat) and the SparkyFitness backend. This server provides tools to create and search food items with complete nutrition data.

## Features

- **Smart Search**: Find existing foods in the database to avoid duplicates
- **Food Creation**: Create new food entries with complete nutrition data
- **Variant Management**: Add multiple serving sizes to the same food (e.g., 100g, 150g, 1 cup)
- **Dual Transport Support**:
  - **stdio**: For local Claude Desktop integration
  - **HTTP/SSE**: For remote deployment and claude.ai web integration
- **Stateless Design**: Pure API translation layer with no local storage

**Works seamlessly with AI agents:** When used with Claude Chat or other vision-capable AI agents, users can upload nutrition label photos and have the AI automatically extract data before calling the MCP tools.

## Quick Start

### Using Claude Desktop (Local)

1. **Download or build the binary** (see [Releases](https://github.com/chickenzord/sparkyfitness-mcp/releases))

2. **Add to Claude Desktop configuration**:

   On macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`

   ```json
   {
     "mcpServers": {
       "sparkyfitness": {
         "command": "/path/to/sparkyfitness-mcp",
         "env": {
           "SPARKYFITNESS_API_URL": "http://localhost:8000",
           "SPARKYFITNESS_API_KEY": "your-api-key"
         }
       }
     }
   }
   ```

   **Note:** Replace `http://localhost:8000` with the actual URL of your SparkyFitness server component (backend).

3. **Restart Claude Desktop**

4. **Start using it**: Upload a nutrition label photo and ask Claude to add it to SparkyFitness!

### Using Docker

```bash
docker run -p 8080:8080 \
  -e SPARKYFITNESS_API_URL=http://localhost:8000 \
  -e SPARKYFITNESS_API_KEY=your-api-key \
  -e MCP_TRANSPORT=http \
  -e MCP_HTTP_HOST=0.0.0.0 \
  -e MCP_HTTP_PORT=8080 \
  sparkyfitness-mcp
```

**Note:** Replace `http://localhost:8000` with the actual URL of your SparkyFitness server component (backend).

Then connect from claude.ai web using the MCP endpoint: `http://localhost:8080/mcp/`

## Configuration

### Required Environment Variables

| Variable | Description |
|----------|-------------|
| `SPARKYFITNESS_API_URL` | Base URL of the SparkyFitness **server component** (backend), not the frontend. Example: `http://localhost:8000` or `https://sparkyfitness-server.example.com` |
| `SPARKYFITNESS_API_KEY` | Your API key for authentication |

### Optional Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_TRANSPORT` | `stdio` | Transport mode: `stdio` or `http` |
| `MCP_HTTP_HOST` | `0.0.0.0` | Host to bind to (HTTP mode only) |
| `MCP_HTTP_PORT` | `8080` | Port to listen on (HTTP mode only) |
| `MCP_HTTP_BASIC_AUTH_USER` | - | Username for HTTP basic auth (optional) |
| `MCP_HTTP_BASIC_AUTH_PASSWORD` | - | Password for HTTP basic auth (optional) |

## Available Tools

This MCP server provides three tools that Claude can use:

### 🔍 `search_foods`

Search for existing foods in the SparkyFitness nutrition database. **Always call this first** before creating any food to prevent duplicates.

**What it does:**
- Searches by food name and optionally filters by brand
- Returns matching foods with complete nutrition data for their default variants
- Provides `food_id` values needed for adding variants to existing foods

**Search modes:**
- `broad_match=true` (default): Fuzzy case-insensitive search that finds similar foods
- `broad_match=false`: Exact matching for precise results

**Example:**
```
User: "Find nutrition info for chicken breast"
Claude: [Calls search_foods]
Result: List of matching foods with nutrition facts
```

### 🆕 `create_food_variant`

Create a **completely new** food entry with its first serving size variant. Only use when no matching food exists or user explicitly wants a separate entry.

**When to use:**
- `search_foods` found NO matches (no duplicates exist), OR
- User explicitly chooses to create a separate food entry despite duplicates

**What it does:**
- Creates a new food entity in the database
- Adds the first serving size variant as the default
- Returns both `food_id` and `variant_id`

**Important:** Not for adding variants to existing foods - use `add_food_variant` for that.

**Example:**
```
User: [Uploads photo of nutrition label for "Organic Quinoa"]
Claude: [Searches first, finds no matches]
Claude: [Creates new food entry with nutrition data]
Result: New food created successfully
```

### ➕ `add_food_variant`

Add a new serving size variant to an **existing** food. Use this to add alternative serving sizes to foods found via `search_foods`.

**When to use:**
- `search_foods` found a matching food, AND
- User wants to add a new serving size to that existing food (not create a separate entry)

**What it does:**
- Adds another serving size option to an existing food entry
- Requires `food_id` from search results
- Example: Food has 100g variant, add 150g variant to the same food

**Example:**
```
User: "I have nutrition data for Enoki Mushroom 150g serving"
Claude: [Searches and finds existing Enoki Mushroom with 100g variant]
Claude: "Found existing Enoki Mushroom with 100g variant. Add 150g variant?"
User: "Yes"
Claude: [Calls add_food_variant with food_id and nutrition data]
Result: Enoki Mushroom now has TWO variants (100g and 150g)
```

## Usage Examples

### Adding a New Food (with Claude Chat)

1. Upload a nutrition label photo to Claude Chat
2. Claude Chat uses its vision capabilities to extract nutrition facts from the photo
3. Claude Chat calls the MCP tools:
   - Searches for duplicates using `search_foods`
   - Creates the food entry using `create_food_variant` if no duplicates found
   - Confirms success with food ID

### Searching for Foods

```
User: "What's the nutrition info for quinoa?"
Claude: [Searches and shows available options with nutrition details]
```

### Managing Variants

```
User: "Add a 1 cup serving size for Brown Rice"
Claude: [Searches, finds existing food, asks to confirm, adds variant]
```

## Security

### HTTP Basic Authentication

When running in HTTP mode, you can optionally enable basic authentication:

```bash
export MCP_HTTP_BASIC_AUTH_USER=admin
export MCP_HTTP_BASIC_AUTH_PASSWORD=your-secret-password
```

**Note**: Both username and password must be set to enable authentication. If either is missing, authentication is disabled.

## Health Check

When running in HTTP mode, the server provides a health check endpoint:

```bash
curl http://localhost:8080/health
```

## Contributing

For development setup, building from source, and contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

See [LICENSE](LICENSE) for details.

## Support

For issues, questions, or feature requests, please open an issue on [GitHub](https://github.com/chickenzord/sparkyfitness-mcp/issues).
