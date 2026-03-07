# RBAC Design: Role-Based Tool Filtering

## Context

go-unifi-mcp runs behind agentgateway, which multiplexes requests from
multiple AI agents. Different agents need different permission levels (e.g.,
one can read network state but not delete firewall rules). Agentgateway
currently does not forward caller identity, so RBAC is instance-level: each
container gets a role via env var.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `UNIFI_ROLE` | *(empty = no filtering)* | Role: `reader`, `operator`, or `admin` |

Invalid values are rejected at config load time.

## Role Definitions

| Role | Allowed verb prefixes | Approx tool count |
|---|---|---|
| `reader` | `list_`, `get_` | ~160 |
| `operator` | `list_`, `get_`, `create_`, `update_` | ~220 |
| `admin` | all (no filtering) | 242 |
| *(unset)* | all (backwards compatible) | 242 |

## Implementation

### Hook point

After `server.New()` registers all tools (eager or lazy mode), call
`rbac.Apply(s, role)` which deletes denied tools from the MCPServer. The
`tools/list` response and `tools/call` dispatch both read from the same map,
so deleting a tool hides it AND prevents execution.

### Package: `internal/rbac`

- `Apply(s, role)` — iterates registered tools, deletes those not matching
  the role's allowed verb prefixes
- Meta-tools (`tool_index`, `execute`, `batch`) are always allowed regardless
  of role
- Roles and verb mappings are hardcoded (no policy file)

### Call site

```go
s, err := r.newServer(server.Options{...})
rbac.Apply(s, cfg.Role)
return r.serve(s, cfg)
```

## Future evolution

- **Identity-aware RBAC**: When agentgateway gains identity forwarding, map
  identity → role via a policy file. The `Apply` function becomes
  per-session rather than per-instance.
- **Category-based filtering**: Add resource category restrictions within a
  role (e.g., reader on Clients but not Firewall).
- **Explicit tool lists**: Allow/deny specific tool names for maximum control.
