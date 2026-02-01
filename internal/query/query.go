// Package query provides generic post-processing for list operation results.
// It supports field filtering, text search, and field projection.
package query

import (
	"fmt"
	"regexp"
	"strings"
)

// Options holds query parameters extracted from MCP request arguments.
type Options struct {
	Filter map[string]any // field -> value | {"contains":"..."} | {"regex":"..."}
	Search string         // case-insensitive text search across top-level string values
	Fields []string       // field projection (nil = all fields)
}

// HasQuery returns true if any query parameters are set.
func (o Options) HasQuery() bool {
	return len(o.Filter) > 0 || o.Search != "" || len(o.Fields) > 0
}

// ParseOptions extracts query options from MCP request arguments.
func ParseOptions(args map[string]any) Options {
	var opts Options
	if f, ok := args["filter"].(map[string]any); ok {
		opts.Filter = f
	}
	if s, ok := args["search"].(string); ok {
		opts.Search = s
	}
	if arr, ok := args["fields"].([]any); ok {
		for _, v := range arr {
			if s, ok := v.(string); ok {
				opts.Fields = append(opts.Fields, s)
			}
		}
	}
	return opts
}

// Apply applies filter, search, and fields operations to a slice of items.
// Order: filter → search → fields (projection last so filter/search see all fields).
func Apply(items []map[string]any, opts Options) []map[string]any {
	if len(items) == 0 {
		return items
	}

	result := items

	if len(opts.Filter) > 0 {
		result = applyFilter(result, opts.Filter)
	}
	if opts.Search != "" {
		result = applySearch(result, opts.Search)
	}
	if len(opts.Fields) > 0 {
		result = applyFields(result, opts.Fields)
	}

	return result
}

func applyFilter(items []map[string]any, filter map[string]any) []map[string]any {
	result := make([]map[string]any, 0)
	for _, item := range items {
		if matchesFilter(item, filter) {
			result = append(result, item)
		}
	}
	return result
}

func matchesFilter(item map[string]any, filter map[string]any) bool {
	for field, filterValue := range filter {
		fieldValue, exists := item[field]
		if !exists {
			return false
		}
		if !matchesFieldFilter(fieldValue, filterValue) {
			return false
		}
	}
	return true
}

func matchesFieldFilter(fieldValue any, filterValue any) bool {
	// Check if filterValue is an operator object like {"contains": "..."} or {"regex": "..."}
	if opMap, ok := filterValue.(map[string]any); ok {
		if containsVal, ok := opMap["contains"].(string); ok {
			fieldStr := fmt.Sprintf("%v", fieldValue)
			return strings.Contains(strings.ToLower(fieldStr), strings.ToLower(containsVal))
		}
		if regexVal, ok := opMap["regex"].(string); ok {
			fieldStr := fmt.Sprintf("%v", fieldValue)
			matched, err := regexp.MatchString(regexVal, fieldStr)
			if err != nil {
				return false // invalid regex → no match
			}
			return matched
		}
		return false // unknown operator
	}

	// Exact match: compare string representations
	return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", filterValue)
}

func applySearch(items []map[string]any, search string) []map[string]any {
	searchLower := strings.ToLower(search)
	result := make([]map[string]any, 0)
	for _, item := range items {
		if matchesSearch(item, searchLower) {
			result = append(result, item)
		}
	}
	return result
}

func matchesSearch(item map[string]any, searchLower string) bool {
	for _, value := range item {
		// Only search string-typed values
		if s, ok := value.(string); ok {
			if strings.Contains(strings.ToLower(s), searchLower) {
				return true
			}
		}
	}
	return false
}

func applyFields(items []map[string]any, fields []string) []map[string]any {
	result := make([]map[string]any, len(items))
	for i, item := range items {
		projected := make(map[string]any, len(fields))
		for _, field := range fields {
			if val, ok := item[field]; ok {
				projected[field] = val
			}
		}
		result[i] = projected
	}
	return result
}
