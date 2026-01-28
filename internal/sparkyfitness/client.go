package sparkyfitness

import (
	"context"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
)

// SparkyClient wraps the generated SparkyFitness API client with authentication and configuration
type SparkyClient struct {
	client *ClientWithResponses
	config *config.Config
}

// authInterceptor adds authentication to outgoing requests
type authInterceptor struct {
	apiKey string
	next   http.RoundTripper
}

func (a *authInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add API key authentication via cookie
	// SparkyFitness uses cookie-based auth with "token" cookie
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: a.apiKey,
	})

	return a.next.RoundTrip(req)
}

// NewSparkyClient creates a new SparkyFitness API client with authentication
func NewSparkyClient(cfg *config.Config) (*SparkyClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if cfg.SparkyFitnessAPIURL == "" {
		return nil, fmt.Errorf("API URL is required")
	}

	if cfg.SparkyFitnessAPIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Create HTTP client with auth interceptor
	httpClient := &http.Client{
		Transport: &authInterceptor{
			apiKey: cfg.SparkyFitnessAPIKey,
			next:   http.DefaultTransport,
		},
	}

	// Create generated client with responses
	generatedClient, err := NewClientWithResponses(cfg.SparkyFitnessAPIURL, WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &SparkyClient{
		client: generatedClient,
		config: cfg,
	}, nil
}

// UpsertFood creates or retrieves a food by name and optional brand
// Returns the food ID
func (c *SparkyClient) UpsertFood(ctx context.Context, name string, brand *string) (openapi_types.UUID, error) {
	// Build food suggestion with name and brand in data
	foodData := map[string]interface{}{
		"name": name,
	}
	if brand != nil && *brand != "" {
		foodData["brand"] = *brand
	}

	body := PostFoodCrudCreateOrGetJSONRequestBody{
		FoodSuggestion: &Food{
			Name: name,
			Data: foodData,
		},
	}

	resp, err := c.client.PostFoodCrudCreateOrGetWithResponse(ctx, body)
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("failed to upsert food: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return openapi_types.UUID{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil || resp.JSON200.FoodId == nil {
		return openapi_types.UUID{}, fmt.Errorf("unexpected response format: missing food ID")
	}

	return *resp.JSON200.FoodId, nil
}

// CreateFoodVariant creates a new food variant with nutrition data
func (c *SparkyClient) CreateFoodVariant(ctx context.Context, variant FoodVariant) (*FoodVariant, error) {
	resp, err := c.client.PostFoodCrudFoodVariantsWithResponse(ctx, variant)
	if err != nil {
		return nil, fmt.Errorf("failed to create food variant: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON201 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}

	return resp.JSON201, nil
}

// SearchFoods searches for foods by name
func (c *SparkyClient) SearchFoods(ctx context.Context, name string, exactMatch bool) ([]Food, error) {
	params := &GetFoodCrudSearchParams{
		Name:       name,
		ExactMatch: &exactMatch,
	}

	resp, err := c.client.GetFoodCrudSearchWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search foods: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}

	return *resp.JSON200, nil
}

// ListFoodVariants lists all variants for a given food
func (c *SparkyClient) ListFoodVariants(ctx context.Context, foodID openapi_types.UUID) ([]FoodVariant, error) {
	params := &GetFoodCrudFoodVariantsParams{
		FoodId: foodID,
	}

	resp, err := c.client.GetFoodCrudFoodVariantsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list food variants: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}

	return *resp.JSON200, nil
}

// GetUnderlyingClient returns the underlying generated client for advanced usage
func (c *SparkyClient) GetUnderlyingClient() *ClientWithResponses {
	return c.client
}
