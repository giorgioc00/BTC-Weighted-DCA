package data

import (
	"time"

	"backtester/util"
)

// LoadPricesFromCSV reads price data from CSV
// Returns prices in chronological order (oldest first) and first date
func LoadPricesFromCSV(filepath string) ([]float64, time.Time, error) {
	firstRowDateRaw, lastRowDateRaw, err := util.LoadFirstAndLastCSVDate(filepath, "Date")
	if err != nil {
		return nil, time.Time{}, err
	}

	var firstDate time.Time
	if firstRowDateRaw != "" && lastRowDateRaw != "" {
		firstParsed, err1 := util.ParseDateFlexible(firstRowDateRaw)
		lastParsed, err2 := util.ParseDateFlexible(lastRowDateRaw)
		if err1 == nil && err2 == nil {
			if lastParsed.Before(firstParsed) {
				firstDate = lastParsed
			} else {
				firstDate = firstParsed
			}
		} else if err1 == nil {
			firstDate = firstParsed
		} else if err2 == nil {
			firstDate = lastParsed
		}
	}

	prices, err := util.LoadCSVColumn(filepath, "Price", true)
	if err != nil {
		return nil, time.Time{}, err
	}

	return prices, firstDate, nil
}
