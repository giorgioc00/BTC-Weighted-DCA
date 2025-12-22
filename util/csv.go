package util

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// LoadCSVColumnWithDefault reads a column with default values for missing/invalid entries
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

// LoadCSVVolumeColumn reads volume column with K/M/B suffixes
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

// LoadCSVColumn reads a column from CSV by name
func LoadCSVColumn(filepath string, columnName string, reverse bool) ([]float64, error) {
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
			continue
		}

		valueStr := records[i][columnIndex]
		if valueStr == "" {
			values = append(values, 50.0)
			continue
		}

		valueStr = strings.ReplaceAll(valueStr, ",", "")
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			values = append(values, 50.0)
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

// LoadLastCSVColumnRow reads the first data row from a column
func LoadLastCSVColumnRow(filepath string, columnName string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return "", err
	}
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	if len(records) < 2 {
		return "", nil
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
		return "", nil
	}

	if columnIndex >= len(records[1]) {
		return "", nil
	}

	return records[1][columnIndex], nil
}

// LoadFirstAndLastCSVDate reads first and last data row values
func LoadFirstAndLastCSVDate(filepath string, columnName string) (string, string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	r, _, err := br.ReadRune()
	if err != nil {
		return "", "", err
	}
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	reader := csv.NewReader(br)
	records, err := reader.ReadAll()
	if err != nil {
		return "", "", err
	}

	if len(records) < 2 {
		return "", "", nil
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
		return "", "", nil
	}

	if columnIndex >= len(records[1]) || columnIndex >= len(records[len(records)-1]) {
		return "", "", nil
	}

	firstRowDateRaw := records[1][columnIndex]
	lastRowDateRaw := records[len(records)-1][columnIndex]

	return firstRowDateRaw, lastRowDateRaw, nil
}
