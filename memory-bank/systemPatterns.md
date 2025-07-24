# System Patterns

## Architectural Patterns

- Pattern 1: Description

## Design Patterns

- Pattern 1: Description

## Common Idioms

- Idiom 1: Description

## Panel+Agent Command Protocol Breaking Changes

The Panel has implemented breaking changes in Issue #27 where commands now use a new standardized format. Instead of individual message types like 'server_start', the Panel now sends unified 'command' messages with an 'action' field specifying the operation. Agent must be updated to handle both new format and maintain backwards compatibility.

### Examples

- New format: {"type":"command","action":"start_server","id":"cmd_123"}
- Old format: {"type":"server_start","data":{...}}
- Agent needs handlePanelCommand() method for new format
- Must send structured responses with success/error fields


## Panel+Agent Communication Protocol with Issue #27 Breaking Changes

Updated with comprehensive Panel+Agent integration patterns. The Panel implemented Issue #27 with breaking changes requiring Agent updates. New unified command protocol uses type:'command' with action fields instead of individual message types. Agent must implement dual compatibility - new handlePanelCommand() for Panel commands and legacy handlers for backwards compatibility. All responses must follow structured format with success/error fields.

### Examples

- Panel sends: {"type":"command","action":"start_server","id":"cmd_123","serverId":"server_456"}
- Agent responds: {"type":"response","id":"cmd_123","success":true,"message":"Server started"}
- Legacy support: {"type":"server_start","data":{...}} still handled
- Error format: {"error":{"code":"CONTAINER_NOT_FOUND","message":"..."}}
- WebSocket connection with Bearer token authentication
- Heartbeat every 30 seconds for connection health
- Docker API integration for container lifecycle management
