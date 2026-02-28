# Release Process

This project uses [Release Please](https://github.com/googleapis/release-please)
for automated releases.

## How It Works

### Every Push to Main

1. **Edge builds**: The `edge.yml` workflow builds snapshot binaries and pushes a
   `ghcr.io/cameronsjo/go-unifi-mcp:edge` Docker image
2. **Release PR**: Release Please creates or updates a pull request that bumps the
   version and updates `CHANGELOG.md` based on
   [Conventional Commits](https://www.conventionalcommits.org/)

### Cutting a Stable Release

Merge the Release Please PR. This automatically:

1. Creates a GitHub Release with the new tag (e.g., `v0.3.0`)
2. Triggers `release.yml`, which runs GoReleaser to:
   - Build binaries for linux/darwin (amd64/arm64)
   - Push Docker images to `ghcr.io/cameronsjo/go-unifi-mcp`
   - Update the Homebrew formula in `cameronsjo/homebrew-tap`
3. Publishes to the MCP Registry

### Commit Message Convention

Release Please determines version bumps from commit prefixes:

| Prefix | Version Bump | Example |
| ------ | ------------ | ------- |
| `fix:` | Patch (0.3.0 → 0.3.1) | `fix: resolve nil pointer in device list` |
| `feat:` | Minor (0.3.0 → 0.4.0) | `feat: add zone policy tools` |
| `feat!:` or `BREAKING CHANGE:` | Major (0.x → 1.0) | `feat!: rename tool_index to catalog` |
| `chore:`, `docs:`, `ci:` | No release | `ci: update golangci-lint` |

## Post-Release Verification

After the release workflow completes:

- [ ] Release created: `gh release view <tag>`
- [ ] Docker image published:
      `gh api user/packages/container/go-unifi-mcp/versions --jq '.[0].metadata.container.tags'`
- [ ] Homebrew formula updated: check `cameronsjo/homebrew-tap` for commit
- [ ] MCP Registry published (runs as part of release workflow)

## If Something Goes Wrong

- **Bad release notes**: Edit via `gh release edit <tag> --notes-file <file>`
- **Broken binary/image**: Fix the issue on main, then cut a patch release
- **Wrong tag**: Delete and retag only if the release workflow hasn't completed;
  otherwise cut a new patch release
