package main

import (
	"fmt"
	"os"

	"github.com/gjtiquia/txt2ig/internal/cli"
	"github.com/gjtiquia/txt2ig/internal/config"
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

	// TODO: Implement the rest of the pipeline
	// 1. Read input file
	// 2. Run pre-processors
	// 3. Load fonts
	// 4. Render image
	// 5. Run post-processors
	// 6. Save output

	fmt.Printf("Input: %s\n", cliArgs.InputFile)
	fmt.Printf("Output: %s\n", cliArgs.Output)
	fmt.Printf("Config: %+v\n", cfg)
}
