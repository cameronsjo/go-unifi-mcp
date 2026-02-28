# ADR 0001: Initial Architecture

## Status

Accepted

## Context

go-unifi-mcp provides a Model Context Protocol (MCP) server that exposes UniFi
Network Controller functionality to AI assistants. The project needs a clear
architecture that supports:

- Type-safe interaction with the UniFi controller API
- Extensible tool registration for MCP
- Testability with high coverage requirements (95% total, 90% per file)

## Decision

### Language and Framework

- **Go** for the MCP server implementation
- **mcp-go** (`github.com/mark3labs/mcp-go`) as the MCP SDK
- **go-unifi** for UniFi API client types

### Project Structure

```
cmd/           # Entry points
internal/      # Private packages (server, handlers)
pkg/           # Public API surface
```

### Development Environment

- **Nix flake** for reproducible development environment
- **Task** (Taskfile.yml) as the task runner
- **golangci-lint** for linting
- **mockery** for test mock generation

### Testing Strategy

- Table-driven tests with testify assertions
- Generated mocks for external dependencies
- Coverage thresholds enforced in CI

## Consequences

- Nix provides reproducible builds but adds a learning curve for contributors
- Generated types from go-unifi reduce maintenance but require a sync process
- High coverage thresholds enforce quality but slow down initial development
