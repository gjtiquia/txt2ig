package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gjtiquia/txt2ig/internal/cli"
	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
	"github.com/gjtiquia/txt2ig/internal/web"
)

func main() {
	cliArgs, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch {
	case cliArgs.Convert.InputFile != "":
		runConvert(&cliArgs.Convert)
	case true:
		runWeb(&cliArgs.Web)
	}
}

func runConvert(cmd *cli.ConvertCmd) {
	loader := config.NewConfigLoader()
	if cmd.Config != "" {
		loader.SetCustomPath(cmd.Config)
	}

	cfg, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	outputPath := renderer.DetermineOutputPath(cmd.InputFile, cmd.Output)

	r := renderer.NewRenderer(cfg)
	defer r.Close()

	fmt.Printf("Converting %s to %s...\n", filepath.Base(cmd.InputFile), filepath.Base(outputPath))
	if err := r.Render(cmd.InputFile, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s\n", outputPath)
}

func runWeb(cmd *cli.WebCmd) {
	server := web.NewServer()
	if err := server.Run(cmd.Port); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
		os.Exit(1)
	}
}
