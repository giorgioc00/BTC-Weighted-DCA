package data

import (
	"backtester/util"
)

// LoadPricesFromCSV reads price data from a CSV file
// Returns prices in chronological order (oldest first)
func LoadPricesFromCSV(filepath string) ([]float64, error) {

	records, err := util.LoadCSVColumn("data/btc_price.csv", "Price", true)
	if err != nil {
		return nil, err
	}

	// Skip header row, parse prices
	var prices []float64
	for i := 1; i < len(records); i++ {
		// Price is column 1 (index 1), format: "87,342.0"
		price := records[i]
		if err != nil {
			continue // Skip invalid rows
		}
		prices = append(prices, price)
	}

	return prices, nil
}
