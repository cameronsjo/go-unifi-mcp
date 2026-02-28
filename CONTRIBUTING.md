# Contributing to go-unifi-mcp

Thank you for considering contributing to go-unifi-mcp! This document explains
how to get started.

## Code of Conduct

Be respectful and constructive. We're all here to build something useful.

## Getting Started

1. Fork the repository
2. Clone your fork:

   ```bash
   git clone https://github.com/<your-username>/go-unifi-mcp.git
   cd go-unifi-mcp
   ```

3. Set up the development environment (requires [Nix](https://nixos.org/)):

   ```bash
   nix develop
   # or use direnv for automatic activation
   ```

4. Create a feature branch:

   ```bash
   git checkout -b feat/your-feature main
   ```

## Development

### Common Commands

```bash
task lint        # Run linters
task test        # Run tests
task coverage    # Run tests with coverage checks
task build       # Build binary
task generate    # Run go generate
```

### Code Style

- Format with `goimports` (local prefix: `github.com/claytono/go-unifi-mcp`)
- Follow standard Go conventions
- Maintain 95% total test coverage, 90% per file

### Testing

- Write table-driven tests where applicable
- Use testify for assertions
- Mock external dependencies using mockery
- Run `task coverage` before submitting to ensure thresholds are met

## Commits

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

feat:     new feature
fix:      bug fix
docs:     documentation only
refactor: code change that neither fixes a bug nor adds a feature
test:     adding or updating tests
chore:    maintenance tasks
```

## Pull Requests

1. Keep PRs focused — one logical change per PR
2. Update tests for any code changes
3. Ensure all checks pass (`task lint && task test`)
4. Write a clear description of what changed and why
5. Link related issues using closing keywords (`Closes #123`)

## Reporting Issues

Use [GitHub Issues](https://github.com/claytono/go-unifi-mcp/issues) for bug
reports and feature requests. Include:

- Steps to reproduce (for bugs)
- Expected vs actual behavior
- Environment details (OS, Go version, UniFi controller version)
