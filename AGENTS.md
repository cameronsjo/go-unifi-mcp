# go-unifi-mcp

MCP server for UniFi Network Controller written in Go.

## ABSOLUTE PROHIBITIONS

**The following actions are FORBIDDEN without explicit user permission. There
are NO exceptions. Do not rationalize, justify, or work around these rules.**

### 1. DO NOT CREATE OR MERGE PULL REQUESTS WITHOUT PERMISSION

- Always ASK before creating a PR - never create one autonomously
- NEVER merge PRs - merging is the user's responsibility
- Even if tests pass, even if it looks ready, even if the user seems busy - ASK
  FIRST

### 2. DO NOT BYPASS PRE-COMMIT HOOKS

- Never use `--no-verify` on any git command
- Never disable or skip hooks "temporarily"
- If hooks fail, FIX THE ISSUE - do not bypass
- Hook failures exist to catch problems - respect them

### 3. DO NOT CHANGE CODE COVERAGE THRESHOLDS

- The coverage thresholds in `.testcoverage.yaml` are SACRED
- Never lower `total` or `file` thresholds for any reason
- If coverage is failing, write more tests - do not lower the bar
- This includes "temporary" changes - there is no such thing

### 4. DO NOT RUSH

- Rushing leads to sloppy work that wastes MORE time later
- Take the time to do things correctly the first time
- If something seems hard, that's a sign to slow down, not speed up
- Quality gates exist for a reason - do not circumvent them to "save time"

### 5. DO NOT COMMIT DIRECTLY TO MAIN

- All changes to main require a pull request
- No direct commits, no matter how small the change
- This includes documentation, changelog updates, and release prep

**If you find yourself tempted to violate any of these rules, STOP and ask the
user for guidance. The answer is almost certainly "no, do it correctly."**

---

## Development

This project uses Nix for development environment management. Run `nix develop`
or use direnv to enter the development shell.

### Common Commands

```bash
task lint        # Run linters
task test        # Run tests
task coverage    # Run tests with coverage checks
task build       # Build binary
task generate    # Run go generate
```

### Project Structure

- `/internal/` - Internal packages (not exported)
- `/pkg/` - Public packages (exported API)
- `/cmd/` - Main entry points

### Code Style

- Use `goimports` for import organization (local prefix:
  `github.com/claytono/go-unifi-mcp`)
- Follow standard Go conventions
- Maintain 95% total test coverage, 90% per file (DO NOT CHANGE THESE VALUES)

### Type Generation

Types in this project are generated from the go-unifi library. Do not manually
edit generated type files. See Phase 2 documentation for sync process.

### Testing

- Write table-driven tests where applicable
- Use testify for assertions
- Mock external dependencies using mockery

### Mocks

- Generate mocks with `go generate ./internal/server/mocks`
- Mock definitions live in `internal/server/mocks`
- Do not edit generated mocks manually

### Releasing

**Read `RELEASE.md` before starting any release.** It contains the full release
rules including changelog style, version selection, README audit, and
post-release verification checklist.

## Issue Tracking

This project tracks work with **GitHub issues**. The `bd` (beads) tool was
retired (cadence-groundwork#138); `.beads/` remains in the repo as dormant
data, not the active queue.

### Session Completion

**When ending a work session**, make sure all work is committed and pushed
before finishing up. Follow-up work should be filed as a GitHub issue.
