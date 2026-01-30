package sparkyfitness

// Backend API types - manually defined based on real API structure

// Food represents a food item from the backend API
type Food struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	Brand              *string       `json:"brand"`
	IsCustom           bool          `json:"is_custom"`
	UserID             string        `json:"user_id"`
	SharedWithPublic   bool          `json:"shared_with_public"`
	ProviderExternalID *string       `json:"provider_external_id"`
	ProviderType       *string       `json:"provider_type"`
	DefaultVariant     *FoodVariant  `json:"default_variant"`
}

// FoodVariant represents a food variant with nutrition information
type FoodVariant struct {
	ID                   string                 `json:"id"`
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
	GlycemicIndex        *string                `json:"glycemic_index"`
	CustomNutrients      map[string]interface{} `json:"custom_nutrients"`
}

// SearchFoodsResponse represents the response from the search foods endpoint
type SearchFoodsResponse struct {
	SearchResults []Food `json:"searchResults"`
}

// AddFoodVariantRequest represents the request to add a variant to an existing food
// Backend endpoint: POST /foods/food-variants
type AddFoodVariantRequest struct {
	FoodID               string                 `json:"food_id"`
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
	GlycemicIndex        *string                `json:"glycemic_index,omitempty"`
	CustomNutrients      map[string]interface{} `json:"custom_nutrients,omitempty"`
}

// AddFoodVariantResponse represents the response from adding a food variant
// Backend returns 201 with just the variant ID
type AddFoodVariantResponse struct {
	ID string `json:"id"`
}
