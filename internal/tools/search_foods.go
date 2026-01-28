package tools

import (
	"context"
	"fmt"

	"github.com/chickenzord/sparkyfitness-mcp/internal/sparkyfitness"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SearchFoodsInput defines the input parameters for the search_foods tool
type SearchFoodsInput struct {
	Name       string  `json:"name" jsonschema:"Food name to search for"`
	Brand      *string `json:"brand,omitempty" jsonschema:"Optional brand name to filter results"`
	BroadMatch *bool   `json:"broad_match,omitempty" jsonschema:"If true, performs broad match search (default: true)"`
	Limit      *int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default: 10)"`
}

// FoodResult represents a single food with its default variant in the search results
type FoodResult struct {
	FoodID           string   `json:"food_id" jsonschema:"Unique identifier of the food"`
	FoodName         string   `json:"food_name" jsonschema:"Name of the food"`
	Brand            *string  `json:"brand,omitempty" jsonschema:"Brand name if available"`
	IsCustom         bool     `json:"is_custom" jsonschema:"Whether this is a custom food"`
	ProviderType     *string  `json:"provider_type,omitempty" jsonschema:"Provider type (e.g., usda, nutritionix)"`
	VariantID        string   `json:"variant_id" jsonschema:"Unique identifier of the default variant"`
	ServingSize      float64  `json:"serving_size" jsonschema:"Serving size amount"`
	ServingUnit      string   `json:"serving_unit" jsonschema:"Unit of measurement for serving"`
	Calories         float64  `json:"calories" jsonschema:"Calories per serving"`
	Protein          float64  `json:"protein" jsonschema:"Protein in grams"`
	Carbs            float64  `json:"carbs" jsonschema:"Carbohydrates in grams"`
	Fat              float64  `json:"fat" jsonschema:"Fat in grams"`
	SaturatedFat     float64  `json:"saturated_fat,omitempty" jsonschema:"Saturated fat in grams"`
	DietaryFiber     float64  `json:"dietary_fiber,omitempty" jsonschema:"Dietary fiber in grams"`
	Sugars           float64  `json:"sugars,omitempty" jsonschema:"Sugars in grams"`
	Sodium           float64  `json:"sodium,omitempty" jsonschema:"Sodium in milligrams"`
	Cholesterol      float64  `json:"cholesterol,omitempty" jsonschema:"Cholesterol in milligrams"`
	Potassium        float64  `json:"potassium,omitempty" jsonschema:"Potassium in milligrams"`
	VitaminA         float64  `json:"vitamin_a,omitempty" jsonschema:"Vitamin A"`
	VitaminC         float64  `json:"vitamin_c,omitempty" jsonschema:"Vitamin C"`
	Calcium          float64  `json:"calcium,omitempty" jsonschema:"Calcium"`
	Iron             float64  `json:"iron,omitempty" jsonschema:"Iron"`
	GlycemicIndex    *string  `json:"glycemic_index,omitempty" jsonschema:"Glycemic index if available"`
}

// SearchFoodsOutput defines the output structure
type SearchFoodsOutput struct {
	Foods []FoodResult `json:"foods" jsonschema:"List of matching foods with their default variants"`
	Total int          `json:"total" jsonschema:"Total number of foods found"`
}

// RegisterSearchFoods registers the search_foods tool with the MCP server
func (r *Registry) RegisterSearchFoods(server *mcp.Server, client *sparkyfitness.Client) error {
	tool := &mcp.Tool{
		Name:        "search_foods",
		Description: "Search for foods in the SparkyFitness database by name and optional brand. Returns matching foods with their default nutrition information.",
	}

	handler := func(ctx context.Context, request *mcp.CallToolRequest, input SearchFoodsInput) (*mcp.CallToolResult, SearchFoodsOutput, error) {
		// Validate required parameters
		if input.Name == "" {
			return nil, SearchFoodsOutput{}, fmt.Errorf("name parameter is required")
		}

		// Set defaults
		broadMatch := true
		if input.BroadMatch != nil {
			broadMatch = *input.BroadMatch
		}

		limit := 10
		if input.Limit != nil {
			limit = *input.Limit
		}

		// Search for foods using backend API
		foods, err := client.SearchFoods(ctx, input.Name, broadMatch, limit)
		if err != nil {
			return nil, SearchFoodsOutput{}, fmt.Errorf("failed to search foods: %w", err)
		}

		// Filter by brand if specified
		var filteredFoods []sparkyfitness.Food
		if input.Brand != nil && *input.Brand != "" {
			for _, food := range foods {
				// Check brand (handle both nil and empty string)
				if food.Brand != nil && *food.Brand == *input.Brand {
					filteredFoods = append(filteredFoods, food)
				}
			}
		} else {
			filteredFoods = foods
		}

		// No foods found
		if len(filteredFoods) == 0 {
			return nil, SearchFoodsOutput{
				Foods: []FoodResult{},
				Total: 0,
			}, nil
		}

		// Convert foods to result format
		var results []FoodResult
		for _, food := range filteredFoods {
			// Skip foods without default variant
			if food.DefaultVariant == nil {
				continue
			}

			result := convertFoodToResult(food)
			results = append(results, result)
		}

		// Prepare output
		output := SearchFoodsOutput{
			Foods: results,
			Total: len(results),
		}

		return nil, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}

// convertFoodToResult converts a Food from backend API to FoodResult
func convertFoodToResult(food sparkyfitness.Food) FoodResult {
	variant := food.DefaultVariant

	return FoodResult{
		FoodID:        food.ID,
		FoodName:      food.Name,
		Brand:         food.Brand,
		IsCustom:      food.IsCustom,
		ProviderType:  food.ProviderType,
		VariantID:     variant.ID,
		ServingSize:   variant.ServingSize,
		ServingUnit:   variant.ServingUnit,
		Calories:      variant.Calories,
		Protein:       variant.Protein,
		Carbs:         variant.Carbs,
		Fat:           variant.Fat,
		SaturatedFat:  variant.SaturatedFat,
		DietaryFiber:  variant.DietaryFiber,
		Sugars:        variant.Sugars,
		Sodium:        variant.Sodium,
		Cholesterol:   variant.Cholesterol,
		Potassium:     variant.Potassium,
		VitaminA:      variant.VitaminA,
		VitaminC:      variant.VitaminC,
		Calcium:       variant.Calcium,
		Iron:          variant.Iron,
		GlycemicIndex: variant.GlycemicIndex,
	}
}
