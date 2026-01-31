package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chickenzord/sparkyfitness-mcp/internal/sparkyfitness"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CreateFoodInput defines the input parameters for the create_food_variant tool
type CreateFoodInput struct {
	Name                 string   `json:"name" jsonschema:"required,Food name"`
	Brand                *string  `json:"brand,omitempty" jsonschema:"Brand name (optional)"`
	ServingSize          float64  `json:"serving_size" jsonschema:"required,Numeric serving size amount (e.g., 100, 1)"`
	ServingUnit          string   `json:"serving_unit" jsonschema:"required,Unit of measurement (e.g., g, ml, cup, piece)"`
	Calories             float64  `json:"calories" jsonschema:"required,Calories per serving"`
	Protein              float64  `json:"protein" jsonschema:"required,Protein in grams"`
	Carbs                float64  `json:"carbs" jsonschema:"required,Carbohydrates in grams"`
	Fat                  float64  `json:"fat" jsonschema:"required,Fat in grams"`
	SaturatedFat         *float64 `json:"saturated_fat,omitempty" jsonschema:"Saturated fat in grams"`
	PolyunsaturatedFat   *float64 `json:"polyunsaturated_fat,omitempty" jsonschema:"Polyunsaturated fat in grams"`
	MonounsaturatedFat   *float64 `json:"monounsaturated_fat,omitempty" jsonschema:"Monounsaturated fat in grams"`
	TransFat             *float64 `json:"trans_fat,omitempty" jsonschema:"Trans fat in grams"`
	Cholesterol          *float64 `json:"cholesterol,omitempty" jsonschema:"Cholesterol in milligrams"`
	Sodium               *float64 `json:"sodium,omitempty" jsonschema:"Sodium in milligrams"`
	Potassium            *float64 `json:"potassium,omitempty" jsonschema:"Potassium in milligrams"`
	DietaryFiber         *float64 `json:"dietary_fiber,omitempty" jsonschema:"Dietary fiber in grams"`
	Sugars               *float64 `json:"sugars,omitempty" jsonschema:"Sugars in grams"`
	VitaminA             *float64 `json:"vitamin_a,omitempty" jsonschema:"Vitamin A"`
	VitaminC             *float64 `json:"vitamin_c,omitempty" jsonschema:"Vitamin C"`
	Calcium              *float64 `json:"calcium,omitempty" jsonschema:"Calcium"`
	Iron                 *float64 `json:"iron,omitempty" jsonschema:"Iron"`
	IsQuickFood          *bool    `json:"is_quick_food,omitempty" jsonschema:"Mark as quick food (default: false)"`
	IsDefault            *bool    `json:"is_default,omitempty" jsonschema:"Set this variant as default (default: true for first variant)"`
	GlycemicIndex        *string  `json:"glycemic_index,omitempty" jsonschema:"Glycemic index if available"`
}

// CreateFoodOutput defines the output structure
type CreateFoodOutput struct {
	FoodID    string `json:"food_id" jsonschema:"ID of the created food"`
	VariantID string `json:"variant_id" jsonschema:"ID of the created variant"`
	Message   string `json:"message" jsonschema:"Success message"`
}

// CreateFoodRequest represents the backend API request structure for POST /foods
type CreateFoodRequest struct {
	Name                 string                 `json:"name"`
	Brand                string                 `json:"brand"`
	IsCustom             bool                   `json:"is_custom"`
	IsQuickFood          bool                   `json:"is_quick_food"`
	ServingSize          float64                `json:"serving_size"`
	ServingUnit          string                 `json:"serving_unit"`
	Calories             float64                `json:"calories"`
	Protein              float64                `json:"protein"`
	Carbs                float64                `json:"carbs"`
	Fat                  float64                `json:"fat"`
	SaturatedFat         float64                `json:"saturated_fat"`
	PolyunsaturatedFat   float64                `json:"polyunsaturated_fat"`
	MonounsaturatedFat   float64                `json:"monounsaturated_fat"`
	TransFat             float64                `json:"trans_fat"`
	Cholesterol          float64                `json:"cholesterol"`
	Sodium               float64                `json:"sodium"`
	Potassium            float64                `json:"potassium"`
	DietaryFiber         float64                `json:"dietary_fiber"`
	Sugars               float64                `json:"sugars"`
	VitaminA             float64                `json:"vitamin_a"`
	VitaminC             float64                `json:"vitamin_c"`
	Calcium              float64                `json:"calcium"`
	Iron                 float64                `json:"iron"`
	IsDefault            bool                   `json:"is_default"`
	GlycemicIndex        string                 `json:"glycemic_index"`
	CustomNutrients      map[string]interface{} `json:"custom_nutrients"`
}

// CreateFoodResponse represents the backend API response for POST /foods
type CreateFoodResponse struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Brand          string                  `json:"brand"`
	IsCustom       bool                    `json:"is_custom"`
	UserID         string                  `json:"user_id"`
	DefaultVariant *CreateFoodVariantNested `json:"default_variant"`
}

// CreateFoodVariantNested represents the nested variant in create food response
type CreateFoodVariantNested struct {
	ID string `json:"id"`
}

// RegisterCreateFoodVariant registers the create_food_variant tool with the MCP server
func (r *Registry) RegisterCreateFoodVariant(server *mcp.Server, client *sparkyfitness.Client) error {
	tool := &mcp.Tool{
		Name:  "create_food_variant",
		Title: "Create New Food Entry",
		Description: "🆕 Create a NEW food entry with default variant in SparkyFitness.\n\n" +
			"**When to Use:**\n" +
			"ONLY use this when:\n" +
			"1. search_foods found NO matches (no duplicates exist), OR\n" +
			"2. User explicitly chooses to create a separate food entry despite duplicates\n\n" +
			"**What This Does:**\n" +
			"Creates a completely new food entity in the database with its first serving size variant. This is NOT for adding variants to existing foods - use add_food_variant for that.\n\n" +
			"**Required Input:**\n" +
			"• name: Food name (e.g., 'Organic Quinoa')\n" +
			"• brand: Brand name (optional, e.g., 'Nature's Best')\n" +
			"• serving_size: Numeric amount (e.g., 100, 1)\n" +
			"• serving_unit: Unit of measurement (g, ml, cup, piece, oz, etc.)\n" +
			"• Core nutrition: calories, protein, carbs, fat (all required)\n" +
			"• Optional nutrition: fiber, sugar, vitamins, minerals, etc.\n\n" +
			"**Output:**\n" +
			"• food_id: UUID of the newly created food\n" +
			"• variant_id: UUID of the default variant\n" +
			"• Success message\n\n" +
			"**Example Workflow:**\n" +
			"User: 'Add nutrition for Organic Quinoa by Nature's Best'\n" +
			"1. search_foods(name='Organic Quinoa', brand='Nature's Best') → no matches\n" +
			"2. create_food_variant(name='Organic Quinoa', brand='Nature's Best', ...nutrition)\n" +
			"3. Result: New food created in database\n\n" +
			"**Important:**\n" +
			"If search_foods found existing foods, show them to the user first and ask whether to:\n" +
			"• Add variant to existing food (use add_food_variant), OR\n" +
			"• Create new separate entry (use this tool)",
	}

	handler := func(ctx context.Context, request *mcp.CallToolRequest, input CreateFoodInput) (*mcp.CallToolResult, CreateFoodOutput, error) {
		// Validate required parameters
		if input.Name == "" {
			return nil, CreateFoodOutput{}, fmt.Errorf("name parameter is required")
		}
		if input.ServingSize <= 0 {
			return nil, CreateFoodOutput{}, fmt.Errorf("serving_size must be greater than 0")
		}
		if input.ServingUnit == "" {
			return nil, CreateFoodOutput{}, fmt.Errorf("serving_unit parameter is required")
		}

		// Build request for backend API
		req := &CreateFoodRequest{
			Name:            input.Name,
			Brand:           "",
			IsCustom:        true, // MCP-created foods are always custom
			IsQuickFood:     false,
			ServingSize:     input.ServingSize,
			ServingUnit:     input.ServingUnit,
			Calories:        input.Calories,
			Protein:         input.Protein,
			Carbs:           input.Carbs,
			Fat:             input.Fat,
			IsDefault:       true, // First variant is always default
			GlycemicIndex:   "None",
			CustomNutrients: make(map[string]interface{}),
		}

		// Set optional brand
		if input.Brand != nil {
			req.Brand = *input.Brand
		}

		// Set optional nutrition fields
		if input.SaturatedFat != nil {
			req.SaturatedFat = *input.SaturatedFat
		}
		if input.PolyunsaturatedFat != nil {
			req.PolyunsaturatedFat = *input.PolyunsaturatedFat
		}
		if input.MonounsaturatedFat != nil {
			req.MonounsaturatedFat = *input.MonounsaturatedFat
		}
		if input.TransFat != nil {
			req.TransFat = *input.TransFat
		}
		if input.Cholesterol != nil {
			req.Cholesterol = *input.Cholesterol
		}
		if input.Sodium != nil {
			req.Sodium = *input.Sodium
		}
		if input.Potassium != nil {
			req.Potassium = *input.Potassium
		}
		if input.DietaryFiber != nil {
			req.DietaryFiber = *input.DietaryFiber
		}
		if input.Sugars != nil {
			req.Sugars = *input.Sugars
		}
		if input.VitaminA != nil {
			req.VitaminA = *input.VitaminA
		}
		if input.VitaminC != nil {
			req.VitaminC = *input.VitaminC
		}
		if input.Calcium != nil {
			req.Calcium = *input.Calcium
		}
		if input.Iron != nil {
			req.Iron = *input.Iron
		}
		if input.IsQuickFood != nil {
			req.IsQuickFood = *input.IsQuickFood
		}
		if input.GlycemicIndex != nil {
			req.GlycemicIndex = *input.GlycemicIndex
		}

		// Call backend API to create food + variant
		resp, err := createFood(ctx, client, req)
		if err != nil {
			return nil, CreateFoodOutput{}, fmt.Errorf("failed to create food: %w", err)
		}

		// Prepare output
		foodName := resp.Name
		if resp.Brand != "" {
			foodName = fmt.Sprintf("%s (%s)", resp.Name, resp.Brand)
		}

		output := CreateFoodOutput{
			FoodID:    resp.ID,
			VariantID: resp.DefaultVariant.ID,
			Message:   fmt.Sprintf("Successfully created new food '%s' with default variant", foodName),
		}

		return nil, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// createFood calls the backend API to create a food with variant
// Backend endpoint: POST /foods
func createFood(ctx context.Context, client *sparkyfitness.Client, req *CreateFoodRequest) (*CreateFoodResponse, error) {
	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build request URL - using reflection to get baseURL from client
	// For now, we'll need to add a helper method to the client
	baseURL := client.BaseURL()
	reqURL := fmt.Sprintf("%s/foods", baseURL)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request (auth interceptor will add Bearer token)
	httpClient := client.HTTPClient()
	resp, err := httpClient.Do(httpReq)
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
	var createResp CreateFoodResponse
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &createResp, nil
}
