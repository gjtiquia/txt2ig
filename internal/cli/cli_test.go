package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_InitCommand(t *testing.T) {
	result, err := Parse([]string{"init"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedName := ".txt2igconfig.jsonc"
	if filepath.Base(result.Cli.Init.Output) != expectedName {
		t.Errorf("Init.Output basename = %s, want %s", filepath.Base(result.Cli.Init.Output), expectedName)
	}
	if result.Cli.Init.Force != false {
		t.Errorf("Init.Force = %v, want false", result.Cli.Init.Force)
	}
	if result.Command != "init" {
		t.Errorf("Command = %s, want init", result.Command)
	}
}

func TestParse_InitWithCustomOutput(t *testing.T) {
	result, err := Parse([]string{"init", "-o", "custom.jsonc"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedName := "custom.jsonc"
	if filepath.Base(result.Cli.Init.Output) != expectedName {
		t.Errorf("Init.Output basename = %s, want %s", filepath.Base(result.Cli.Init.Output), expectedName)
	}
}

func TestParse_InitWithForce(t *testing.T) {
	result, err := Parse([]string{"init", "--force"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Init.Force != true {
		t.Errorf("Init.Force = %v, want true", result.Cli.Init.Force)
	}
}

func TestParse_ConvertCommand(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{"convert", testFile})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if result.Cli.Convert.Output != "" {
		t.Errorf("Convert.Output = %s, want empty", result.Cli.Convert.Output)
	}
	if result.Cli.Convert.Watch != false {
		t.Errorf("Convert.Watch = %v, want false", result.Cli.Convert.Watch)
	}
	// Command is "convert <file>" not just "convert"
	if result.Command == "" {
		t.Errorf("Command should not be empty")
	}
}

func TestParse_ConvertWithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{"convert", testFile, "-o", "output.jpg"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	expectedName := "output.jpg"
	if filepath.Base(result.Cli.Convert.Output) != expectedName {
		t.Errorf("Convert.Output basename = %s, want %s", filepath.Base(result.Cli.Convert.Output), expectedName)
	}
}

func TestParse_ConvertWithWatchFlag(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{"convert", testFile, "-w"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if result.Cli.Convert.Watch != true {
		t.Errorf("Convert.Watch = %v, want true", result.Cli.Convert.Watch)
	}
}

func TestParse_ConvertWithWatchAndPort(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{"convert", testFile, "-w", "-p", "3000"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if result.Cli.Convert.Watch != true {
		t.Errorf("Convert.Watch = %v, want true", result.Cli.Convert.Watch)
	}
	if result.Cli.Convert.Port != 3000 {
		t.Errorf("Convert.Port = %d, want 3000", result.Cli.Convert.Port)
	}
}

func TestParse_ConvertWithDebugFlag(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{"convert", testFile, "--debug"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if result.Cli.Convert.Debug != true {
		t.Errorf("Convert.Debug = %v, want true", result.Cli.Convert.Debug)
	}
}

func TestParse_ConvertWithConfig(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	configFile := filepath.Join(tmpDir, "custom.jsonc")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(configFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	result, err := Parse([]string{"convert", testFile, "-c", configFile})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if result.Cli.Convert.Config != configFile {
		t.Errorf("Convert.Config = %s, want %s", result.Cli.Convert.Config, configFile)
	}
}

func TestParse_WebCommand(t *testing.T) {
	result, err := Parse([]string{"web"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Web.Port != 3000 {
		t.Errorf("Web.Port = %d, want 3000 (default)", result.Cli.Web.Port)
	}
	if result.Command != "web" {
		t.Errorf("Command = %s, want web", result.Command)
	}
}

func TestParse_WebWithCustomPort(t *testing.T) {
	result, err := Parse([]string{"web", "-p", "8080"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Web.Port != 8080 {
		t.Errorf("Web.Port = %d, want 8080", result.Cli.Web.Port)
	}
}
