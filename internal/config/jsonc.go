package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJSONC(data []byte, v interface{}) error {
	cleaned := stripComments(data)
	return json.Unmarshal(cleaned, v)
}

func stripComments(data []byte) []byte {
	var result []byte
	inString := false
	inSingleLineComment := false
	inMultiLineComment := false

	for i := 0; i < len(data); i++ {
		// If we're in a string, check for end quote
		if inString {
			if data[i] == '"' && (i == 0 || data[i-1] != '\\') {
				inString = false
			}
			result = append(result, data[i])
			continue
		}

		// If we're in a single-line comment, check for newline
		if inSingleLineComment {
			if data[i] == '\n' {
				inSingleLineComment = false
				result = append(result, '\n')
			}
			continue
		}

		// If we're in a multi-line comment, check for */
		if inMultiLineComment {
			if i < len(data)-1 && data[i] == '*' && data[i+1] == '/' {
				inMultiLineComment = false
				i++
			}
			continue
		}

		// Check for start of string
		if data[i] == '"' {
			inString = true
			result = append(result, data[i])
			continue
		}

		// Check for start of comment
		if i < len(data)-1 && data[i] == '/' {
			// Single-line comment
			if data[i+1] == '/' {
				inSingleLineComment = true
				i++
				continue
			}
			// Multi-line comment
			if data[i+1] == '*' {
				inMultiLineComment = true
				i++
				continue
			}
		}

		result = append(result, data[i])
	}

	return result
}

func LoadJSONCFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	return ParseJSONC(data, v)
}
