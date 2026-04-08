package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/gjtiquia/txt2ig/internal/config"
	"github.com/gjtiquia/txt2ig/internal/renderer"
)

var cli struct {
	Convert ConvertCmd `cmd:"" default:"withargs" help:"Convert text file to image"`
	Init    InitCmd    `cmd:"" help:"Create default config file"`
}

type ConvertCmd struct {
	InputFile string `arg:"" name:"file" help:"Text file to convert" type:"existingfile" optional:""`
	Output    string `short:"o" help:"Output file name (.jpg or .png)" type:"path"`
	Config    string `short:"c" long:"config" help:"Custom config file" type:"existingfile"`
	Debug     bool   `short:"d" long:"debug" help:"Print config info and exit"`
}

type InitCmd struct {
	Output string `short:"o" long:"output" help:"Output file path" type:"path" default:".txt2igconfig.jsonc"`
	Force  bool   `short:"f" long:"force" help:"Overwrite existing file"`
}

func main() {
	ctx := kong.Parse(&cli, kong.Name("txt2ig"), kong.Description("Convert text files to images for Instagram"))

	switch ctx.Command() {
	case "init":
		runInit(&cli.Init)
	case "convert":
		runConvert(&cli.Convert)
	default:
		runConvert(&cli.Convert)
	}
}

func runConvert(cmd *ConvertCmd) {
	loader := config.NewConfigLoader()
	if cmd.Config != "" {
		loader.SetCustomPath(cmd.Config)
	}

	cfg, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cmd.Debug {
		printDebugInfo(loader, cfg)
		os.Exit(0)
	}

	if cmd.InputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input file is required\n")
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

func runInit(cmd *InitCmd) {
	outputPath := cmd.Output
	if outputPath == "" {
		outputPath = ".txt2igconfig.jsonc"
	}

	if _, err := os.Stat(outputPath); err == nil && !cmd.Force {
		fmt.Fprintf(os.Stderr, "Error: %s already exists (use --force to overwrite)\n", outputPath)
		os.Exit(1)
	}

	content := config.DefaultConfigContent()
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		os.Exit(1)
	}

	if cmd.Force {
		fmt.Printf("Overwrote %s\n", outputPath)
	} else {
		fmt.Printf("Created %s\n", outputPath)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func printDebugInfo(loader *config.ConfigLoader, cfg *config.Config) {
	fmt.Println("Config search chain:")
	fmt.Println()

	paths := loader.GetConfigPaths()
	hasCustomConfig := len(paths) == 4
	pathIdx := 0

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