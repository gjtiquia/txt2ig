package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gjtiquia/txt2ig/internal/config"
)

func TestRunInit_CreatesDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-config.jsonc")

	cmd := &InitCmd{
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

	cmd := &InitCmd{
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

	cmd := &InitCmd{
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