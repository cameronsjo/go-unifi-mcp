package mcpgen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindFieldsDir(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		setup     func() string
		wantMatch string
		wantErr   bool
	}{
		{
			name: "finds versioned subdirectory",
			setup: func() string {
				baseDir := filepath.Join(tmpDir, "test1")
				vDir := filepath.Join(baseDir, "v9.0.114")
				if err := os.MkdirAll(vDir, 0755); err != nil {
					t.Fatal(err)
				}
				return baseDir
			},
			wantMatch: "v9.0.114",
			wantErr:   false,
		},
		{
			name: "returns base dir if no versioned subdirectory",
			setup: func() string {
				baseDir := filepath.Join(tmpDir, "test2")
				if err := os.MkdirAll(baseDir, 0755); err != nil {
					t.Fatal(err)
				}
				// Create a non-versioned file
				if err := os.WriteFile(filepath.Join(baseDir, "test.json"), []byte("{}"), 0644); err != nil {
					t.Fatal(err)
				}
				return baseDir
			},
			wantMatch: "test2",
			wantErr:   false,
		},
		{
			name: "returns error for non-existent directory",
			setup: func() string {
				return filepath.Join(tmpDir, "nonexistent")
			},
			wantMatch: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseDir := tt.setup()
			got, err := findFieldsDir(baseDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("findFieldsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !strings.Contains(got, tt.wantMatch) {
				t.Errorf("findFieldsDir() = %v, want to contain %v", got, tt.wantMatch)
			}
		})
	}
}

func TestGenerate_Integration(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check that field definitions exist
	fieldsDir := ".tmp/fields"
	if _, err := os.Stat("../../" + fieldsDir); os.IsNotExist(err) {
		t.Skip("field definitions not downloaded, run 'task download-fields' first")
	}

	// Create temporary output directory
	outDir := t.TempDir()

	cfg := GeneratorConfig{
		FieldsDir: "../../.tmp/fields",
		V2Dir:     "../../internal/gounifi/v2",
		OutDir:    outDir,
	}

	// Generate
	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify files were created
	toolsFile := filepath.Join(outDir, "tools.gen.go")
	registryFile := filepath.Join(outDir, "registry.gen.go")

	if _, err := os.Stat(toolsFile); os.IsNotExist(err) {
		t.Errorf("tools.gen.go was not created")
	}

	if _, err := os.Stat(registryFile); os.IsNotExist(err) {
		t.Errorf("registry.gen.go was not created")
	}

	// Check that the generated files have content
	toolsContent, err := os.ReadFile(toolsFile)
	if err != nil {
		t.Fatalf("Failed to read tools.gen.go: %v", err)
	}
	if len(toolsContent) < 1000 {
		t.Errorf("tools.gen.go seems too small: %d bytes", len(toolsContent))
	}

	registryContent, err := os.ReadFile(registryFile)
	if err != nil {
		t.Fatalf("Failed to read registry.gen.go: %v", err)
	}
	if len(registryContent) < 1000 {
		t.Errorf("registry.gen.go seems too small: %d bytes", len(registryContent))
	}

	// Verify expected content
	if !strings.Contains(string(toolsContent), "package generated") {
		t.Error("tools.gen.go missing package declaration")
	}
	if !strings.Contains(string(registryContent), "RegisterAllTools") {
		t.Error("registry.gen.go missing RegisterAllTools function")
	}

	// Verify generated code compiles
	cmd := exec.Command("go", "build", toolsFile, registryFile)
	cmd.Env = append(os.Environ(),
		"GOPATH="+filepath.Join(os.Getenv("PWD"), "../../.go"),
		"GOMODCACHE="+filepath.Join(os.Getenv("PWD"), "../../.go/pkg/mod"),
		"GOCACHE="+filepath.Join(os.Getenv("PWD"), "../../.go/cache"),
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("Generated code does not compile: %v\n%s", err, output)
	}
}

func TestGenerate_ToolCounts(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check that field definitions exist
	if _, err := os.Stat("../../.tmp/fields"); os.IsNotExist(err) {
		t.Skip("field definitions not downloaded, run 'task download-fields' first")
	}

	outDir := t.TempDir()

	cfg := GeneratorConfig{
		FieldsDir: "../../.tmp/fields",
		V2Dir:     "../../internal/gounifi/v2",
		OutDir:    outDir,
	}

	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Read the generated tools file
	toolsFile := filepath.Join(outDir, "tools.gen.go")
	content, err := os.ReadFile(toolsFile)
	if err != nil {
		t.Fatalf("Failed to read tools.gen.go: %v", err)
	}

	// Count function definitions - should be at least 200
	funcCount := strings.Count(string(content), "\nfunc ")
	if funcCount < 200 {
		t.Errorf("Expected at least 200 tool functions, got %d", funcCount)
	}
}

func TestGenerate_MissingFieldsDir(t *testing.T) {
	outDir := t.TempDir()

	cfg := GeneratorConfig{
		FieldsDir: "/nonexistent/path",
		V2Dir:     "../../internal/gounifi/v2",
		OutDir:    outDir,
	}

	err := Generate(cfg)
	if err == nil {
		t.Error("Generate() should return error for missing fields dir")
	}
}

func TestGenerate_MissingV2Dir(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("skipping test requiring field definitions in short mode")
	}

	// Check that field definitions exist
	if _, err := os.Stat("../../.tmp/fields"); os.IsNotExist(err) {
		t.Skip("field definitions not downloaded, run 'task download-fields' first")
	}

	outDir := t.TempDir()

	cfg := GeneratorConfig{
		FieldsDir: "../../.tmp/fields",
		V2Dir:     "/nonexistent/path",
		OutDir:    outDir,
	}

	err := Generate(cfg)
	if err == nil {
		t.Error("Generate() should return error for missing V2 dir")
	}
}

func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	outDir := t.TempDir()
	outPath := filepath.Join(outDir, "test.go")

	// Test with a template that doesn't exist
	err := renderTemplate("templates/nonexistent.tmpl", outPath, nil)
	if err == nil {
		t.Error("renderTemplate() should return error for non-existent template")
	}
	if !strings.Contains(err.Error(), "failed to read template") {
		t.Errorf("Expected 'failed to read template' error, got: %v", err)
	}
}

func TestRenderTemplate_WriteError(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Create a readonly directory
	tmpDir := t.TempDir()
	readonlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readonlyDir, 0555); err != nil {
		t.Skip("failed to create readonly directory")
	}
	defer func() { _ = os.Chmod(readonlyDir, 0755) }()

	outPath := filepath.Join(readonlyDir, "test.go")

	// Use a valid template with simple data
	data := []ToolInfo{
		{Name: "Test", SnakeName: "test", Operations: []string{"Get"}, IsSetting: false},
	}

	err := renderTemplate("templates/tools.go.tmpl", outPath, data)
	if err == nil {
		t.Error("renderTemplate() should return error when write fails")
	}
}

func TestGenerate_OutputDirCreationFailure(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Check that field definitions exist
	if _, err := os.Stat("../../.tmp/fields"); os.IsNotExist(err) {
		t.Skip("field definitions not downloaded, run 'task download-fields' first")
	}

	// Use a path that will fail to create (inside a file)
	tmpFile := filepath.Join(t.TempDir(), "afile")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := GeneratorConfig{
		FieldsDir: "../../.tmp/fields",
		V2Dir:     "../../internal/gounifi/v2",
		OutDir:    filepath.Join(tmpFile, "subdir"), // Can't create dir inside a file
	}

	err := Generate(cfg)
	if err == nil {
		t.Error("Generate() should return error when output dir creation fails")
	}
}

func TestToolInfoSnakeName(t *testing.T) {
	tests := []struct {
		name      string
		toolName  string
		wantSnake string
	}{
		{
			name:      "simple name",
			toolName:  "Network",
			wantSnake: "network",
		},
		{
			name:      "camel case",
			toolName:  "FirewallRule",
			wantSnake: "firewall_rule",
		},
		{
			name:      "setting",
			toolName:  "SettingMgmt",
			wantSnake: "setting_mgmt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := ToolInfo{
				Name:      tt.toolName,
				SnakeName: tt.wantSnake,
			}

			if tool.SnakeName != tt.wantSnake {
				t.Errorf("SnakeName = %v, want %v", tool.SnakeName, tt.wantSnake)
			}
		})
	}
}
