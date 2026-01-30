package sparkyfitness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
)

// Client is a manual HTTP client for the SparkyFitness backend API
type Client struct {
	httpClient *http.Client
	baseURL    string
	config     *config.Config
}

// authInterceptor adds authentication to outgoing requests
type authInterceptor struct {
	apiKey string
	next   http.RoundTripper
}

func (a *authInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add API key authentication via Bearer token
	// Cookie-based auth is for frontend UI, API clients use Bearer tokens
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))

	return a.next.RoundTrip(req)
}

// NewClient creates a new SparkyFitness API client with authentication
func NewClient(cfg *config.Config) (*Client, error) {
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

	return &Client{
		httpClient: httpClient,
		baseURL:    cfg.SparkyFitnessAPIURL,
		config:     cfg,
	}, nil
}

// SearchFoods searches for foods by name using the backend API
// Note: broadMatch and exactMatch are mutually exclusive. If broadMatch=false, exactMatch will be used.
func (c *Client) SearchFoods(ctx context.Context, name string, broadMatch bool, limit int) ([]Food, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("name", name)

	// broadMatch and exactMatch are mutually exclusive - one must be true
	if broadMatch {
		params.Set("broadMatch", "true")
	} else {
		params.Set("exactMatch", "true")
	}

	params.Set("limit", strconv.Itoa(limit))

	// Build request URL
	reqURL := fmt.Sprintf("%s/foods?%s", c.baseURL, params.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request (auth interceptor will add Bearer token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var searchResp SearchFoodsResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return searchResp.SearchResults, nil
}

// AddFoodVariant adds a new variant to an existing food
// Backend endpoint: POST /foods/food-variants
// Returns 201 Created with variant ID
func (c *Client) AddFoodVariant(ctx context.Context, req *AddFoodVariantRequest) (*AddFoodVariantResponse, error) {
	// Build request URL
	reqURL := fmt.Sprintf("%s/foods/food-variants", c.baseURL)

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request (auth interceptor will add Bearer token)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code (backend returns 201 Created)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var addVariantResp AddFoodVariantResponse
	if err := json.Unmarshal(respBody, &addVariantResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &addVariantResp, nil
}
