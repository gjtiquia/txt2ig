package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gjtiquia/txt2ig/internal/cli"
	"github.com/gjtiquia/txt2ig/internal/config"
)

func TestRunInit_CreatesDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-config.jsonc")

	cmd := &cli.InitCmd{
		Output: outputPath,
		Force:  false,
	}

	runInit(cmd)

	// Check file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %s", outputPath)
	}

	// Check file content matches default config
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	expectedContent := config.DefaultConfigContent()
	if string(content) != string(expectedContent) {
		t.Errorf("Config file content does not match default config")
	}
}

func TestRunInit_ErrorIfFileExists(t *testing.T) {
	t.Skip("Cannot test os.Exit behavior directly - requires refactoring runInit to return errors")
}

func TestRunInit_OverwriteWithForce(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "force-config.jsonc")

	// Create existing file with different content
	existingContent := []byte("existing content")
	if err := os.WriteFile(outputPath, existingContent, 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	cmd := &cli.InitCmd{
		Output: outputPath,
		Force:  true,
	}

	runInit(cmd)

	// Check file was overwritten
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	expectedContent := config.DefaultConfigContent()
	if string(content) != string(expectedContent) {
		t.Errorf("Config file was not overwritten with default config")
	}
}

func TestRunInit_DefaultOutputPath(t *testing.T) {
	// Create a temp directory and change to it
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cmd := &cli.InitCmd{
		Output: "", // Empty should default to .txt2igconfig.jsonc
		Force:  true,
	}

	runInit(cmd)

	// Check default file was created
	defaultPath := filepath.Join(tmpDir, ".txt2igconfig.jsonc")
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Errorf("Default config file was not created at .txt2igconfig.jsonc")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test existing file
	existingPath := filepath.Join(tmpDir, "existing.txt")
	os.WriteFile(existingPath, []byte("test"), 0644)
	if !fileExists(existingPath) {
		t.Errorf("fileExists should return true for existing file")
	}

	// Test non-existing file
	nonExistingPath := filepath.Join(tmpDir, "non-existing.txt")
	if fileExists(nonExistingPath) {
		t.Errorf("fileExists should return false for non-existing file")
	}
}

func TestRunCommand_InitRouting(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "init-test.jsonc")

	c := &cli.CLI{
		Init: cli.InitCmd{
			Output: outputPath,
			Force:  true,
		},
	}

	runCommand(c, "init")

	// Verify init command created the file
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Init command was not routed correctly - file not created")
	}
}

func TestRunCommand_ConvertRouting(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	c := &cli.CLI{
		Convert: cli.ConvertCmd{
			InputFile: testFile,
			Watch:     false,
		},
	}

	// This should NOT panic or fail
	// Note: runConvert will try to create image, but we're just testing routing
	runCommand(c, "convert")
}

func TestRunCommand_WatchRouting(t *testing.T) {
	t.Skip("Skipping: Watch mode starts server and blocks, requires integration test")
}

func TestRunCommand_WatchWithPortRouting(t *testing.T) {
	t.Skip("Skipping: Watch mode with port starts server and blocks, requires integration test")
}

func TestRunCommand_WebRouting(t *testing.T) {
	t.Skip("Skipping: Web mode starts server and blocks, requires integration test")
}

func TestRunCommand_DefaultRouting(t *testing.T) {
	t.Skip("Skipping: Requires input file or CLI refactoring to avoid os.Exit(1)")
}

func TestRunCommand_ImplicitConvertCommand(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Parse implicit command: "txt2ig test.md" (no "convert" subcommand)
	result, err := cli.Parse([]string{testFile})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Test that routing works correctly for implicit convert command
	runCommand(result.Cli, result.Command)
}

func TestRunCommand_ImplicitConvertWithWatch(t *testing.T) {
	t.Skip("Skipping: Watch mode starts server and blocks, requires integration test")
}

func TestRunCommand_ExplicitVsImplicit(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "compare.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Parse both explicit and implicit forms
	explicitResult, err := cli.Parse([]string{"convert", testFile})
	if err != nil {
		t.Fatalf("Parse explicit failed: %v", err)
	}

	implicitResult, err := cli.Parse([]string{testFile})
	if err != nil {
		t.Fatalf("Parse implicit failed: %v", err)
	}

	// Verify command routing is the same
	if explicitResult.Command != implicitResult.Command {
		t.Errorf("Command strings differ: explicit=%q, implicit=%q", explicitResult.Command, implicitResult.Command)
	}

	// Verify both route to convert
	explicitCmdName := strings.Fields(explicitResult.Command)[0]
	implicitCmdName := strings.Fields(implicitResult.Command)[0]

	if explicitCmdName != "convert" {
		t.Errorf("Explicit command name = %q, want 'convert'", explicitCmdName)
	}
	if implicitCmdName != "convert" {
		t.Errorf("Implicit command name = %q, want 'convert'", implicitCmdName)
	}
}
