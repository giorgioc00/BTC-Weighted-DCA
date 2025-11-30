package util

import (
	"encoding/json"
	"os"
)

// LoadJSONConfig reads a JSON file and unmarshals it into the provided target
func LoadJSONConfig(filepath string, target interface{}) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, target) // it wants a pointer
	if err != nil {
		return err
	}

	return nil
}
