# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.1] - 2025-08-01

### Added

- **Enhanced Documentation**: Updated project documentation and architecture details
- **Version Consistency**: Synchronized versioning across all project files

### Changed

- **Build System**: Improved build and test reliability
- **Project Structure**: Better organization of documentation and version files

### Fixed

- **Build Issues**: Resolved compilation and test failures
- **Documentation**: Updated version references and project status

## [1.1.0] - 2025-07-25

### Added
- **Dual Communication Architecture**: HTTP REST API server alongside WebSocket client
- **HTTP API Endpoints**: `/health` for discovery and `/api/command` for execution
- **Panel Discovery Compatibility**: Multi-port scanning support (8081, 8080)
- **Enhanced Authentication**: X-API-Key header and Bearer token support
- **CORS Support**: Full browser compatibility for panel UI integration
- **Graceful Error Handling**: Agent continues running if panel unavailable
- **Combined Server**: Single port (8081) hosting both health and API endpoints
- **Comprehensive Documentation**: API docs, deployment guide, architecture docs
- **Versioning System**: Aligned with panel versioning (package.json + VERSION file)
- **Docker Command Support**: Complete container lifecycle management via API
- **System Monitoring**: System status commands (uptime, memory, disk)

### Changed
- **Architecture**: Moved from WebSocket-only to dual HTTP/WebSocket architecture
- **Port Configuration**: Standardized on port 8081 for HTTP API
- **Authentication**: Enhanced to support both panel authentication methods
- **Error Handling**: Non-fatal panel connection failures with degraded mode
- **Documentation**: Complete overhaul with proper API, deployment, and architecture docs
- **Version Management**: Updated to v1.1.0 to align with panel versioning

### Fixed
- **Panel Integration**: Resolved communication protocol mismatch with panel services
- **Discovery Issues**: Agent now responds correctly to panel's AgentDiscoveryService
- **Command Execution**: Fixed command routing and response formatting
- **Connection Stability**: Improved WebSocket connection handling and recovery

### Technical Details
- **ExternalAgentService Compatibility**: Full support for panel's HTTP command execution
- **AgentDiscoveryService Integration**: Health endpoint for automatic agent registration
- **Docker API Context**: Proper context handling for all Docker operations
- **JSON Response Format**: Standardized success/error response structure

## [1.0.0] - 2025-07-24

### Added
- Panel Issue #27 unified command protocol support
- Comprehensive test suites for all components
- GitHub Actions CI/CD pipeline with multi-platform builds
- Security scanning with Trivy vulnerability scanner
- Docker container publishing to registry
- Memory bank system for project management
- Protocol compatibility with backwards support
- Health check endpoint for monitoring
- Cross-platform binary builds (Linux, macOS, Windows)

### Changed
- Updated message protocol to support Panel Issue #27 format
- Enhanced error handling with structured error codes
- Improved WebSocket client with automatic reconnection
- Modernized project structure and dependencies

### Fixed
- Protocol compatibility issues between Agent and Panel
- Build system configuration for multi-platform support
- Test coverage gaps in core components
- Documentation inconsistencies

### Security
- Added security scanning to CI/CD pipeline
- Implemented proper JWT Bearer token authentication
- Enhanced connection security with TLS support

## [v1.0.0] - Initial Release

### Added
- Basic WebSocket client for Panel communication
- Docker container management functionality
- Core message handling system
- Basic configuration management
- Health monitoring capabilities
