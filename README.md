# go-unifi-mcp

A Model Context Protocol (MCP) server for UniFi Network Controller, written in
Go.

## Overview

go-unifi-mcp provides an MCP interface to UniFi Network Controller, enabling AI
assistants and other MCP clients to interact with your UniFi infrastructure.

## Status

This project is under development. See the issue tracker for current progress.

## Development

### Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled
- [direnv](https://direnv.net/) (optional but recommended)

### Developing

```bash
# Clone the repository
git clone https://github.com/claytono/go-unifi-mcp.git
cd go-unifi-mcp

# Enter the development environment
nix develop
# Or with direnv:
direnv allow

# Install pre-commit hooks
pre-commit install

# Run linters
task lint

# Run tests
task test

# Run tests with coverage
task coverage
```

### Available Tasks

```bash
task lint        # Run linters via pre-commit
task test        # Run tests
task coverage    # Run tests with coverage checks
task build       # Build the binary
task generate    # Run go generate
```

## Credits

This project builds upon:

- [go-unifi](https://github.com/paultyng/go-unifi) - Go client library for UniFi
  Network Controller
- [unifi-network-mcp](https://github.com/sirkirby/unifi-network-mcp) - Python
  MCP server for UniFi that inspired this project
- [mcp-go](https://github.com/mark3labs/mcp-go) - Go SDK for Model Context
  Protocol

## License

MIT
