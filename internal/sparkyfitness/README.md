# SparkyFitness API Client

This directory contains the auto-generated API client for SparkyFitness with a convenient wrapper.

## Code Generation

The client code in `sparkyfitness_client.go` is **generated from OpenAPI spec** and should not be edited manually.

To regenerate the client:

```bash
make generate
```

This runs `oapi-codegen` with the configuration in `oapi-codegen.yaml` against `swagger.json`.

## Wrapper Client Usage

The `SparkyClient` wrapper in `client.go` provides a convenient interface to the generated client with automatic authentication.

### Creating a Client

```go
import (
    "github.com/chickenzord/sparkyfitness-mcp/internal/config"
    "github.com/chickenzord/sparkyfitness-mcp/internal/sparkyfitness"
)

cfg := &config.Config{
    SparkyFitnessAPIURL: "https://api.sparkyfitness.example.com",
    SparkyFitnessAPIKey: "your-api-key",
}

client, err := sparkyfitness.NewSparkyClient(cfg)
if err != nil {
    // Handle error
}
```

### Available Methods

- **UpsertFood**: Create or retrieve a food by name and brand
  ```go
  brand := "Acme Foods"
  foodID, err := client.UpsertFood(ctx, "Chicken Breast", &brand)
  ```

- **CreateFoodVariant**: Create a new food variant with nutrition data
  ```go
  variant := sparkyfitness.FoodVariant{
      FoodId:        foodID,
      ServingSize:   "100g",
      ServingWeight: 100,
      Data: map[string]interface{}{
          "calories": 165.0,
          "protein":  31.0,
          "carbs":    0.0,
          "fat":      3.6,
      },
  }
  result, err := client.CreateFoodVariant(ctx, variant)
  ```

- **SearchFoods**: Search for foods by name
  ```go
  foods, err := client.SearchFoods(ctx, "chicken", false)
  ```

- **ListFoodVariants**: List all variants for a food
  ```go
  variants, err := client.ListFoodVariants(ctx, foodID)
  ```

### Authentication

The wrapper automatically handles authentication using cookie-based auth (token cookie) as required by the SparkyFitness API.

### Testing

Tests are in `client_test.go` and use `httptest` to mock the API server. Run tests with:

```bash
go test ./internal/sparkyfitness/...
```

## Manual Wrapper

Any custom wrapper code or convenience methods should be added in separate files (e.g., `client.go`, `wrapper.go`) to avoid being overwritten during regeneration.
