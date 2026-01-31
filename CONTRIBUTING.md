# Contributing to SparkyFitness MCP Server

This guide is for developers who want to contribute to or modify the SparkyFitness MCP server.

## Architecture

### Package Structure

- `/cmd/sparkyfitness-mcp` - Main entry point
- `/internal/config` - Configuration management
- `/internal/sparkyfitness` - Manual HTTP API client implementation
- `/internal/tools` - MCP tool implementations (search_foods, create_food_variant, add_food_variant)
- `/internal/logger` - Structured logging with slog

### API Client Implementation

The SparkyFitness API client is **manually implemented** using standard Go `net/http`:

- `internal/sparkyfitness/client.go` - HTTP client with Bearer token authentication
- `internal/sparkyfitness/types.go` - Request/response type definitions
- No code generation - simple, maintainable code

See `docs/backend_api.md` for API endpoint documentation.

## Prerequisites

- Go 1.25 or later
- SparkyFitness API access (URL and API key)
- Make (optional, for using Makefile commands)

## Development Setup

### Clone the Repository

```bash
git clone https://github.com/chickenzord/sparkyfitness-mcp.git
cd sparkyfitness-mcp
```

### Install Dependencies

```bash
go mod download
```

## Building

### Using Make

```bash
make build
```

### Using Go Directly

```bash
go build -o sparkyfitness-mcp ./cmd/sparkyfitness-mcp
```

## Testing

### Run All Tests

```bash
make test
# or
go test ./...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
make test-coverage
# or
go test -cover ./...
```

## Running Locally

### Stdio Transport (for Claude Desktop development)

```bash
# Set environment variables
export SPARKYFITNESS_API_URL=http://localhost:8000
export SPARKYFITNESS_API_KEY=your-api-key

# Run the server
make run
# or
./sparkyfitness-mcp
```

### HTTP Transport (for web/remote development)

```bash
export SPARKYFITNESS_API_URL=http://localhost:8000
export SPARKYFITNESS_API_KEY=your-api-key
export MCP_TRANSPORT=http
export MCP_HTTP_HOST=0.0.0.0
export MCP_HTTP_PORT=8080

./sparkyfitness-mcp
```

The server will start on `http://0.0.0.0:8080` with:
- MCP endpoint: `http://localhost:8080/mcp/`
- Health check: `http://localhost:8080/health`

## Docker Development

### Build Docker Image

```bash
make docker-build
# or
docker build -t sparkyfitness-mcp .
```

### Run Docker Container (Stdio)

```bash
docker run \
  -e SPARKYFITNESS_API_URL=http://localhost:8000 \
  -e SPARKYFITNESS_API_KEY=your-api-key \
  sparkyfitness-mcp
```

### Run Docker Container (HTTP)

```bash
docker run -p 8080:8080 \
  -e SPARKYFITNESS_API_URL=http://localhost:8000 \
  -e SPARKYFITNESS_API_KEY=your-api-key \
  -e MCP_TRANSPORT=http \
  -e MCP_HTTP_HOST=0.0.0.0 \
  -e MCP_HTTP_PORT=8080 \
  sparkyfitness-mcp
```

## Code Quality

### Format Code

```bash
make fmt
# or
go fmt ./...
```

### Run Linter

```bash
make vet
# or
go vet ./...
```

### Run All Checks

```bash
make check
```

This runs formatting, linting, and tests.

## Project Guidelines

### Design Principles

- **Stateless**: No local storage or database
- **Pure API Bridge**: MCP protocol → SparkyFitness backend API
- **Simple & Maintainable**: Manual HTTP client, no code generation
- **Configuration via Environment**: No configuration files

### Adding New MCP Tools

1. Create a new file in `internal/tools/` (e.g., `my_new_tool.go`)
2. Define input/output structs with JSON schema tags
3. Implement the tool handler function
4. Register the tool in `internal/tools/registry.go`
5. Add tests for the new tool
6. Update documentation in README.md and CLAUDE.md

### API Client Development

When adding new backend API endpoints:

1. Add request/response types to `internal/sparkyfitness/types.go`
2. Implement the client method in `internal/sparkyfitness/client.go`
3. Use Bearer token authentication (handled by `authInterceptor`)
4. Return descriptive errors for API failures
5. Document the endpoint in `docs/backend_api.md`

### Testing Strategy

- Unit tests for business logic
- Integration tests for API client (use test fixtures)
- Manual testing with real backend (use test scripts in project root)

## Environment Variables

### Required

- `SPARKYFITNESS_API_URL` - Base URL of SparkyFitness server component (backend), not the frontend. Example: `http://localhost:8000`
- `SPARKYFITNESS_API_KEY` - Authentication credential for SparkyFitness API

### Optional

- `MCP_TRANSPORT` - Transport mode: `stdio` (default) or `http`
- `MCP_HTTP_HOST` - Host to bind to when using HTTP transport (default: `0.0.0.0`)
- `MCP_HTTP_PORT` - Port to listen on when using HTTP transport (default: `8080`)
- `MCP_HTTP_BASIC_AUTH_USER` - Username for HTTP basic auth (optional)
- `MCP_HTTP_BASIC_AUTH_PASSWORD` - Password for HTTP basic auth (optional)

## Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run code quality checks (`make check`)
6. Commit your changes (`git commit -am 'Add new feature'`)
7. Push to the branch (`git push origin feature/my-feature`)
8. Create a Pull Request

## Questions?

For questions or issues, please open an issue on GitHub.
