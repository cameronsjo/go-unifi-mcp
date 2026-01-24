// mcpgen generates MCP tool handlers from UniFi API field definitions.
package main

import (
	"flag"
	"log"

	"github.com/claytono/go-unifi-mcp/internal/mcpgen"
)

func main() {
	fieldsDir := flag.String("fields", ".tmp/fields", "Path to v1 field definitions")
	v2Dir := flag.String("v2", "internal/gounifi/v2", "Path to v2 field definitions")
	outDir := flag.String("out", "internal/tools/generated", "Output directory")
	flag.Parse()

	cfg := mcpgen.GeneratorConfig{
		FieldsDir: *fieldsDir,
		V2Dir:     *v2Dir,
		OutDir:    *outDir,
	}

	if err := mcpgen.Generate(cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("Generated MCP tools to %s", *outDir)
}
