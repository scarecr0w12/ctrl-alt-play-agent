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
