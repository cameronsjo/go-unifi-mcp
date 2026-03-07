package rbac

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(toolNames ...string) *server.MCPServer {
	s := server.NewMCPServer("test", "0.0.0")
	noop := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return nil, nil
	}
	for _, name := range toolNames {
		s.AddTools(server.ServerTool{
			Tool:    mcp.Tool{Name: name},
			Handler: noop,
		})
	}
	return s
}

func toolNames(s *server.MCPServer) []string {
	tools := s.ListTools()
	names := make([]string, 0, len(tools))
	for name := range tools {
		names = append(names, name)
	}
	return names
}

func TestApply_EmptyRoleNoFiltering(t *testing.T) {
	s := newTestServer("list_network", "get_network", "delete_network")
	Apply(s, "")
	assert.Len(t, s.ListTools(), 3)
}

func TestApply_AdminNoFiltering(t *testing.T) {
	s := newTestServer("list_network", "get_network", "delete_network")
	Apply(s, "admin")
	assert.Len(t, s.ListTools(), 3)
}

func TestApply_ReaderAllowsListAndGet(t *testing.T) {
	s := newTestServer(
		"list_network", "list_client",
		"get_network", "get_client",
		"create_network", "update_network", "delete_network",
	)

	Apply(s, "reader")

	names := toolNames(s)
	assert.Len(t, names, 4)
	assert.Contains(t, names, "list_network")
	assert.Contains(t, names, "list_client")
	assert.Contains(t, names, "get_network")
	assert.Contains(t, names, "get_client")
	assert.NotContains(t, names, "create_network")
	assert.NotContains(t, names, "update_network")
	assert.NotContains(t, names, "delete_network")
}

func TestApply_OperatorAllowsListGetCreateUpdate(t *testing.T) {
	s := newTestServer(
		"list_network", "get_network",
		"create_network", "update_network",
		"delete_network",
	)

	Apply(s, "operator")

	names := toolNames(s)
	assert.Len(t, names, 4)
	assert.Contains(t, names, "list_network")
	assert.Contains(t, names, "get_network")
	assert.Contains(t, names, "create_network")
	assert.Contains(t, names, "update_network")
	assert.NotContains(t, names, "delete_network")
}

func TestApply_PreservesMetaTools(t *testing.T) {
	s := newTestServer(
		"tool_index", "execute", "batch",
		"list_network", "delete_network",
	)

	Apply(s, "reader")

	names := toolNames(s)
	assert.Contains(t, names, "tool_index")
	assert.Contains(t, names, "execute")
	assert.Contains(t, names, "batch")
	assert.Contains(t, names, "list_network")
	assert.NotContains(t, names, "delete_network")
}

func TestApply_UnknownRoleNoOp(t *testing.T) {
	s := newTestServer("list_network", "delete_network")
	Apply(s, "unknown")
	assert.Len(t, s.ListTools(), 2)
}

func TestApply_ReaderWithNoMatchingTools(t *testing.T) {
	s := newTestServer("delete_network", "create_network")
	Apply(s, "reader")
	assert.Empty(t, s.ListTools())
}

func TestHasAllowedPrefix(t *testing.T) {
	tests := []struct {
		name     string
		tool     string
		prefixes []string
		expected bool
	}{
		{"match first", "list_network", []string{"list_", "get_"}, true},
		{"match second", "get_client", []string{"list_", "get_"}, true},
		{"no match", "delete_network", []string{"list_", "get_"}, false},
		{"empty prefixes", "list_network", []string{}, false},
		{"exact prefix", "list_", []string{"list_"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, hasAllowedPrefix(tt.tool, tt.prefixes))
		})
	}
}

func TestRolesDefinition(t *testing.T) {
	// Verify all expected roles exist
	require.Contains(t, roles, "reader")
	require.Contains(t, roles, "operator")
	require.Contains(t, roles, "admin")

	// Reader is a subset of operator
	for _, prefix := range roles["reader"] {
		assert.Contains(t, roles["operator"], prefix)
	}

	// Admin has empty prefixes (allow all)
	assert.Empty(t, roles["admin"])
}
