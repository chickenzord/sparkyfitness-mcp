# SparkyFitness MCP Server

A stateless Go-based MCP (Model Context Protocol) server that acts as an API bridge between Claude Chat and the SparkyFitness backend. This server enables you to create and search food items in SparkyFitness by uploading nutrition label photos to Claude Chat.

## Features

- **Vision-Powered Food Entry**: Upload nutrition labels to Claude Chat, which automatically extracts nutrition data and creates food entries
- **Smart Search**: Find existing foods to avoid duplicates
- **Variant Management**: Add multiple serving sizes to the same food (e.g., 100g, 150g, 1 cup)
- **Dual Transport Support**:
  - **stdio**: For local Claude Desktop integration
  - **HTTP/SSE**: For remote deployment and claude.ai web integration
- **Stateless Design**: Pure API translation layer with no local storage

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

### `search_foods`

Search for foods in the SparkyFitness database to avoid duplicates.

**Example:**
```
User: "Find nutrition info for chicken breast"
Claude: [Uses search_foods to find existing entries]
```

### `create_food_variant`

Create a completely new food entry with its first serving size variant.

**Example:**
```
User: [Uploads photo of nutrition label]
Claude: [Extracts data, searches for duplicates, then creates new entry if none found]
```

### `add_food_variant`

Add a new serving size variant to an existing food.

**Example:**
```
User: "I have nutrition data for Enoki Mushroom 150g serving"
Claude: "Found existing Enoki Mushroom with 100g variant. Add 150g variant?"
User: "Yes"
Claude: [Adds 150g variant to the existing food]
```

## Usage Examples

### Adding a New Food

1. Upload a nutrition label photo to Claude Chat
2. Claude automatically:
   - Extracts nutrition facts using vision
   - Searches for duplicates
   - Creates the food entry if it doesn't exist
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
