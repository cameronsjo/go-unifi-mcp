package mcpgen

import (
	"github.com/claytono/go-unifi-mcp/internal/gounifi"
)

// InferOperations determines which CRUD operations are available for a resource.
func InferOperations(r *gounifi.Resource) []string {
	// Settings resources only have Get and Update
	if r.IsSetting() {
		return []string{"Get", "Update"}
	}

	// Device resource is read-only (List, Get only)
	if r.StructName == "Device" {
		return []string{"List", "Get"}
	}

	// All other resources have full CRUD
	return []string{"List", "Get", "Create", "Update", "Delete"}
}

// HasOperation checks if a resource supports a specific operation.
func HasOperation(r *gounifi.Resource, op string) bool {
	ops := InferOperations(r)
	for _, o := range ops {
		if o == op {
			return true
		}
	}
	return false
}
