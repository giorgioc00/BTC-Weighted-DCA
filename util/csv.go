package util

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// CSV default values for missing/invalid data
const (
	DefaultIndicatorValue = 50.0
	DefaultVolumeValue    = 0.0
)

// csvReader holds CSV records and header info
type csvReader struct {
	records     [][]string
	columnIndex int
}

// openCSVFile handles BOM, file opening, and CSV parsing
func openCSVFile(filepath string, columnName string) (*csvReader, error) {
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

	columnIndex := findColumnIndex(records[0], columnName)
	if columnIndex == -1 {
		return nil, nil
	}

	return &csvReader{
		records:     records,
		columnIndex: columnIndex,
	}, nil
}

// findColumnIndex finds the index of a column in the header
func findColumnIndex(header []string, columnName string) int {
	for i, col := range header {
		if col == columnName {
			return i
		}
	}
	return -1
}

// reverseSlice reverses a float64 slice in place
func reverseSlice(values []float64) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

// parseFloatValue removes commas and parses a float64
func parseFloatValue(valueStr string) (float64, error) {
	valueStr = strings.ReplaceAll(valueStr, ",", "")
	return strconv.ParseFloat(valueStr, 64)
}

// LoadCSVColumnWithDefault reads a column with default values for missing/invalid entries
func LoadCSVColumnWithDefault(filepath string, columnName string, reverse bool, defaultVal float64) ([]float64, error) {
	csv, err := openCSVFile(filepath, columnName)
	if err != nil || csv == nil {
		return nil, err
	}

	var values []float64
	for i := 1; i < len(csv.records); i++ {
		if csv.columnIndex >= len(csv.records[i]) {
			values = append(values, defaultVal)
			continue
		}

		valueStr := csv.records[i][csv.columnIndex]
		if valueStr == "" {
			values = append(values, defaultVal)
			continue
		}

		value, err := parseFloatValue(valueStr)
		if err != nil {
			values = append(values, defaultVal)
			continue
		}

		values = append(values, value)
	}

	if reverse {
		reverseSlice(values)
	}

	return values, nil
}

// LoadCSVVolumeColumn reads volume column with K/M/B suffixes
func LoadCSVVolumeColumn(filepath string, columnName string, reverse bool) ([]float64, error) {
	csv, err := openCSVFile(filepath, columnName)
	if err != nil || csv == nil {
		return nil, err
	}

	var values []float64
	for i := 1; i < len(csv.records); i++ {
		if csv.columnIndex >= len(csv.records[i]) {
			values = append(values, DefaultVolumeValue)
			continue
		}

		valueStr := strings.TrimSpace(csv.records[i][csv.columnIndex])
		if valueStr == "" || valueStr == "-" {
			values = append(values, DefaultVolumeValue)
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

		v, err := parseFloatValue(valueStr)
		if err != nil {
			values = append(values, DefaultVolumeValue)
			continue
		}

		values = append(values, v*multiplier)
	}

	if reverse {
		reverseSlice(values)
	}

	return values, nil
}

// LoadCSVColumn reads a column from CSV by name, using DefaultIndicatorValue for missing/invalid entries
func LoadCSVColumn(filepath string, columnName string, reverse bool) ([]float64, error) {
	return LoadCSVColumnWithDefault(filepath, columnName, reverse, DefaultIndicatorValue)
}

// LoadLastCSVColumnRow reads the first data row from a column
func LoadLastCSVColumnRow(filepath string, columnName string) (string, error) {
	csv, err := openCSVFile(filepath, columnName)
	if err != nil || csv == nil {
		return "", err
	}

	if csv.columnIndex >= len(csv.records[1]) {
		return "", nil
	}

	return csv.records[1][csv.columnIndex], nil
}

// LoadFirstAndLastCSVDate reads first and last data row values
func LoadFirstAndLastCSVDate(filepath string, columnName string) (string, string, error) {
	csv, err := openCSVFile(filepath, columnName)
	if err != nil || csv == nil {
		return "", "", err
	}

	lastIndex := len(csv.records) - 1
	if csv.columnIndex >= len(csv.records[1]) || csv.columnIndex >= len(csv.records[lastIndex]) {
		return "", "", nil
	}

	firstRowDateRaw := csv.records[1][csv.columnIndex]
	lastRowDateRaw := csv.records[lastIndex][csv.columnIndex]

	return firstRowDateRaw, lastRowDateRaw, nil
}
