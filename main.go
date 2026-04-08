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
	cliArgs, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	loader := config.NewConfigLoader()
	if cliArgs.Config != "" {
		loader.SetCustomPath(cliArgs.Config)
	}

	cfg, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliArgs.Debug {
		printDebugInfo(loader, cfg)
		os.Exit(0)
	}

	if cliArgs.InputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input file is required\n")
		os.Exit(1)
	}

	outputPath := renderer.DetermineOutputPath(cliArgs.InputFile, cliArgs.Output)

	r := renderer.NewRenderer(cfg)
	defer r.Close()

	fmt.Printf("Converting %s to %s...\n", filepath.Base(cliArgs.InputFile), filepath.Base(outputPath))
	if err := r.Render(cliArgs.InputFile, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s\n", outputPath)
}

func printDebugInfo(loader *config.ConfigLoader, cfg *config.Config) {
	fmt.Println("Config search chain:")
	fmt.Println()

	paths := loader.GetConfigPaths()
	hasCustomConfig := len(paths) == 4
	pathIdx := 0

	homeDir, _ := os.UserHomeDir()
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(homeDir, ".config")
	}

	labels := []string{"--config flag", "./.txt2igconfig.jsonc", "$XDG_CONFIG_HOME/txt2ig/config.jsonc", "~/.txt2ig/config.jsonc"}

	for i := 0; i < 5; i++ {
		if i < 4 {
			if i == 0 && !hasCustomConfig {
				fmt.Printf("  %d. %s: (not specified)\n", i+1, labels[i])
				continue
			}
			if pathIdx < len(paths) {
				path := paths[pathIdx]
				fmt.Printf("  %d. %s (%s): ", i+1, labels[i], path)
				if _, err := os.Stat(path); err == nil {
					fmt.Println("found ✓")
				} else {
					fmt.Println("not found")
				}
				pathIdx++
			}
		} else {
			fmt.Println("  5. embedded defaults: (used if no config file found)")
		}
	}

	fmt.Println()

	usedPath := loader.UsedPath()
	if usedPath == "" {
		fmt.Println("Using: embedded defaults")
	} else {
		fmt.Printf("Using: %s\n", usedPath)
	}

	fmt.Println()

	jsonOutput, err := cfg.ToJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config:")
	fmt.Println(string(jsonOutput))
}