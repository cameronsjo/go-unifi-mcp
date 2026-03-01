# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0](https://github.com/cameronsjo/go-unifi-mcp/compare/v0.2.0...v0.3.0) (2026-03-01)


### Features

* add AGENTS.md, AI symlinks, CONTRIBUTING, SECURITY, and initial ADR ([b0a8260](https://github.com/cameronsjo/go-unifi-mcp/commit/b0a8260c9053b5e3e09f17fe0f7a7e8f65bbcd76))
* apply essentials scaffold ([809c5b1](https://github.com/cameronsjo/go-unifi-mcp/commit/809c5b178944184304d994a8e3fc4b8c84c4da6a))
* configure Homebrew tap and Docker for cameronsjo fork ([b1f585b](https://github.com/cameronsjo/go-unifi-mcp/commit/b1f585bc8b41a00cef99d7c621a5671b54c6fbfa))
* **mcpgen:** add filterable enum hints to list tool descriptions ([5c7896d](https://github.com/cameronsjo/go-unifi-mcp/commit/5c7896d1dc4fa4437aec6878c2b68cd093f41136))


### Bug Fixes

* add filterable enum hints and remove broken kacl hook ([e28a795](https://github.com/cameronsjo/go-unifi-mcp/commit/e28a795df559108e946a73be052aaf5f8a2c21c5))
* **ci:** remove python-kacl pre-commit hook ([e3115cd](https://github.com/cameronsjo/go-unifi-mcp/commit/e3115cdc7d3e2d8868edfdb89db14b557af76052))
* correct release-please-action SHA pin ([#7](https://github.com/cameronsjo/go-unifi-mcp/issues/7)) ([2aa1dfd](https://github.com/cameronsjo/go-unifi-mcp/commit/2aa1dfdd96408fdc41f8ee4e85eb8c782fed69f2))
* **mcpgen:** use regex assertion for enum hint test ([0cd941b](https://github.com/cameronsjo/go-unifi-mcp/commit/0cd941b67a240edf9ac678042970de7918df40e5))
* remove --fix from golangci-lint pre-commit hook ([#6](https://github.com/cameronsjo/go-unifi-mcp/issues/6)) ([48531b0](https://github.com/cameronsjo/go-unifi-mcp/commit/48531b0379d3a95e38c7b9116929b68b0c9b88f5))

## [0.2.0] - 2026-02-03

### Added

- Add automatic ID reference resolution for tool responses — list, get, create,
  and update tools now resolve opaque `_id` fields to human-readable `_name`
  siblings using per-request caching. For example,
  `"network_id": "609fbf24e3ae433962e000de"` becomes
  `"network_id": "609fbf24e3ae433962e000de", "network_name": "IOT"`. Resolution
  is on by default; pass `"resolve": false` to disable (#63)
- Add configurable log level via `UNIFI_LOG_LEVEL` environment variable,
  defaulting to `error` to prevent go-unifi INFO messages from breaking piped
  JSON workflows (#62)
- Add `filter`, `fields`, and `search` post-processing parameters to all list
  operations for client-side filtering. `filter` supports exact match, contains,
  and regex operators (e.g. `"filter": {"name": {"contains": "office"}}`).
  `fields` projects the response to specific keys (e.g.
  `"fields": ["name", "ip"]`). `search` does a full-text search across all
  string values (#61)

### Changed

- Pin and update GitHub Actions dependencies

## [0.1.1] - 2026-01-30

### Added

- `--help` and `--version` CLI flags
- Automatic publishing to MCP Registry on release

## [0.1.0] - 2026-01-29

### Added

- MCP server for UniFi Network Controller with 240+ auto-generated tools
- Lazy mode with 3 meta-tools (tool_index, execute, batch) for LLM-friendly
  operation
- Eager mode exposing all tools directly for MCP clients that handle large tool
  sets
- Docker multi-arch images (linux/amd64, linux/arm64) published to ghcr.io
- Homebrew tap install via `brew install claytono/tap/go-unifi-mcp`
- Nix flake support for reproducible builds and installation
- Pre-built binaries for linux and macOS (amd64 and arm64)
- Configurable site selection and authentication via environment variables

[0.2.0]: https://github.com/claytono/go-unifi-mcp/releases/tag/v0.2.0
[0.1.1]: https://github.com/claytono/go-unifi-mcp/releases/tag/v0.1.1
[0.1.0]: https://github.com/claytono/go-unifi-mcp/releases/tag/v0.1.0

<!-- Versions 0.1.0–0.2.0 were released from the upstream repo (claytono/go-unifi-mcp).
     Fork releases (cameronsjo/go-unifi-mcp) begin at 0.3.0. -->
