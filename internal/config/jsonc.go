package config

import (
	"fmt"
	"os"

	"github.com/titanous/json5"
)

func ParseJSONC(data []byte, v interface{}) error {
	// JSON5 handles:
	// - Comments (// and /* */)
	// - Trailing commas
	// - Unquoted keys (optional)
	// - Single quotes (optional)
	return json5.Unmarshal(data, v)
}

func LoadJSONCFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	return ParseJSONC(data, v)
}
