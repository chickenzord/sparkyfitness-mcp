package tools

import (
	"context"
	"fmt"

	"github.com/chickenzord/sparkyfitness-mcp/internal/sparkyfitness"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddFoodVariantInput defines the input parameters for the add_food_variant tool
type AddFoodVariantInput struct {
	FoodID               string   `json:"food_id" jsonschema:"required,Unique identifier of the existing food (from search_foods results)"`
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
	IsDefault            *bool    `json:"is_default,omitempty" jsonschema:"Set this variant as the food's default variant (default: false)"`
	GlycemicIndex        *string  `json:"glycemic_index,omitempty" jsonschema:"Glycemic index if available"`
}

// AddFoodVariantOutput defines the output structure
type AddFoodVariantOutput struct {
	FoodID    string `json:"food_id" jsonschema:"ID of the food this variant was added to"`
	VariantID string `json:"variant_id" jsonschema:"ID of the newly created variant"`
	Message   string `json:"message" jsonschema:"Success message"`
}

// RegisterAddFoodVariant registers the add_food_variant tool with the MCP server
func (r *Registry) RegisterAddFoodVariant(server *mcp.Server, client *sparkyfitness.Client) error {
	tool := &mcp.Tool{
		Name: "add_food_variant",
		Description: "Add a new variant (serving size + nutrition) to an EXISTING food. " +
			"Use this when search_foods found a matching food and the user confirms they want to add a variant to it. " +
			"Requires food_id from search_foods results. " +
			"This adds another serving size option to an existing food entry rather than creating a new food.",
	}

	handler := func(ctx context.Context, request *mcp.CallToolRequest, input AddFoodVariantInput) (*mcp.CallToolResult, AddFoodVariantOutput, error) {
		// Validate required parameters
		if input.FoodID == "" {
			return nil, AddFoodVariantOutput{}, fmt.Errorf("food_id parameter is required")
		}
		if input.ServingSize <= 0 {
			return nil, AddFoodVariantOutput{}, fmt.Errorf("serving_size must be greater than 0")
		}
		if input.ServingUnit == "" {
			return nil, AddFoodVariantOutput{}, fmt.Errorf("serving_unit parameter is required")
		}

		// Build request for backend API
		req := &sparkyfitness.AddFoodVariantRequest{
			FoodID:      input.FoodID,
			ServingSize: input.ServingSize,
			ServingUnit: input.ServingUnit,
			Calories:    input.Calories,
			Protein:     input.Protein,
			Carbs:       input.Carbs,
			Fat:         input.Fat,
		}

		// Add optional fields
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
		if input.IsDefault != nil {
			req.IsDefault = *input.IsDefault
		}
		if input.GlycemicIndex != nil {
			req.GlycemicIndex = input.GlycemicIndex
		}

		// Call backend API to add variant
		resp, err := client.AddFoodVariant(ctx, req)
		if err != nil {
			return nil, AddFoodVariantOutput{}, fmt.Errorf("failed to add food variant: %w", err)
		}

		// Prepare output
		output := AddFoodVariantOutput{
			FoodID:    input.FoodID,
			VariantID: resp.ID,
			Message:   fmt.Sprintf("Successfully added variant to existing food (variant ID: %s)", resp.ID),
		}

		return nil, output, nil
	}

	mcp.AddTool(server, tool, handler)
	return nil
}
