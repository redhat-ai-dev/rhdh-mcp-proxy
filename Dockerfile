# Use UBI9 Go toolset for building
FROM registry.access.redhat.com/ubi9/go-toolset:1.22 AS builder

WORKDIR /opt/app-root/src

# Copy go mod files
COPY go.mod go.mod

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o rhdh-mcp-proxy main.go

# Final stage using UBI minimal for smaller footprint
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Install ca-certificates using microdnf
RUN microdnf update -y && \
    microdnf install -y ca-certificates && \
    microdnf clean all

# Create a non-root user for security
RUN useradd -r -u 1001 -g 0 rhdh-mcp-proxy

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /opt/app-root/src/rhdh-mcp-proxy .

# Change ownership to the non-root user
RUN chown -R 1001:0 /app && chmod -R g=u /app

# Switch to non-root user
USER 1001

# Run the binary
CMD ["./rhdh-mcp-proxy"] 