# Architecture Documentation

## System Overview

The Ctrl-Alt-Play Agent is designed as a lightweight, distributed node agent that works in conjunction with the Ctrl-Alt-Play Panel to provide comprehensive game server management capabilities.

## High-Level Architecture

```text
┌─────────────────────────────────────────────────────────────────┐
│                    Ctrl-Alt-Play Panel                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Web UI        │  │  ExternalAgent  │  │ AgentDiscovery  │ │
│  │   (Frontend)    │  │    Service      │  │    Service      │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│           │                     │                     │        │
│           │              HTTP API Calls      Multi-Port Scan   │
│           │                     │                     │        │
│           └─────────────────────┼─────────────────────┼────────┘
│                                 │                     │
│                          WebSocket Connection         │
│                                 │                     │
└─────────────────────────────────┼─────────────────────┼────────┘
                                  │                     │
                    ┌─────────────┼─────────────────────┼─────────┐
                    │             ▼                     ▼         │
                    │    ┌─────────────────┐  ┌─────────────────┐ │
                    │    │  WebSocket      │  │   HTTP API      │ │
                    │    │   Client        │  │    Server       │ │
                    │    └─────────────────┘  └─────────────────┘ │
                    │             │                     │         │
                    │             └─────────┬───────────┘         │
                    │                       │                     │
                    │         ┌─────────────▼─────────────┐       │
                    │         │       Agent Core          │       │
                    │         │   (Command Processing)    │       │
                    │         └─────────────┬─────────────┘       │
                    │                       │                     │
                    │         ┌─────────────▼─────────────┐       │
                    │         │     Docker Manager        │       │
                    │         │  (Container Operations)   │       │
                    │         └─────────────┬─────────────┘       │
                    │                       │                     │
                    │                       ▼                     │
                    │              Docker Daemon                  │
                    │         (Game Server Containers)            │
                    └─────────────────────────────────────────────┘
                                   Agent Node
```

## Component Architecture

### 1. Communication Layer

#### HTTP API Server
- **Purpose**: Provides REST API endpoints for panel communication
- **Port**: 8081 (configurable via `HEALTH_PORT`)
- **Endpoints**:
  - `/health` - Health check and status
  - `/api/command` - Command execution
- **Authentication**: X-API-Key header validation
- **Features**: CORS support, JSON request/response

#### WebSocket Client
- **Purpose**: Real-time bidirectional communication with panel
- **Connection**: Connects to panel WebSocket server
- **Authentication**: Bearer token in connection headers
- **Features**: Automatic reconnection, heartbeat monitoring

### 2. Core Services

#### Agent Core
- **Responsibilities**:
  - Command routing and execution
  - Request validation and authentication
  - Response formatting
  - Error handling and logging
- **Components**:
  - Command dispatcher
  - Authentication middleware
  - Response serializer

#### Configuration Manager
- **Purpose**: Centralized configuration management
- **Sources**:
  - Environment variables
  - Default values
  - Runtime configuration
- **Parameters**:
  - Panel connection details
  - Authentication secrets
  - Network configuration

#### Health Monitor
- **Purpose**: System health tracking and reporting
- **Metrics**:
  - Connection status (panel connectivity)
  - System uptime
  - Version information
  - Node identification

### 3. Container Management Layer

#### Docker Manager
- **Purpose**: Direct Docker API integration
- **Capabilities**:
  - Container lifecycle management (create, start, stop, remove)
  - Container inspection and monitoring
  - Resource management and limits
  - Network and volume management
- **Context**: Uses Go's context package for cancellation and timeouts

#### Game Server Abstraction
- **Purpose**: Higher-level game server operations
- **Features**:
  - Server configuration management
  - Port mapping and networking
  - Environment variable injection
  - Resource limit enforcement

## Communication Patterns

### 1. Panel Discovery Process

```text
Panel AgentDiscoveryService
         │
         ├─ Scan Port 8081 ──────┐
         ├─ Scan Port 8080       │
         ├─ Scan Port+1          │
         └─ Scan Port+100        │
                                 │
                                 ▼
                    Agent HTTP Server (8081)
                                 │
                                 ▼
                         GET /health
                                 │
                                 ▼
                    {"status": "healthy", "nodeId": "node-1"}
                                 │
                                 ▼
                    Panel registers agent with baseURL and apiKey
```

### 2. Command Execution Flow

```text
Panel ExternalAgentService
         │
         ▼
POST /api/command
Header: X-API-Key: secret
Body: {"action": "docker.list", "data": {}}
         │
         ▼
Agent HTTP API Server
         │
         ├─ Authenticate request
         ├─ Parse command
         ├─ Route to handler
         └─ Execute via Docker Manager
         │
         ▼
Docker API
         │
         ▼
Response: {"success": true, "data": {"containers": [...]}}
         │
         ▼
Panel receives result
```

### 3. Real-time Updates (WebSocket)

```text
Agent WebSocket Client ←────→ Panel WebSocket Server
         │                           │
         ├─ Connect with Bearer       │
         ├─ Send heartbeats          │
         ├─ Receive commands         │
         ├─ Send status updates      │
         └─ Handle reconnection      │
                                     │
                    Real-time event streaming
```

## Data Flow

### 1. Request Processing Pipeline

```text
HTTP Request → Authentication → Validation → Command Dispatch → Docker API → Response
     │              │              │              │              │           │
     ▼              ▼              ▼              ▼              ▼           ▼
 Raw JSON    Check X-API-Key   Parse action   Route to        Call Docker  JSON
 Body        or Bearer token   and data       handler method  daemon       Response
```

### 2. Error Handling Flow

```text
Error Occurrence
     │
     ├─ Log error details
     ├─ Determine error type
     ├─ Format error response
     └─ Set appropriate HTTP status
     │
     ▼
Client receives structured error:
{"success": false, "error": "description"}
```

## Security Architecture

### 1. Authentication Layers

```text
External Request
     │
     ▼
┌─────────────────┐
│  CORS Headers   │  ← Browser security
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Authentication  │  ← X-API-Key or Bearer token
│   Middleware    │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Authorization   │  ← Command-level permissions
│    Check        │
└─────────────────┘
     │
     ▼
Command Execution
```

### 2. Network Security

```text
Internet/WAN
     │
     ▼
┌─────────────────┐
│   Firewall      │  ← Port restrictions (8081, panel connection)
└─────────────────┘
     │
     ▼
┌─────────────────┐
│  Reverse Proxy  │  ← Optional: nginx, Apache (HTTPS termination)
│   (Optional)    │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│   Agent Host    │  ← Agent process
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Docker Daemon   │  ← Unix socket access required
└─────────────────┘
```

## Deployment Patterns

### 1. Single-Node Deployment

```text
┌─────────────────────────────────────┐
│           Host System               │
│  ┌─────────────────────────────────┐│
│  │        Agent Process            ││
│  │  ┌─────────────────────────────┐││
│  │  │     Docker Daemon           │││
│  │  │  ┌─────────────────────────┐│││
│  │  │  │   Game Containers       ││││
│  │  │  │  ┌─────┐ ┌─────┐ ┌─────┐││││
│  │  │  │  │Game1│ │Game2│ │Game3│││││
│  │  │  │  └─────┘ └─────┘ └─────┘││││
│  │  │  └─────────────────────────┘│││
│  │  └─────────────────────────────┘││
│  └─────────────────────────────────┘│
└─────────────────────────────────────┘
```

### 2. Multi-Node Deployment

```text
┌─────────────────────────────────────┐
│            Panel Host               │
│         (Central Control)           │
└─────────────────┬───────────────────┘
                  │
        ┌─────────┼─────────┐
        │         │         │
        ▼         ▼         ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│   Agent 1   │ │   Agent 2   │ │   Agent N   │
│ (US-East)   │ │ (US-West)   │ │  (Europe)   │
└─────────────┘ └─────────────┘ └─────────────┘
```

### 3. Container Deployment

```text
┌─────────────────────────────────────┐
│           Docker Host               │
│  ┌─────────────────────────────────┐│
│  │       Agent Container           ││
│  │  /var/run/docker.sock mounted   ││
│  └─────────────────────────────────┘│
│  ┌─────────────────────────────────┐│
│  │      Game Container 1           ││
│  └─────────────────────────────────┘│
│  ┌─────────────────────────────────┐│
│  │      Game Container 2           ││
│  └─────────────────────────────────┘│
└─────────────────────────────────────┘
```

## Performance Considerations

### 1. Scalability

- **Horizontal Scaling**: Deploy multiple agents across different hosts
- **Vertical Scaling**: Increase host resources for more containers per agent
- **Load Distribution**: Panel can distribute workload across multiple agents

### 2. Resource Management

- **Memory Usage**: Minimal agent footprint (~50MB base memory)
- **CPU Usage**: Low baseline CPU usage, spikes during container operations
- **Network**: Persistent WebSocket connection, periodic HTTP health checks
- **Disk I/O**: Primarily for Docker operations and logging

### 3. Optimization Strategies

- **Connection Pooling**: Reuse Docker API connections
- **Caching**: Cache container state to reduce Docker API calls
- **Batching**: Batch multiple operations where possible
- **Async Operations**: Non-blocking container operations

## Monitoring and Observability

### 1. Health Metrics

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Health Status  │    │ Connection Status│   │ System Metrics  │
│                 │    │                 │    │                 │
│ • Healthy       │    │ • Connected     │    │ • Uptime        │
│ • Degraded      │    │ • Disconnected  │    │ • Memory Usage  │
│ • Unhealthy     │    │ • Reconnecting  │    │ • CPU Usage     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 2. Logging Architecture

```text
Application Logs → Structured JSON → Log Aggregation → Monitoring Dashboard
     │                    │                │                    │
     ▼                    ▼                ▼                    ▼
 • Commands          • Timestamp      • ELK Stack          • Grafana
 • Errors            • Level          • Fluentd            • Kibana
 • Connections       • Message        • Loki               • Prometheus
 • Docker Events     • Context        • Splunk             • Custom UI
```

## Future Architecture Considerations

### 1. Microservices Evolution

- **Service Decomposition**: Split agent into focused microservices
- **Message Queues**: Implement async messaging (Redis, RabbitMQ)
- **Service Discovery**: Automatic agent registration and discovery
- **Circuit Breakers**: Resilience patterns for external dependencies

### 2. Cloud-Native Features

- **Kubernetes Integration**: Native K8s deployment and scaling
- **Service Mesh**: Istio/Linkerd for inter-service communication
- **Cloud Storage**: Integration with cloud storage for game data
- **Auto-scaling**: Horizontal pod autoscaling based on workload

### 3. Enhanced Security

- **mTLS**: Mutual TLS for all inter-service communication
- **Vault Integration**: HashiCorp Vault for secret management
- **RBAC**: Role-based access control for agent operations
- **Audit Logging**: Comprehensive audit trail for all operations

This architecture provides a solid foundation for scalable, secure, and maintainable game server management while maintaining simplicity and performance.
