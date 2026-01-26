package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	// SparkyFitnessAPIURL is the base URL for the SparkyFitness API
	SparkyFitnessAPIURL string
	// SparkyFitnessAPIKey is the authentication credential for the API
	SparkyFitnessAPIKey string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	apiURL := os.Getenv("SPARKYFITNESS_API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("SPARKYFITNESS_API_URL environment variable is required")
	}

	apiKey := os.Getenv("SPARKYFITNESS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SPARKYFITNESS_API_KEY environment variable is required")
	}

	return &Config{
		SparkyFitnessAPIURL: apiURL,
		SparkyFitnessAPIKey: apiKey,
	}, nil
}
