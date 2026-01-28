package config

import (
	"fmt"
	"net/http"
	"os"
)

// TransportMode defines the transport protocol for the MCP server
type TransportMode string

const (
	// TransportStdio uses standard input/output for MCP communication
	TransportStdio TransportMode = "stdio"
	// TransportHTTP uses HTTP/SSE for MCP communication
	TransportHTTP TransportMode = "http"
)

// Config holds the application configuration
type Config struct {
	// SparkyFitnessAPIURL is the base URL for the SparkyFitness API
	SparkyFitnessAPIURL string
	// SparkyFitnessAPIKey is the authentication credential for the API
	SparkyFitnessAPIKey string
	// Transport defines the transport mode (stdio or http)
	Transport TransportMode
	// HTTPHost is the host to bind to when using HTTP transport
	HTTPHost string
	// HTTPPort is the port to listen on when using HTTP transport
	HTTPPort string
	// HTTPBasicAuthUser is the username for HTTP basic authentication (optional)
	HTTPBasicAuthUser string
	// HTTPBasicAuthPassword is the password for HTTP basic authentication (optional)
	HTTPBasicAuthPassword string
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

	// Transport mode (default: stdio)
	transport := TransportMode(os.Getenv("MCP_TRANSPORT"))
	if transport == "" {
		transport = TransportStdio
	}
	if transport != TransportStdio && transport != TransportHTTP {
		return nil, fmt.Errorf("invalid MCP_TRANSPORT value: %s (must be 'stdio' or 'http')", transport)
	}

	// HTTP configuration (only needed for HTTP transport)
	httpHost := os.Getenv("MCP_HTTP_HOST")
	if httpHost == "" {
		httpHost = "0.0.0.0"
	}

	httpPort := os.Getenv("MCP_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// HTTP Basic Auth (optional)
	httpBasicAuthUser := os.Getenv("MCP_HTTP_BASIC_AUTH_USER")
	httpBasicAuthPassword := os.Getenv("MCP_HTTP_BASIC_AUTH_PASSWORD")

	return &Config{
		SparkyFitnessAPIURL:   apiURL,
		SparkyFitnessAPIKey:   apiKey,
		Transport:             transport,
		HTTPHost:              httpHost,
		HTTPPort:              httpPort,
		HTTPBasicAuthUser:     httpBasicAuthUser,
		HTTPBasicAuthPassword: httpBasicAuthPassword,
	}, nil
}

// BasicAuthEnabled returns true if basic authentication is configured
func (c *Config) BasicAuthEnabled() bool {
	return c.HTTPBasicAuthUser != "" && c.HTTPBasicAuthPassword != ""
}

// BasicAuthMiddleware returns an HTTP middleware that enforces basic authentication
// If basic auth is not configured, it returns a no-op middleware
func (c *Config) BasicAuthMiddleware(next http.Handler) http.Handler {
	if !c.BasicAuthEnabled() {
		// No auth configured, pass through
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != c.HTTPBasicAuthUser || pass != c.HTTPBasicAuthPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="MCP Server"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
