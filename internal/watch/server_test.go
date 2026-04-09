package watch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gjtiquia/txt2ig/internal/config"
)

func TestNewServer_TracksUsedConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	configFile := filepath.Join(tmpDir, "custom.jsonc")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := os.WriteFile(configFile, []byte(`{"fontSize": 40}`), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	server, err := NewServer(testFile, configFile)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	if server.usedConfigPath != configFile {
		t.Errorf("usedConfigPath = %s, want %s", server.usedConfigPath, configFile)
	}

	if server.config.FontSize != 40 {
		t.Errorf("config.FontSize = %d, want 40", server.config.FontSize)
	}
}

func TestNewServer_WithoutConfig(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	server, err := NewServer(testFile, "")
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	if server.usedConfigPath != "" {
		t.Errorf("usedConfigPath = %s, want empty string (embedded defaults)", server.usedConfigPath)
	}

	defaultConfig := config.DefaultConfig()
	if server.config.FontSize != defaultConfig.FontSize {
		t.Errorf("config.FontSize = %d, want %d (default)", server.config.FontSize, defaultConfig.FontSize)
	}
}

func TestNewServer_LocalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	localConfig := filepath.Join(tmpDir, ".txt2igconfig.jsonc")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := os.WriteFile(localConfig, []byte(`{"fontSize": 50}`), 0644); err != nil {
		t.Fatalf("Failed to create local config: %v", err)
	}

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	server, err := NewServer(testFile, "")
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// ConfigLoader returns relative paths for local config
	// Compare basenames instead of full paths
	if filepath.Base(server.usedConfigPath) != ".txt2igconfig.jsonc" {
		t.Errorf("usedConfigPath basename = %s, want .txt2igconfig.jsonc", filepath.Base(server.usedConfigPath))
	}

	if server.config.FontSize != 50 {
		t.Errorf("config.FontSize = %d, want 50", server.config.FontSize)
	}
}
