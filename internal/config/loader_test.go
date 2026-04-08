package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoader_GetConfigPaths_WithoutCustomPath(t *testing.T) {
	loader := NewConfigLoader()
	paths := loader.GetConfigPaths()

	// Should return 3 default paths when no custom path is set
	if len(paths) != 3 {
		t.Errorf("GetConfigPaths should return 3 paths without custom path, got %d", len(paths))
	}

	// First path should be local config (relative)
	if paths[0] != "./.txt2igconfig.jsonc" {
		t.Errorf("First path should be local config, got %s", paths[0])
	}

	// Second and third paths should be absolute (XDG and home)
	for i := 1; i < len(paths); i++ {
		if !filepath.IsAbs(paths[i]) {
			t.Errorf("Path %d should be absolute: %s", i, paths[i])
		}
	}
}

func TestConfigLoader_GetConfigPaths_WithCustomPath(t *testing.T) {
	loader := NewConfigLoader()
	loader.SetCustomPath("/custom/config.jsonc")
	paths := loader.GetConfigPaths()

	// Should return 4 paths: custom + 3 defaults
	if len(paths) != 4 {
		t.Errorf("GetConfigPaths should return 4 paths with custom path, got %d", len(paths))
	}

	if paths[0] != "/custom/config.jsonc" {
		t.Errorf("First path should be custom path, got %s", paths[0])
	}
}

func TestConfigLoader_UsedPath_ReturnsEmptyByDefault(t *testing.T) {
	loader := NewConfigLoader()

	// Before Load() is called, UsedPath should return empty string
	if usedPath := loader.UsedPath(); usedPath != "" {
		t.Errorf("UsedPath should return empty string before Load(), got %s", usedPath)
	}
}

func TestConfigLoader_UsedPath_ReturnsPathAfterLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.jsonc")
	configContent := `{"fontSize": 24}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}

	loader := NewConfigLoader()
	loader.SetCustomPath(configPath)

	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify the config was loaded correctly
	if cfg.FontSize != 24 {
		t.Errorf("Expected FontSize 24, got %d", cfg.FontSize)
	}

	// Verify UsedPath returns the loaded path
	usedPath := loader.UsedPath()
	if usedPath != configPath {
		t.Errorf("UsedPath should return %s, got %s", configPath, usedPath)
	}
}

func TestConfigLoader_UsedPath_ReturnsEmptyWhenUsingDefaults(t *testing.T) {
	loader := NewConfigLoader()

	// Load without any custom path - will use defaults
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify defaults are used
	if cfg.FontSize != 32 {
		t.Errorf("Expected default FontSize 32, got %d", cfg.FontSize)
	}

	// UsedPath should be empty when using defaults
	usedPath := loader.UsedPath()
	if usedPath != "" {
		t.Errorf("UsedPath should return empty string when using defaults, got %s", usedPath)
	}
}

func TestConfigLoader_Load_PreferCustomPath(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()

	// Custom config
	customConfigPath := filepath.Join(tmpDir, "custom.jsonc")
	customConfig := `{"fontSize": 100}`
	if err := os.WriteFile(customConfigPath, []byte(customConfig), 0644); err != nil {
		t.Fatalf("Failed to write custom config: %v", err)
	}

	// Local config (should be ignored when custom is provided)
	localConfigPath := filepath.Join(tmpDir, ".txt2igconfig.jsonc")
	localConfig := `{"fontSize": 50}`
	if err := os.WriteFile(localConfigPath, []byte(localConfig), 0644); err != nil {
		t.Fatalf("Failed to write local config: %v", err)
	}

	// Change working directory to temp dir
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	loader := NewConfigLoader()
	loader.SetCustomPath(customConfigPath)

	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Should use custom config (fontSize: 100), not local config (fontSize: 50)
	if cfg.FontSize != 100 {
		t.Errorf("Expected FontSize 100 from custom config, got %d", cfg.FontSize)
	}
}