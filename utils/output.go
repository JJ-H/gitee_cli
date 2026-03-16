package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// PrintJSON outputs the given data as compact JSON to stdout
func PrintJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}
