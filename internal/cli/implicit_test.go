package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_ImplicitConvertCommand(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{testFile})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if result.Cli.Convert.Watch != false {
		t.Errorf("Convert.Watch = %v, want false", result.Cli.Convert.Watch)
	}
	if result.Cli.Convert.Output != "" {
		t.Errorf("Convert.Output = %s, want empty", result.Cli.Convert.Output)
	}
}

func TestParse_ImplicitConvertWithWatch(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{testFile, "-w"})
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

func TestParse_ImplicitConvertWithWatchAndPort(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{testFile, "-w", "-p", "3000"})
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

func TestParse_ImplicitConvertWithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{testFile, "-o", "output.jpg"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Cli.Convert.InputFile != testFile {
		t.Errorf("Convert.InputFile = %s, want %s", result.Cli.Convert.InputFile, testFile)
	}
	if filepath.Base(result.Cli.Convert.Output) != "output.jpg" {
		t.Errorf("Convert.Output basename = %s, want output.jpg", filepath.Base(result.Cli.Convert.Output))
	}
}

func TestParse_ImplicitConvertWithDebug(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := Parse([]string{testFile, "--debug"})
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

func TestParse_CommandStringEquals_ExplicitVsImplicit(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	explicitResult, err := Parse([]string{"convert", testFile})
	if err != nil {
		t.Fatalf("Parse explicit failed: %v", err)
	}

	implicitResult, err := Parse([]string{testFile})
	if err != nil {
		t.Fatalf("Parse implicit failed: %v", err)
	}

	if explicitResult.Command != implicitResult.Command {
		t.Errorf("Command strings differ: explicit=%q, implicit=%q", explicitResult.Command, implicitResult.Command)
	}
}
