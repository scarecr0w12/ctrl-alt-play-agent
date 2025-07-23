#!/bin/bash

# Development script for Ctrl-Alt-Play Agent
set -e

echo "🚀 Ctrl-Alt-Play Agent Development Script"
echo "=========================================="

# Function to show usage
show_usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  test      - Run module tests"
    echo "  build     - Build the agent binary"
    echo "  run       - Build and run the agent"
    echo "  docker    - Build and run Docker image"
    echo "  clean     - Clean build artifacts"
    echo "  help      - Show this help"
    echo ""
    echo "Environment Variables:"
    echo "  PANEL_URL      - Panel WebSocket URL (default: ws://localhost:8080)"
    echo "  NODE_ID        - Node identifier (default: node-1)"
    echo "  AGENT_SECRET   - Authentication secret (default: agent-secret)"
}

# Function to run tests
run_tests() {
    echo "🧪 Running tests..."
    go run ./cmd/test
    echo "✅ Tests completed successfully!"
}

# Function to build
build_agent() {
    echo "🔨 Building agent..."
    make build
    echo "✅ Build completed!"
}

# Function to run agent
run_agent() {
    build_agent
    echo "🚀 Starting agent..."
    echo "Press Ctrl+C to stop"
    echo ""
    ./build/ctrl-alt-play-agent
}

# Function to run with Docker
run_docker() {
    echo "🐳 Building and running with Docker..."
    echo "Press Ctrl+C to stop"
    echo ""
    make docker-run
}

# Function to clean
clean_build() {
    echo "🧹 Cleaning build artifacts..."
    make clean
    echo "✅ Clean completed!"
}

# Main script logic
case "${1:-help}" in
    "test")
        run_tests
        ;;
    "build")
        build_agent
        ;;
    "run")
        run_agent
        ;;
    "docker")
        run_docker
        ;;
    "clean")
        clean_build
        ;;
    "help"|*)
        show_usage
        ;;
esac
