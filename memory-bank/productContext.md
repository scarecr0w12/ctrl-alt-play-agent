# Product Context

## Overview

Ctrl-Alt-Play Agent is a lightweight, production-ready game server management agent designed to work seamlessly with the Ctrl-Alt-Play Panel system. It serves as the distributed execution layer in a Panel+Agent architecture, enabling centralized management of game servers across multiple remote Linux systems.

## Project Description

Ctrl-Alt-Play Agent v1.1.0 - A lightweight, high-performance remote game server management agent with dual communication architecture (HTTP API + WebSocket) for seamless integration with Ctrl-Alt-Play Panel



The Agent functions similarly to the "Wings" system found in Pelican Panel and Pterodactyl Panel, providing a secure and efficient bridge between the central Panel and Docker-based game servers. It handles the complete lifecycle of game server containers, from creation and configuration to monitoring and cleanup, while maintaining real-time communication with the Panel for immediate status updates and command execution.

**Current Status**: Requires protocol updates to align with Panel's Issue #27 breaking changes that introduced a new unified command format.

## Core Features

- **Real-time Communication**: WebSocket-based bidirectional messaging with the Panel
- **Docker Integration**: Complete container lifecycle management via Docker API
- **Security**: Bearer token authentication and secure container operations
- **Health Monitoring**: HTTP endpoint for system health checks and status reporting
- **Automatic Reconnection**: Robust connection handling with heartbeat mechanism
- **Resource Management**: Monitor and report system and container resource usage
- **Log Streaming**: Real-time server console output streaming to Panel
- **File Operations**: Basic file management capabilities for server configurations
- **Signal Handling**: Graceful shutdown and container cleanup on termination

## Architecture

Dual communication architecture with HTTP REST API server (port 8081) for panel discovery and command execution, plus WebSocket client for real-time communication. Single combined server hosting health endpoint and API commands with CORS support and multiple authentication methods (X-API-Key and Bearer token).



**Distributed Panel+Agent System**:
- Agent runs on remote nodes with Docker access
- Communicates with central Panel via WebSocket (port 8080)
- Uses Docker API for container lifecycle management
- Bearer token authentication for secure Panel communication
- Heartbeat mechanism maintains connection health (30-second intervals)
- HTTP health check endpoint for monitoring (configurable port)

**Recent Changes**: Panel Issue #27 introduced breaking changes requiring Agent protocol updates for new unified command format.

## Technologies

- Go 1.23+
- Docker API
- WebSocket
- HTTP REST API
- JSON
- CORS
- Docker Compose
- Kubernetes
- Systemd



- **Go 1.21+**: Core runtime and development platform
- **Docker Engine**: Container management and orchestration
- **WebSocket**: Real-time bidirectional communication protocol
- **Linux**: Target deployment platform
- **JSON**: Message serialization and configuration format
- **Bearer Token Authentication**: Secure Panel authentication
- **Docker API**: Container lifecycle and monitoring interface
- **HTTP**: Health check and status endpoint

## Libraries and Dependencies

- github.com/docker/docker v28.3.2+incompatible
- github.com/gorilla/websocket v1.5.3
- github.com/stretchr/testify v1.10.0



**Core Dependencies**:
- `github.com/gorilla/websocket` - WebSocket client implementation
- `github.com/docker/docker` - Docker API client library

**Standard Library**:
- `context` - Request context and cancellation
- `net/http` - HTTP server for health checks
- `encoding/json` - JSON message serialization
- `time` - Timestamp and duration handling
- `log` - Application logging
- `os` - Operating system interface
- `syscall` - System call interface
- `runtime` - Go runtime information
- `strconv` - String conversion utilities
- `sync` - Synchronization primitives
- `net/url` - URL parsing for WebSocket connections

