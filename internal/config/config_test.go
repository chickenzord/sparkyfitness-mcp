package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		apiURL      string
		apiKey      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid config",
			apiURL:  "https://api.sparkyfitness.com",
			apiKey:  "test-key-123",
			wantErr: false,
		},
		{
			name:        "missing API URL",
			apiURL:      "",
			apiKey:      "test-key",
			wantErr:     true,
			errContains: "SPARKYFITNESS_API_URL",
		},
		{
			name:        "missing API key",
			apiURL:      "https://api.sparkyfitness.com",
			apiKey:      "",
			wantErr:     true,
			errContains: "SPARKYFITNESS_API_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.apiURL != "" {
				os.Setenv("SPARKYFITNESS_API_URL", tt.apiURL)
				defer os.Unsetenv("SPARKYFITNESS_API_URL")
			} else {
				os.Unsetenv("SPARKYFITNESS_API_URL")
			}

			if tt.apiKey != "" {
				os.Setenv("SPARKYFITNESS_API_KEY", tt.apiKey)
				defer os.Unsetenv("SPARKYFITNESS_API_KEY")
			} else {
				os.Unsetenv("SPARKYFITNESS_API_KEY")
			}

			// Test LoadFromEnv
			cfg, err := LoadFromEnv()

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadFromEnv() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadFromEnv() unexpected error: %v", err)
				return
			}

			if cfg.SparkyFitnessAPIURL != tt.apiURL {
				t.Errorf("LoadFromEnv() API URL = %v, want %v", cfg.SparkyFitnessAPIURL, tt.apiURL)
			}

			if cfg.SparkyFitnessAPIKey != tt.apiKey {
				t.Errorf("LoadFromEnv() API Key = %v, want %v", cfg.SparkyFitnessAPIKey, tt.apiKey)
			}
		})
	}
}
