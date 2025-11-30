package util

import (
	"encoding/json"
	"os"
)

// LoadJSONConfig reads a JSON file and unmarshals it into the provided target
// The target parameter should be a pointer to the struct you want to populate
func LoadJSONConfig(filepath string, target interface{}) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, target)
	if err != nil {
		return err
	}

	return nil
}
