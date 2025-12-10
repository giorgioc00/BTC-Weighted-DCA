package data

import (
	"backtester/util"
)

// LoadPricesFromCSV reads price data from a CSV file.
// Returns prices in chronological order (oldest first).
func LoadPricesFromCSV(filepath string) ([]float64, error) {
	return util.LoadCSVColumn(filepath, "Price", true)
}
