# Build stage
FROM golang:1.21-alpine AS builder

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o agent ./cmd/agent

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections and Docker CLI
RUN apk --no-cache add ca-certificates docker-cli

# Create a non-root user
RUN addgroup -g 1000 agent && \
    adduser -D -s /bin/sh -u 1000 -G agent agent

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/agent .

# Change ownership
RUN chown agent:agent /app/agent

# Switch to non-root user
USER agent

# Expose any required ports (if needed for future features)
# EXPOSE 8080

# Set environment variables
ENV PANEL_URL=""
ENV NODE_ID=""
ENV AGENT_SECRET=""

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ps aux | grep '[a]gent' || exit 1

# Run the agent
CMD ["./agent"]
