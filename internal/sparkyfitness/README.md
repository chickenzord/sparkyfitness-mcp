# SparkyFitness API Client

This directory contains the auto-generated API client for SparkyFitness.

## Code Generation

The client code in `sparkyfitness_client.go` is **generated from OpenAPI spec** and should not be edited manually.

To regenerate the client:

```bash
make generate
```

This runs `oapi-codegen` with the configuration in `oapi-codegen.yaml` against `swagger.json`.

## Manual Wrapper

Any custom wrapper code or convenience methods should be added in separate files (e.g., `client.go`, `wrapper.go`) to avoid being overwritten during regeneration.
