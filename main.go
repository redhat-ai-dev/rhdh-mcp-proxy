// Assisted-by: claude-4-sonnet

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Config holds the proxy configuration
type Config struct {
	BackstageURL string
	MCPToken     string
	Port         string
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	backstageURL := os.Getenv("BACKSTAGE_URL")
	if backstageURL == "" {
		return nil, fmt.Errorf("BACKSTAGE_URL environment variable is required")
	}

	mcpToken := os.Getenv("MCP_TOKEN")
	if mcpToken == "" {
		return nil, fmt.Errorf("MCP_TOKEN environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Ensure BACKSTAGE_URL doesn't end with a slash
	backstageURL = strings.TrimSuffix(backstageURL, "/")

	return &Config{
		BackstageURL: backstageURL,
		MCPToken:     mcpToken,
		Port:         port,
	}, nil
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Parse the target URL
	targetURL, err := url.Parse(config.BackstageURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the proxy to add authentication
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Add Bearer token authentication
		req.Header.Set("Authorization", "Bearer "+config.MCPToken)
		// Log the request
		log.Printf("Proxying %s %s", req.Method, req.URL.Path)
	}

	// Create HTTP server with the new ServeMux from Go 1.22
	mux := http.NewServeMux()

	// Handle only /api/mcp-actions requests with the proxy
	mux.HandleFunc("/api/mcp-actions/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	// Handle all other requests with a 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Printf("Starting MCP Proxy Server on port %s", config.Port)
	log.Printf("Proxying requests to: %s", config.BackstageURL)

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
