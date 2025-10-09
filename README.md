# RHDH MCP Proxy

A simple Go-based HTTP proxy service that forwards requests to MCP servers running on a Backstage instance with authentication support.

## Features

- **HTTP Proxy**: Forwards all HTTP methods (GET, POST, PUT, DELETE, etc.) to the target MCP server
- **Authentication**: Automatically adds Bearer token authentication from environment variables
- **Streaming Support**: Handles Server-Sent Events (SSE) and other streaming responses
- **Header Forwarding**: Properly forwards request headers while filtering hop-by-hop headers
- **Error Handling**: Comprehensive error handling with appropriate HTTP status codes
- **Logging**: Request/response logging for debugging and monitoring

## Configuration

The proxy requires two environment variables:

- `BACKSTAGE_URL`: The base URL of your Backstage instance (e.g., `https://backstage.example.com`)
- `MCP_TOKEN`: The Bearer token for authenticating with the Backstage MCP server
- `PORT`: (Optional) The port to run the proxy on (defaults to 8080)

## Usage

1. **Set environment variables:**
   ```bash
   export BACKSTAGE_URL="https://your-backstage-instance.com"
   export MCP_TOKEN="your-bearer-token-here"
   export PORT="8080"  # Optional, defaults to 8080
   ```

2. **Run the proxy:**
   ```bash
   go run main.go
   ```

3. **Make requests through the proxy:**
   ```bash
   # The proxy only handles /api/mcp-actions paths
   curl http://localhost:8080/api/mcp-actions/v1
   curl http://localhost:8080/api/mcp-actions/v1/sse
   curl -X POST http://localhost:8080/api/mcp-actions/v1/call \
     -H "Content-Type: application/json" \
     -d '{"tool": "example", "params": {}}'
   
   # Other paths will return 404
   curl http://localhost:8080/api/mcp/tools  # Returns 404
   ```

## How it Works

1. The proxy receives incoming HTTP requests
2. **Path Filtering**: Only requests to `/api/mcp-actions` and its subpaths are proxied; all other requests return 404
3. It forwards the request to the configured Backstage URL with the same path and query parameters
4. It adds the Bearer token from `MCP_TOKEN` to the Authorization header
5. It forwards all relevant headers while filtering out hop-by-hop headers
6. For streaming responses (SSE, chunked encoding), it properly streams the response back to the client
7. For regular responses, it copies the response body back to the client

## Supported Endpoints

The proxy specifically handles these paths:
- `http://localhost:8080/api/mcp-actions` - Root MCP actions endpoint
- `http://localhost:8080/api/mcp-actions/v1` - MCP actions v1 endpoint
- `http://localhost:8080/api/mcp-actions/v1/sse` - Server-Sent Events endpoint
- Any other subpaths under `/api/mcp-actions/`

## Building

To build the proxy as a binary:

```bash
go build -o mcp-proxy main.go
./mcp-proxy
```

## Requirements

- Go 1.22 or newer
- Access to a Backstage instance with MCP server
- Valid Bearer token for the Backstage MCP server
