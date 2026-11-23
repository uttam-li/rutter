package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tailscale/hujson"
)

// ReadDocument parses a JSONC file into generic Go values. It is used by
// `rutter config show --json` so agents get valid JSON even though the file
// on disk contains comments.
func ReadDocument(path string) (any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	standard, err := hujson.Standardize(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: parse %s: %w", ErrInvalid, path, err)
	}

	var value any
	decoder := json.NewDecoder(bytes.NewReader(standard))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("%w: decode %s: %w", ErrInvalid, path, err)
	}
	return value, nil
}
