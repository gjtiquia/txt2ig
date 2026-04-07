package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gjtiquia/txt2ig/internal/cli"
	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
)

func main() {
	// Parse CLI arguments
	cliArgs, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load config
	loader := config.NewConfigLoader()
	if cliArgs.Config != "" {
		loader.SetCustomPath(cliArgs.Config)
	}

	cfg, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Determine output path
	outputPath := renderer.DetermineOutputPath(cliArgs.InputFile, cliArgs.Output)

	// Create renderer
	r := renderer.NewRenderer(cfg)
	defer r.Close()

	// Render the image
	fmt.Printf("Converting %s to %s...\n", filepath.Base(cliArgs.InputFile), filepath.Base(outputPath))
	if err := r.Render(cliArgs.InputFile, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s\n", outputPath)
}
