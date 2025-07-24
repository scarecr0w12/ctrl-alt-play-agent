# MemoriPilot: System Architect

## Overview
This file contains the architectural decisions and design patterns for the MemoriPilot project.

## Architectural Decisions

- Use Go for performance and Docker API integration
- WebSocket for real-time Panel communication
- Environment variables for configuration management
- HTTP endpoint for health monitoring
- Dual protocol support for Panel compatibility
- Structured logging for operational visibility
- Graceful shutdown for container safety



1. **Decision 1**: Description of the decision and its rationale.
2. **Decision 2**: Description of the decision and its rationale.
3. **Decision 3**: Description of the decision and its rationale.



## Design Considerations

- Security: Bearer token authentication and secure Docker operations
- Performance: Efficient message handling and Docker API usage
- Reliability: Automatic reconnection and error recovery
- Compatibility: Support both new and legacy Panel protocols
- Monitoring: Comprehensive health checks and status reporting
- Deployment: Easy configuration via environment variables



## Components

### WebSocket Client

Manages real-time communication with Panel via WebSocket connection

**Responsibilities:**

- Establish and maintain WebSocket connection
- Handle authentication with bearer tokens
- Send/receive messages with Panel
- Implement automatic reconnection logic
- Process heartbeat messages

### Message Handler

Processes incoming commands and routes them to appropriate handlers

**Responsibilities:**

- Parse incoming WebSocket messages
- Route legacy message types to existing handlers
- Handle new Panel command format from Issue #27
- Maintain backwards compatibility
- Send structured responses

### Docker Manager

Interfaces with Docker API for container lifecycle management

**Responsibilities:**

- Create and configure game server containers
- Start, stop, restart container operations
- Monitor container status and resource usage
- Handle container cleanup and deletion
- Stream container logs and output

### Health Monitor

Provides health check endpoint and system monitoring

**Responsibilities:**

- HTTP health check endpoint on configurable port
- System resource monitoring
- Container status tracking
- Connection status reporting
- Uptime and performance metrics

### Configuration Manager

Handles environment-based configuration and validation

**Responsibilities:**

- Load configuration from environment variables
- Validate required settings
- Provide sensible defaults
- Support for different deployment environments
- Security credential management



