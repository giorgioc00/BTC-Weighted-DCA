package util

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// LoadCSVColumnWithDefault reads a specific column by name and substitutes defaultVal
// for missing or invalid entries. Values can be reversed to chronological order.
func LoadCSVColumnWithDefault(filepath string, columnName string, reverse bool, defaultVal float64) ([]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return nil, err
	}
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, nil
	}

	header := records[0]
	columnIndex := -1
	for i, col := range header {
		if col == columnName {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		return nil, nil
	}

	var values []float64
	for i := 1; i < len(records); i++ {
		if columnIndex >= len(records[i]) {
			values = append(values, defaultVal)
			continue
		}

		valueStr := records[i][columnIndex]
		if valueStr == "" {
			values = append(values, defaultVal)
			continue
		}

		valueStr = strings.ReplaceAll(valueStr, ",", "")
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			values = append(values, defaultVal)
			continue
		}

		values = append(values, value)
	}

	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}

	return values, nil
}

// LoadCSVVolumeColumn reads a volume column with suffixes (K/M/B) and returns values.
// Set reverse to true to order chronologically (oldest first).
func LoadCSVVolumeColumn(filepath string, columnName string, reverse bool) ([]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return nil, err
	}
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, nil
	}

	header := records[0]
	columnIndex := -1
	for i, col := range header {
		if col == columnName {
			columnIndex = i
			break
		}
	}
	if columnIndex == -1 {
		return nil, nil
	}

	var values []float64
	for i := 1; i < len(records); i++ {
		if columnIndex >= len(records[i]) {
			values = append(values, 0)
			continue
		}

		valueStr := strings.TrimSpace(records[i][columnIndex])
		if valueStr == "" || valueStr == "-" {
			values = append(values, 0)
			continue
		}

		multiplier := 1.0
		lastChar := valueStr[len(valueStr)-1]
		switch lastChar {
		case 'K':
			multiplier = 1_000
			valueStr = valueStr[:len(valueStr)-1]
		case 'M':
			multiplier = 1_000_000
			valueStr = valueStr[:len(valueStr)-1]
		case 'B':
			multiplier = 1_000_000_000
			valueStr = valueStr[:len(valueStr)-1]
		}

		valueStr = strings.ReplaceAll(valueStr, ",", "")
		v, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			values = append(values, 0)
			continue
		}

		values = append(values, v*multiplier)
	}

	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}

	return values, nil
}

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

// LoadLastCSVColumnRow reads the first data row (after header) from a specific column
// This represents the "last" date chronologically before data reversal (oldest date)
func LoadLastCSVColumnRow(filepath string, columnName string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Skip BOM if present
	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return "", err
	}
	if r != '\uFEFF' {
		br.UnreadRune() // Not a BOM, put it back
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	if len(records) < 2 {
		return "", nil // Need at least header + one data row
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
		return "", nil // Column not found
	}

	// Return the first data row value (row 1, after header at row 0)
	if columnIndex >= len(records[1]) {
		return "", nil // Column doesn't exist in first data row
	}

	return records[1][columnIndex], nil
}

// LoadFirstAndLastCSVDate reads both the first and last data row values from a specific column
// Returns (firstRow, lastRow, error) - useful for date validation regardless of sort order
func LoadFirstAndLastCSVDate(filepath string, columnName string) (string, string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// Skip BOM if present
	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return "", "", err
	}
	if r != '\uFEFF' {
		br.UnreadRune() // Not a BOM, put it back
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return "", "", err
	}

	if len(records) < 2 {
		return "", "", nil // Need at least header + one data row
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
		return "", "", nil // Column not found
	}

	// Return first and last data row values
	if columnIndex >= len(records[1]) || columnIndex >= len(records[len(records)-1]) {
		return "", "", nil // Column doesn't exist in data rows
	}

	firstRowDateRaw := records[1][columnIndex]
	lastRowDateRaw := records[len(records)-1][columnIndex]

	return firstRowDateRaw, lastRowDateRaw, nil
}
