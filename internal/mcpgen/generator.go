// Package mcpgen generates MCP tool handlers from UniFi API field definitions.
package mcpgen

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/claytono/go-unifi-mcp/internal/gounifi"
	"github.com/iancoleman/strcase"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// ToolInfo contains metadata about a tool to be generated.
type ToolInfo struct {
	Name       string   // e.g., "Network"
	SnakeName  string   // e.g., "network"
	Operations []string // e.g., ["List", "Get", "Create", "Update", "Delete"]
	IsSetting  bool
	IsV2       bool
}

// GeneratorConfig holds configuration for the generator.
type GeneratorConfig struct {
	FieldsDir string // Path to v1 field JSONs
	V2Dir     string // Path to v2 field JSONs
	OutDir    string // Output directory
}

// Generate generates MCP tool handlers from UniFi API field definitions.
func Generate(cfg GeneratorConfig) error {
	// Load customizer for field processing
	customizer, err := gounifi.NewCodeCustomizer("")
	if err != nil {
		return fmt.Errorf("failed to create customizer: %w", err)
	}

	// Find the versioned fields directory (e.g., .tmp/fields/v9.0.114)
	fieldsDir, err := findFieldsDir(cfg.FieldsDir)
	if err != nil {
		return fmt.Errorf("failed to find fields directory: %w", err)
	}

	// Parse v1 resources from downloaded fields
	v1Resources, err := gounifi.BuildResourcesFromDownloadedFields(fieldsDir, *customizer, false)
	if err != nil {
		return fmt.Errorf("failed to parse v1 resources: %w", err)
	}

	// Parse v2 resources from internal/gounifi/v2/
	v2Resources, err := gounifi.BuildResourcesFromDownloadedFields(cfg.V2Dir, *customizer, true)
	if err != nil {
		return fmt.Errorf("failed to parse v2 resources: %w", err)
	}

	// Combine and process resources
	allResources := make([]*gounifi.Resource, 0, len(v1Resources)+len(v2Resources))
	allResources = append(allResources, v1Resources...)
	allResources = append(allResources, v2Resources...)
	tools := make([]ToolInfo, 0, len(allResources))

	for _, r := range allResources {
		if customizer.IsExcludedFromClient(r.Name()) {
			continue
		}

		tool := ToolInfo{
			Name:       r.StructName,
			SnakeName:  strcase.ToSnake(r.StructName),
			IsSetting:  r.IsSetting(),
			IsV2:       r.IsV2(),
			Operations: InferOperations(r),
		}
		tools = append(tools, tool)
	}

	// Sort for deterministic output
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})

	// Ensure output directory exists
	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Render templates
	if err := renderTemplate("templates/tools.go.tmpl", filepath.Join(cfg.OutDir, "tools.gen.go"), tools); err != nil {
		return fmt.Errorf("failed to render tools template: %w", err)
	}

	if err := renderTemplate("templates/registry.go.tmpl", filepath.Join(cfg.OutDir, "registry.gen.go"), tools); err != nil {
		return fmt.Errorf("failed to render registry template: %w", err)
	}

	return nil
}

// findFieldsDir finds the versioned subdirectory in the fields directory.
func findFieldsDir(baseDir string) (string, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 0 && entry.Name()[0] == 'v' {
			return filepath.Join(baseDir, entry.Name()), nil
		}
	}

	// If no versioned directory, assume files are directly in baseDir
	return baseDir, nil
}

// renderTemplate renders a template to a file.
func renderTemplate(templatePath, outputPath string, data interface{}) error {
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	funcMap := template.FuncMap{
		"has": func(needle string, haystack []string) bool {
			for _, s := range haystack {
				if s == needle {
					return true
				}
			}
			return false
		},
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	// Format the generated Go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Write unformatted for debugging
		_ = os.WriteFile(outputPath+".unformatted", buf.Bytes(), 0644)
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outputPath, err)
	}

	return nil
}
