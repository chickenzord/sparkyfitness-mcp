package sparkyfitness

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
)

func TestNewSparkyClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
			errMsg:  "config cannot be nil",
		},
		{
			name: "empty API URL",
			cfg: &config.Config{
				SparkyFitnessAPIURL: "",
				SparkyFitnessAPIKey: "test-key",
			},
			wantErr: true,
			errMsg:  "API URL is required",
		},
		{
			name: "empty API key",
			cfg: &config.Config{
				SparkyFitnessAPIURL: "https://api.example.com",
				SparkyFitnessAPIKey: "",
			},
			wantErr: true,
			errMsg:  "API key is required",
		},
		{
			name: "valid config",
			cfg: &config.Config{
				SparkyFitnessAPIURL: "https://api.example.com",
				SparkyFitnessAPIKey: "test-key",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSparkyClient(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSparkyClient() expected error but got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("NewSparkyClient() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("NewSparkyClient() unexpected error = %v", err)
				}
				if client == nil {
					t.Errorf("NewSparkyClient() returned nil client")
				}
			}
		})
	}
}

func TestSparkyClient_UpsertFood(t *testing.T) {
	// Create test server
	testFoodID := openapi_types.UUID{}
	_ = testFoodID.UnmarshalText([]byte("550e8400-e29b-41d4-a716-446655440000"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth cookie
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value != "test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify endpoint
		if r.URL.Path != "/api/food-crud/create-or-get" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Decode request body
		var reqBody PostFoodCrudCreateOrGetJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify request
		if reqBody.FoodSuggestion == nil {
			t.Error("foodSuggestion is nil")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"foodId": testFoodID.String(),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		SparkyFitnessAPIURL: server.URL + "/api",
		SparkyFitnessAPIKey: "test-key",
	}

	client, err := NewSparkyClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test upsert food
	brand := "Test Brand"
	foodID, err := client.UpsertFood(context.Background(), "Test Food", &brand)
	if err != nil {
		t.Fatalf("UpsertFood() error = %v", err)
	}

	if foodID.String() != testFoodID.String() {
		t.Errorf("UpsertFood() foodID = %v, want %v", foodID, testFoodID)
	}
}

func TestSparkyClient_CreateFoodVariant(t *testing.T) {
	testFoodID := openapi_types.UUID{}
	_ = testFoodID.UnmarshalText([]byte("550e8400-e29b-41d4-a716-446655440000"))

	testVariantID := openapi_types.UUID{}
	_ = testVariantID.UnmarshalText([]byte("660e8400-e29b-41d4-a716-446655440000"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth cookie
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value != "test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify endpoint
		if r.URL.Path != "/api/food-crud/food-variants" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Decode request
		var variant FoodVariant
		if err := json.NewDecoder(r.Body).Decode(&variant); err != nil {
			t.Errorf("failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify variant
		if variant.ServingSize == "" {
			t.Error("serving_size is empty")
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		variant.Id = testVariantID
		json.NewEncoder(w).Encode(variant)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		SparkyFitnessAPIURL: server.URL + "/api",
		SparkyFitnessAPIKey: "test-key",
	}

	client, err := NewSparkyClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test create food variant
	variant := FoodVariant{
		FoodId:        testFoodID,
		ServingSize:   "100g",
		ServingWeight: 100,
		Data: map[string]interface{}{
			"calories": 250.0,
			"protein":  20.0,
			"carbs":    30.0,
			"fat":      10.0,
		},
	}

	result, err := client.CreateFoodVariant(context.Background(), variant)
	if err != nil {
		t.Fatalf("CreateFoodVariant() error = %v", err)
	}

	if result.Id.String() != testVariantID.String() {
		t.Errorf("CreateFoodVariant() variant ID = %v, want %v", result.Id, testVariantID)
	}

	if result.ServingSize != "100g" {
		t.Errorf("CreateFoodVariant() serving_size = %v, want 100g", result.ServingSize)
	}
}

func TestSparkyClient_SearchFoods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth cookie
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value != "test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify endpoint
		if r.URL.Path != "/api/food-crud/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify query params
		name := r.URL.Query().Get("name")
		if name == "" {
			t.Error("name query param is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		foods := []Food{
			{
				Name: name,
				Data: map[string]interface{}{},
			},
		}
		json.NewEncoder(w).Encode(foods)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		SparkyFitnessAPIURL: server.URL + "/api",
		SparkyFitnessAPIKey: "test-key",
	}

	client, err := NewSparkyClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test search foods
	foods, err := client.SearchFoods(context.Background(), "chicken", false)
	if err != nil {
		t.Fatalf("SearchFoods() error = %v", err)
	}

	if len(foods) != 1 {
		t.Errorf("SearchFoods() returned %d foods, want 1", len(foods))
	}

	if foods[0].Name != "chicken" {
		t.Errorf("SearchFoods() food name = %v, want chicken", foods[0].Name)
	}
}

func TestAuthInterceptor(t *testing.T) {
	// Test that auth interceptor adds the correct cookie
	var receivedCookie *http.Cookie

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			t.Errorf("failed to get cookie: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		receivedCookie = cookie
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foodId":"550e8400-e29b-41d4-a716-446655440000"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		SparkyFitnessAPIURL: server.URL,
		SparkyFitnessAPIKey: "secret-token-123",
	}

	client, err := NewSparkyClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Make a request that will trigger the auth interceptor
	brand := "Test"
	_, err = client.UpsertFood(context.Background(), "Test Food", &brand)
	if err != nil {
		t.Fatalf("UpsertFood() error = %v", err)
	}

	// Verify cookie was set correctly
	if receivedCookie == nil {
		t.Fatal("no cookie received")
	}

	if receivedCookie.Name != "token" {
		t.Errorf("cookie name = %v, want token", receivedCookie.Name)
	}

	if receivedCookie.Value != "secret-token-123" {
		t.Errorf("cookie value = %v, want secret-token-123", receivedCookie.Value)
	}
}
