package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type ConfigLoader struct {
	customPath string
	usedPath   string
}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

func (l *ConfigLoader) SetCustomPath(path string) {
	l.customPath = path
}

func (l *ConfigLoader) UsedPath() string {
	return l.usedPath
}

func (l *ConfigLoader) GetConfigPaths() []string {
	return l.getConfigPaths()
}

func (l *ConfigLoader) Load() (*Config, error) {
	configPaths := l.getConfigPaths()

	for _, path := range configPaths {
		if _, err := os.Stat(path); err != nil {
			continue
		}

		cfg := DefaultConfig()
		if err := LoadJSONCFile(path, cfg); err != nil {
			return nil, fmt.Errorf("parse config from %s: %w", path, err)
		}

		l.usedPath = path
		return cfg, nil
	}

	l.usedPath = ""
	return DefaultConfig(), nil
}

func (l *ConfigLoader) getConfigPaths() []string {
	paths := make([]string, 0)

	// 1. Custom config (from CLI flag)
	if l.customPath != "" {
		paths = append(paths, l.customPath)
	}

	// 2. Local config
	paths = append(paths, "./.txt2igconfig.jsonc")

	// 3. XDG global config
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		homeDir, _ := os.UserHomeDir()
		xdgConfigHome = filepath.Join(homeDir, ".config")
	}
	paths = append(paths, filepath.Join(xdgConfigHome, "txt2ig", "config.jsonc"))

	// 4. Home global config
	homeDir, _ := os.UserHomeDir()
	paths = append(paths, filepath.Join(homeDir, ".txt2ig", "config.jsonc"))

	return paths
}

func MergeConfigs(base, override *Config) *Config {
	result := &Config{}

	// FontFamily
	if len(override.FontFamily.Regular) > 0 || len(override.FontFamily.Bold) > 0 ||
		len(override.FontFamily.Italic) > 0 || len(override.FontFamily.BoldItalic) > 0 {
		result.FontFamily = override.FontFamily
	} else {
		result.FontFamily = base.FontFamily
	}

	// FontSize
	if override.FontSize != 0 {
		result.FontSize = override.FontSize
	} else {
		result.FontSize = base.FontSize
	}

	// FontColor
	if override.FontColor != "" {
		result.FontColor = override.FontColor
	} else {
		result.FontColor = base.FontColor
	}

	// BgColor
	if override.BgColor != "" {
		result.BgColor = override.BgColor
	} else {
		result.BgColor = base.BgColor
	}

	// TextJustify
	if override.TextJustify != "" {
		result.TextJustify = override.TextJustify
	} else {
		result.TextJustify = base.TextJustify
	}

	// TextBoxJustify
	if override.TextBoxJustify != "" {
		result.TextBoxJustify = override.TextBoxJustify
	} else {
		result.TextBoxJustify = base.TextBoxJustify
	}

	// TextBoxAlign
	if override.TextBoxAlign != "" {
		result.TextBoxAlign = override.TextBoxAlign
	} else {
		result.TextBoxAlign = base.TextBoxAlign
	}

	// TextBoxOffset
	if len(override.TextBoxOffset) > 0 {
		result.TextBoxOffset = override.TextBoxOffset
	} else {
		result.TextBoxOffset = base.TextBoxOffset
	}

	// TextBoxMaxWidth
	if override.TextBoxMaxWidth != 0 {
		result.TextBoxMaxWidth = override.TextBoxMaxWidth
	} else {
		result.TextBoxMaxWidth = base.TextBoxMaxWidth
	}

	// ScreenSize
	if len(override.ScreenSize) > 0 {
		result.ScreenSize = override.ScreenSize
	} else {
		result.ScreenSize = base.ScreenSize
	}

	// TextWrap (only set if explicitly different from zero value)
	result.TextWrap = base.TextWrap
	if override.TextWrap {
		result.TextWrap = true
	}

	// LineHeight
	if override.LineHeight != 0 {
		result.LineHeight = override.LineHeight
	} else {
		result.LineHeight = base.LineHeight
	}

	// PreProcessors
	if len(override.PreProcessors) > 0 {
		result.PreProcessors = override.PreProcessors
	} else {
		result.PreProcessors = base.PreProcessors
	}

	// PostProcessors
	if len(override.PostProcessors) > 0 {
		result.PostProcessors = override.PostProcessors
	} else {
		result.PostProcessors = base.PostProcessors
	}

	return result
}
