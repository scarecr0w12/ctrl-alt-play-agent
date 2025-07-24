# MemoriPilot: System Architect

## Overview
This file contains the architectural decisions and design patterns for the MemoriPilot project.

## Architectural Decisions

- Use Go for performance, Docker API integration, and concurrent operations
- WebSocket for real-time Panel communication with automatic reconnection
- Environment variables for configuration management and deployment flexibility
- HTTP endpoint for health monitoring and external service integration
- Dual protocol support for Panel compatibility during migration period
- Structured JSON logging for operational visibility and debugging
- Graceful shutdown with proper container cleanup and resource deallocation
- Bearer token authentication for secure Panel communication
- Docker API for direct container management without shell dependencies



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

- Security: Bearer token authentication, secure Docker operations, and minimal privilege execution
- Performance: Efficient message handling, Docker API optimization, and resource management
- Reliability: Automatic reconnection, error recovery, and graceful degradation
- Compatibility: Support both new Panel protocol and legacy message types during migration
- Monitoring: Comprehensive health checks, metrics collection, and status reporting
- Deployment: Easy configuration via environment variables and containerized deployment
- Scalability: Support for multiple concurrent server management and resource optimization
- Maintainability: Clean code architecture, comprehensive logging, and documentation



- Security: Bearer token authentication and secure Docker operations
- Performance: Efficient message handling and Docker API usage
- Reliability: Automatic reconnection and error recovery
- Compatibility: Support both new and legacy Panel protocols
- Monitoring: Comprehensive health checks and status reporting
- Deployment: Easy configuration via environment variables



## Components

### WebSocket Client

Manages secure real-time communication with Panel via WebSocket connection

**Responsibilities:**

- Establish and maintain WebSocket connection to Panel
- Handle bearer token authentication
- Send/receive JSON messages with Panel
- Implement automatic reconnection with exponential backoff
- Process heartbeat messages every 30 seconds
- Handle connection state management

### Message Handler

Processes incoming commands and routes them to appropriate handlers with protocol compatibility

**Responsibilities:**

- Parse incoming WebSocket messages
- Route legacy message types to existing handlers
- Handle new Panel command format from Issue #27
- Maintain backwards compatibility for smooth migration
- Send structured responses with proper message ID tracking
- Implement error handling and validation

### Docker Manager

Interfaces with Docker API for complete container lifecycle management

**Responsibilities:**

- Create and configure game server containers
- Start, stop, restart container operations with signal support
- Monitor container status and resource usage
- Handle container cleanup and deletion
- Stream container logs and output to Panel
- Manage container networks and volumes

### Health Monitor

Provides health check endpoint and comprehensive system monitoring

**Responsibilities:**

- HTTP health check endpoint on configurable port
- System resource monitoring (CPU, memory, disk)
- Container status tracking and reporting
- Connection status reporting to Panel
- Uptime and performance metrics collection
- Error rate and response time tracking

### Configuration Manager

Handles environment-based configuration with validation and security

**Responsibilities:**

- Load configuration from environment variables
- Validate required settings and provide helpful error messages
- Provide sensible defaults for optional settings
- Support for different deployment environments
- Security credential management and validation
- Configuration hot-reload capabilities

### Logging System

Manages structured logging and operational visibility

**Responsibilities:**

- Structured JSON logging for operational visibility
- Log level management (debug, info, warn, error)
- Log rotation and retention policies
- Error tracking and aggregation
- Performance metrics logging
- Audit trail for security events





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



