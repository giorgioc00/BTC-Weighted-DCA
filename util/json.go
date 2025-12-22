package util

import (
	"encoding/json"
	"os"
)

// Validator interface for types that can validate themselves
type Validator interface {
	Validate() error
}

// LoadJSONConfig reads JSON file and unmarshals into target
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

// LoadAndValidateJSONConfig reads, unmarshals, and validates JSON config
func LoadAndValidateJSONConfig(filepath string, target Validator) error {
	err := LoadJSONConfig(filepath, target)
	if err != nil {
		return err
	}

	return target.Validate()
}
