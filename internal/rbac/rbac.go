package rbac

import (
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

// Allowed verb prefixes per role. Tools whose names start with any of these
// prefixes are permitted; all others are deleted from the server.
var roles = map[string][]string{
	"reader":   {"list_", "get_"},
	"operator": {"list_", "get_", "create_", "update_"},
	"admin":    {}, // empty = allow all
}

// Meta-tools used by lazy mode are always allowed regardless of role.
var alwaysAllowed = map[string]bool{
	"tool_index": true,
	"execute":    true,
	"batch":      true,
}

// Apply filters the MCPServer's registered tools based on the given role.
// An empty role or "admin" performs no filtering. Unknown roles are a no-op
// (validation happens at config load time).
func Apply(s *server.MCPServer, role string) {
	if role == "" || role == "admin" {
		return
	}

	prefixes, ok := roles[role]
	if !ok {
		return
	}

	var denied []string
	for name := range s.ListTools() {
		if alwaysAllowed[name] {
			continue
		}
		if !hasAllowedPrefix(name, prefixes) {
			denied = append(denied, name)
		}
	}

	if len(denied) > 0 {
		s.DeleteTools(denied...)
		log.Printf("RBAC role=%q: removed %d tools, %d tools available", role, len(denied), len(s.ListTools()))
	}
}

func hasAllowedPrefix(name string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}
