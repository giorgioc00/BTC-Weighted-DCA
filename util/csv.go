package util

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// LoadCSVColumn reads a specific column from a CSV file by column name
// Returns values in chronological order (oldest first) by reversing the data
func LoadCSVColumn(filepath string, columnName string, reverse bool) ([]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Skip BOM if present
	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return nil, err
	}
	if r != '\uFEFF' {
		br.UnreadRune() // Not a BOM, put it back
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, nil
	}

	// Find the column index by name
	header := records[0]
	columnIndex := -1
	for i, col := range header {
		if col == columnName {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		return nil, nil // Column not found
	}

	// Parse values from the column
	var values []float64
	for i := 1; i < len(records); i++ {
		if columnIndex >= len(records[i]) {
			continue // Skip rows with insufficient columns
		}

		valueStr := records[i][columnIndex]
		if valueStr == "" {
			values = append(values, 50.0) // Use neutral value for missing data
			continue
		}

		valueStr = strings.ReplaceAll(valueStr, ",", "") // Remove commas
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			values = append(values, 50.0) // Use neutral value for invalid data
			continue
		}

		values = append(values, value)
	}

	// Reverse to get chronological order (oldest first) if requested
	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}

	return values, nil
}
