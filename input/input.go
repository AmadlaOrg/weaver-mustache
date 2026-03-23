package input

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads input data from a file path or stdin ("-").
// Auto-detects JSON vs YAML.
func Load(filePath string) (any, error) {
	var reader io.Reader
	if filePath == "-" {
		reader = os.Stdin
	} else {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return map[string]any{}, nil
	}

	return Parse(raw)
}

// Parse auto-detects JSON or YAML and returns parsed data.
func Parse(raw []byte) (any, error) {
	trimmed := raw
	for len(trimmed) > 0 && (trimmed[0] == ' ' || trimmed[0] == '\t' || trimmed[0] == '\n' || trimmed[0] == '\r') {
		trimmed = trimmed[1:]
	}

	// Try JSON first if starts with { or [
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		var result any
		if err := json.Unmarshal(raw, &result); err == nil {
			return result, nil
		}
	}

	// Try YAML
	var result any
	if err := yaml.Unmarshal(raw, &result); err == nil && result != nil {
		return result, nil
	}

	// Fallback JSON
	var fallback any
	if err := json.Unmarshal(raw, &fallback); err == nil {
		return fallback, nil
	}

	return nil, errors.New("input is neither valid JSON nor YAML")
}
