package data

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// LoadPricesFromCSV reads price data from a CSV file
// Returns prices in chronological order (oldest first)
func LoadPricesFromCSV(filepath string) ([]float64, error) {
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

	// Skip header row, parse prices
	var prices []float64
	for i := 1; i < len(records); i++ {
		// Price is column 1 (index 1), format: "87,342.0"
		priceStr := records[i][1]
		priceStr = strings.ReplaceAll(priceStr, ",", "") // Remove commas

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			continue // Skip invalid rows
		}
		prices = append(prices, price)
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(prices)-1; i < j; i, j = i+1, j-1 {
		prices[i], prices[j] = prices[j], prices[i]
	}

	return prices, nil
}
