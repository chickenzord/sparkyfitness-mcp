package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		wantErr     bool
		errContains string
		wantConfig  *Config
	}{
		{
			name: "valid stdio config (default)",
			env: map[string]string{
				"SPARKYFITNESS_API_URL": "https://api.sparkyfitness.com",
				"SPARKYFITNESS_API_KEY": "test-key-123",
			},
			wantErr: false,
			wantConfig: &Config{
				SparkyFitnessAPIURL: "https://api.sparkyfitness.com",
				SparkyFitnessAPIKey: "test-key-123",
				Transport:           TransportStdio,
				HTTPHost:            "0.0.0.0",
				HTTPPort:            "8080",
			},
		},
		{
			name: "valid http config with custom host and port",
			env: map[string]string{
				"SPARKYFITNESS_API_URL": "https://api.sparkyfitness.com",
				"SPARKYFITNESS_API_KEY": "test-key-123",
				"MCP_TRANSPORT":         "http",
				"MCP_HTTP_HOST":         "127.0.0.1",
				"MCP_HTTP_PORT":         "9090",
			},
			wantErr: false,
			wantConfig: &Config{
				SparkyFitnessAPIURL: "https://api.sparkyfitness.com",
				SparkyFitnessAPIKey: "test-key-123",
				Transport:           TransportHTTP,
				HTTPHost:            "127.0.0.1",
				HTTPPort:            "9090",
			},
		},
		{
			name: "valid http config with defaults",
			env: map[string]string{
				"SPARKYFITNESS_API_URL": "https://api.sparkyfitness.com",
				"SPARKYFITNESS_API_KEY": "test-key-123",
				"MCP_TRANSPORT":         "http",
			},
			wantErr: false,
			wantConfig: &Config{
				SparkyFitnessAPIURL: "https://api.sparkyfitness.com",
				SparkyFitnessAPIKey: "test-key-123",
				Transport:           TransportHTTP,
				HTTPHost:            "0.0.0.0",
				HTTPPort:            "8080",
			},
		},
		{
			name: "missing API URL",
			env: map[string]string{
				"SPARKYFITNESS_API_KEY": "test-key",
			},
			wantErr:     true,
			errContains: "SPARKYFITNESS_API_URL",
		},
		{
			name: "missing API key",
			env: map[string]string{
				"SPARKYFITNESS_API_URL": "https://api.sparkyfitness.com",
			},
			wantErr:     true,
			errContains: "SPARKYFITNESS_API_KEY",
		},
		{
			name: "invalid transport mode",
			env: map[string]string{
				"SPARKYFITNESS_API_URL": "https://api.sparkyfitness.com",
				"SPARKYFITNESS_API_KEY": "test-key",
				"MCP_TRANSPORT":         "invalid",
			},
			wantErr:     true,
			errContains: "invalid MCP_TRANSPORT",
		},
		{
			name: "http config with basic auth",
			env: map[string]string{
				"SPARKYFITNESS_API_URL":      "https://api.sparkyfitness.com",
				"SPARKYFITNESS_API_KEY":      "test-key",
				"MCP_TRANSPORT":              "http",
				"MCP_HTTP_BASIC_AUTH_USER":   "admin",
				"MCP_HTTP_BASIC_AUTH_PASSWORD": "secret",
			},
			wantErr: false,
			wantConfig: &Config{
				SparkyFitnessAPIURL:   "https://api.sparkyfitness.com",
				SparkyFitnessAPIKey:   "test-key",
				Transport:             TransportHTTP,
				HTTPHost:              "0.0.0.0",
				HTTPPort:              "8080",
				HTTPBasicAuthUser:     "admin",
				HTTPBasicAuthPassword: "secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all relevant environment variables
			os.Unsetenv("SPARKYFITNESS_API_URL")
			os.Unsetenv("SPARKYFITNESS_API_KEY")
			os.Unsetenv("MCP_TRANSPORT")
			os.Unsetenv("MCP_HTTP_HOST")
			os.Unsetenv("MCP_HTTP_PORT")
			os.Unsetenv("MCP_HTTP_BASIC_AUTH_USER")
			os.Unsetenv("MCP_HTTP_BASIC_AUTH_PASSWORD")

			// Set test environment variables
			for k, v := range tt.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Test LoadFromEnv
			cfg, err := LoadFromEnv()

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadFromEnv() expected error, got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("LoadFromEnv() error = %v, should contain %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("LoadFromEnv() unexpected error: %v", err)
				return
			}

			if cfg.SparkyFitnessAPIURL != tt.wantConfig.SparkyFitnessAPIURL {
				t.Errorf("API URL = %v, want %v", cfg.SparkyFitnessAPIURL, tt.wantConfig.SparkyFitnessAPIURL)
			}

			if cfg.SparkyFitnessAPIKey != tt.wantConfig.SparkyFitnessAPIKey {
				t.Errorf("API Key = %v, want %v", cfg.SparkyFitnessAPIKey, tt.wantConfig.SparkyFitnessAPIKey)
			}

			if cfg.Transport != tt.wantConfig.Transport {
				t.Errorf("Transport = %v, want %v", cfg.Transport, tt.wantConfig.Transport)
			}

			if cfg.HTTPHost != tt.wantConfig.HTTPHost {
				t.Errorf("HTTPHost = %v, want %v", cfg.HTTPHost, tt.wantConfig.HTTPHost)
			}

			if cfg.HTTPPort != tt.wantConfig.HTTPPort {
				t.Errorf("HTTPPort = %v, want %v", cfg.HTTPPort, tt.wantConfig.HTTPPort)
			}

			if cfg.HTTPBasicAuthUser != tt.wantConfig.HTTPBasicAuthUser {
				t.Errorf("HTTPBasicAuthUser = %v, want %v", cfg.HTTPBasicAuthUser, tt.wantConfig.HTTPBasicAuthUser)
			}

			if cfg.HTTPBasicAuthPassword != tt.wantConfig.HTTPBasicAuthPassword {
				t.Errorf("HTTPBasicAuthPassword = %v, want %v", cfg.HTTPBasicAuthPassword, tt.wantConfig.HTTPBasicAuthPassword)
			}
		})
	}
}

func TestBasicAuthEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "both user and password set",
			config: &Config{
				HTTPBasicAuthUser:     "admin",
				HTTPBasicAuthPassword: "secret",
			},
			expected: true,
		},
		{
			name: "only user set",
			config: &Config{
				HTTPBasicAuthUser:     "admin",
				HTTPBasicAuthPassword: "",
			},
			expected: false,
		},
		{
			name: "only password set",
			config: &Config{
				HTTPBasicAuthUser:     "",
				HTTPBasicAuthPassword: "secret",
			},
			expected: false,
		},
		{
			name: "neither set",
			config: &Config{
				HTTPBasicAuthUser:     "",
				HTTPBasicAuthPassword: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.BasicAuthEnabled()
			if got != tt.expected {
				t.Errorf("BasicAuthEnabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}
