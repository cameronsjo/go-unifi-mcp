package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		expected Options
	}{
		{
			name:     "empty args",
			args:     map[string]any{},
			expected: Options{},
		},
		{
			name:     "nil args",
			args:     nil,
			expected: Options{},
		},
		{
			name: "wrong types ignored",
			args: map[string]any{
				"filter": "not an object",
				"search": 123,
				"fields": "not an array",
			},
			expected: Options{},
		},
		{
			name: "valid filter",
			args: map[string]any{
				"filter": map[string]any{"name": "test"},
			},
			expected: Options{
				Filter: map[string]any{"name": "test"},
			},
		},
		{
			name: "valid search",
			args: map[string]any{
				"search": "amazon",
			},
			expected: Options{
				Search: "amazon",
			},
		},
		{
			name: "valid fields",
			args: map[string]any{
				"fields": []any{"name", "ip"},
			},
			expected: Options{
				Fields: []string{"name", "ip"},
			},
		},
		{
			name: "fields with non-string elements skipped",
			args: map[string]any{
				"fields": []any{"name", 123, "ip"},
			},
			expected: Options{
				Fields: []string{"name", "ip"},
			},
		},
		{
			name: "all params together",
			args: map[string]any{
				"filter": map[string]any{"type": "uap"},
				"search": "living",
				"fields": []any{"name", "ip"},
			},
			expected: Options{
				Filter: map[string]any{"type": "uap"},
				Search: "living",
				Fields: []string{"name", "ip"},
			},
		},
		{
			name: "extra args ignored",
			args: map[string]any{
				"site":   "default",
				"search": "test",
			},
			expected: Options{
				Search: "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseOptions(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasQuery(t *testing.T) {
	tests := []struct {
		name     string
		opts     Options
		expected bool
	}{
		{
			name:     "zero value",
			opts:     Options{},
			expected: false,
		},
		{
			name:     "filter only",
			opts:     Options{Filter: map[string]any{"name": "x"}},
			expected: true,
		},
		{
			name:     "search only",
			opts:     Options{Search: "x"},
			expected: true,
		},
		{
			name:     "fields only",
			opts:     Options{Fields: []string{"name"}},
			expected: true,
		},
		{
			name:     "empty filter is false",
			opts:     Options{Filter: map[string]any{}},
			expected: false,
		},
		{
			name:     "empty fields is false",
			opts:     Options{Fields: []string{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.opts.HasQuery())
		})
	}
}

var testItems = []map[string]any{
	{"name": "switch-1", "ip": "10.0.0.1", "type": "usw", "port_count": float64(24)},
	{"name": "ap-living-room", "ip": "10.0.0.2", "type": "uap", "port_count": float64(1)},
	{"name": "amazon-echo", "ip": "10.0.0.3", "type": "uap", "port_count": float64(0)},
}

func TestApply_NoOptions(t *testing.T) {
	result := Apply(testItems, Options{})
	assert.Equal(t, testItems, result)
}

func TestApply_EmptyInput(t *testing.T) {
	result := Apply([]map[string]any{}, Options{Search: "anything"})
	assert.Empty(t, result)
}

func TestApply_NilInput(t *testing.T) {
	result := Apply(nil, Options{Search: "anything"})
	assert.Nil(t, result)
}

func TestApply_FilterNoMatch_ReturnsEmptySlice(t *testing.T) {
	result := Apply(testItems, Options{Filter: map[string]any{"type": "nonexistent"}})
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestApply_SearchNoMatch_ReturnsEmptySlice(t *testing.T) {
	result := Apply(testItems, Options{Search: "nonexistent"})
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestApply_FilterExact(t *testing.T) {
	tests := []struct {
		name     string
		filter   map[string]any
		expected int
	}{
		{
			name:     "exact string match",
			filter:   map[string]any{"type": "usw"},
			expected: 1,
		},
		{
			name:     "exact number match",
			filter:   map[string]any{"port_count": float64(24)},
			expected: 1,
		},
		{
			name:     "exact match no results",
			filter:   map[string]any{"type": "ugw"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Apply(testItems, Options{Filter: tt.filter})
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestApply_FilterExactBool(t *testing.T) {
	items := []map[string]any{
		{"name": "a", "enabled": true},
		{"name": "b", "enabled": false},
	}
	result := Apply(items, Options{Filter: map[string]any{"enabled": true}})
	assert.Len(t, result, 1)
	assert.Equal(t, "a", result[0]["name"])
}

func TestApply_FilterContains(t *testing.T) {
	result := Apply(testItems, Options{
		Filter: map[string]any{
			"name": map[string]any{"contains": "amazon"},
		},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "amazon-echo", result[0]["name"])
}

func TestApply_FilterContainsCaseInsensitive(t *testing.T) {
	result := Apply(testItems, Options{
		Filter: map[string]any{
			"name": map[string]any{"contains": "LIVING"},
		},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "ap-living-room", result[0]["name"])
}

func TestApply_FilterRegex(t *testing.T) {
	result := Apply(testItems, Options{
		Filter: map[string]any{
			"name": map[string]any{"regex": "^ap-.*"},
		},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "ap-living-room", result[0]["name"])
}

func TestApply_FilterRegexInvalid(t *testing.T) {
	// Invalid regex should gracefully not match (no error)
	result := Apply(testItems, Options{
		Filter: map[string]any{
			"name": map[string]any{"regex": "[invalid"},
		},
	})
	assert.Empty(t, result)
}

func TestApply_FilterMultipleFields(t *testing.T) {
	// AND logic: must match both conditions
	result := Apply(testItems, Options{
		Filter: map[string]any{
			"type": "uap",
			"name": map[string]any{"contains": "echo"},
		},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "amazon-echo", result[0]["name"])
}

func TestApply_FilterNonExistentField(t *testing.T) {
	result := Apply(testItems, Options{
		Filter: map[string]any{"nonexistent": "value"},
	})
	assert.Empty(t, result)
}

func TestApply_FilterNilFieldValue(t *testing.T) {
	items := []map[string]any{
		{"name": "a", "ip": nil},
		{"name": "b", "ip": "10.0.0.1"},
	}
	result := Apply(items, Options{
		Filter: map[string]any{"ip": "10.0.0.1"},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, "b", result[0]["name"])
}

func TestApply_FilterUnknownOperator(t *testing.T) {
	// Unknown operator object should not match
	result := Apply(testItems, Options{
		Filter: map[string]any{
			"name": map[string]any{"startsWith": "ap"},
		},
	})
	assert.Empty(t, result)
}

func TestApply_Search(t *testing.T) {
	result := Apply(testItems, Options{Search: "amazon"})
	assert.Len(t, result, 1)
	assert.Equal(t, "amazon-echo", result[0]["name"])
}

func TestApply_SearchCaseInsensitive(t *testing.T) {
	result := Apply(testItems, Options{Search: "SWITCH"})
	assert.Len(t, result, 1)
	assert.Equal(t, "switch-1", result[0]["name"])
}

func TestApply_SearchNoMatch(t *testing.T) {
	result := Apply(testItems, Options{Search: "nonexistent"})
	assert.Empty(t, result)
}

func TestApply_SearchMatchesAnyStringField(t *testing.T) {
	// Search should match across any string field
	result := Apply(testItems, Options{Search: "10.0.0.2"})
	assert.Len(t, result, 1)
	assert.Equal(t, "ap-living-room", result[0]["name"])
}

func TestApply_SearchSkipsNonStringValues(t *testing.T) {
	// "24" appears in port_count (float64) but search only checks strings
	result := Apply(testItems, Options{Search: "24"})
	assert.Empty(t, result)
}

func TestApply_Fields(t *testing.T) {
	result := Apply(testItems, Options{Fields: []string{"name", "ip"}})
	assert.Len(t, result, 3)
	for _, item := range result {
		assert.Len(t, item, 2)
		assert.Contains(t, item, "name")
		assert.Contains(t, item, "ip")
	}
}

func TestApply_FieldsNonExistentSilentlyIgnored(t *testing.T) {
	result := Apply(testItems, Options{Fields: []string{"name", "nonexistent"}})
	assert.Len(t, result, 3)
	for _, item := range result {
		assert.Len(t, item, 1)
		assert.Contains(t, item, "name")
	}
}

func TestApply_Combined(t *testing.T) {
	result := Apply(testItems, Options{
		Filter: map[string]any{"type": "uap"},
		Search: "echo",
		Fields: []string{"name", "ip"},
	})
	assert.Len(t, result, 1)
	assert.Equal(t, map[string]any{"name": "amazon-echo", "ip": "10.0.0.3"}, result[0])
}
