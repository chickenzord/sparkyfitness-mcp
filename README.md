# sparkyfitness-mcp

MCP server companion for SparkyFitness - A stateless API bridge that enables creating food entries in SparkyFitness via the Model Context Protocol.

## Overview

This MCP server acts as a translation layer between the Model Context Protocol and SparkyFitness REST API. It enables MCP clients to create and search food entries with nutrition data through simple tool calls.

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

Set the following environment variables:

```bash
export SPARKYFITNESS_API_URL=https://api.sparkyfitness.com
export SPARKYFITNESS_API_KEY=your-api-key-here
```

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

```bash
# Set environment variables first
export SPARKYFITNESS_API_URL=https://api.sparkyfitness.com
export SPARKYFITNESS_API_KEY=your-api-key

make run
# or
./sparkyfitness-mcp
```

## Docker

### Build Image

```bash
make docker-build
# or
docker build -t sparkyfitness-mcp .
```

### Run Container

```bash
docker run -e SPARKYFITNESS_API_URL=https://api.sparkyfitness.com \
           -e SPARKYFITNESS_API_KEY=your-api-key \
           sparkyfitness-mcp
```

## MCP Tools

### Available Tools

- `create_food_variant` - Create food variants with nutrition data
- `search_food_variants` - Search for existing food variants

