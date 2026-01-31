package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zdzh/goctl-env-template/internal/generator"
	"github.com/zdzh/goctl-env-template/internal/parser"
)

var (
	configFile = flag.String("c", "config/config.go", "Config file path")
	outputFile = flag.String("o", ".env.template", "Output file path")
)

func main() {
	flag.Parse()

	configPath := *configFile
	outputPath := *outputFile

	if err := generateEnvTemplate(configPath, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated .env template: %s\n", outputPath)
}

func generateEnvTemplate(configPath, outputPath string) error {
	p, err := parser.NewParser(configPath)
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	structs, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if len(structs) == 0 {
		return fmt.Errorf("no config structs found with environment variable tags")
	}

	gen := generator.NewGenerator(structs)
	content := gen.Generate()

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
