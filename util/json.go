package util

import (
	"encoding/json"
	"os"
)

// Validator is an interface for types that can validate themselves
type Validator interface {
	Validate() error
}

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

// LoadAndValidateJSONConfig reads a JSON file, unmarshals it, and validates it
// The target must implement the Validator interface
func LoadAndValidateJSONConfig(filepath string, target Validator) error {
	err := LoadJSONConfig(filepath, target)
	if err != nil {
		return err
	}

	return target.Validate()
}
