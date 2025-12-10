package main

import (
	"backtester/data"
	backtester "backtester/engine"
	"backtester/strategy"
	"backtester/util"
	"fmt"
	"log"
)

func main() {
	// Load BTC price data from CSV
	prices, err := data.LoadPricesFromCSV("data/btc_price.csv")
	if err != nil {
		log.Fatal("Error loading CSV:", err)
	}

	// Load MVRV data from CSV
	mvrvData, err := util.LoadCSVColumn("data/btc_mvrv.csv", "CapMVRVCur", false)
	if err != nil {
		log.Printf("Warning: Could not load MVRV data: %v (using fallback)\n", err)
		mvrvData = nil
	}

	fmt.Println("Dynamic DCA Backtester with Multi-Indicator Score")
	fmt.Println("==================================================")
	fmt.Printf("Data points: %d\n", len(prices))
	fmt.Printf("Starting price: $%.2f\n", prices[0])
	fmt.Printf("Ending price: $%.2f\n", prices[len(prices)-1])
	if mvrvData != nil {
		fmt.Printf("MVRV data loaded: %d points\n", len(mvrvData))
	}

	// Configure the dynamic DCA strategy by reading the rsi_config.json
	config := strategy.DefaultConfig()
	err = util.LoadJSONConfig("config/rsi_config.json", &config)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	if err := config.Weights.Validate(); err != nil {
		log.Fatal("Invalid weights configuration: ", err)
	}

	strat := strategy.NewDynamicDCA(config)
	strat.MVRVData = mvrvData

	// Create and run the backtester
	engine := backtester.NewEngine(strat)
	result := engine.Run(prices)

	if result == nil {
		fmt.Println("Error: Could not run backtest")
		return
	}

	// Print results
	result.PrintSummary()

	// Optional: Print trade history (uncomment to see all trades)
	// result.PrintTrades()

	// Compare with standard DCA (fixed amount)
	fmt.Println("\n--- Comparison with Standard DCA ---")
	compareWithStandardDCA(prices, config.RSI.BaseAmount, 7)
}

// compareWithStandardDCA runs a simple fixed-amount DCA for comparison
func compareWithStandardDCA(prices []float64, amount float64, interval int) {
	var totalInvested, totalUnits float64

	for i := 14; i < len(prices); i += interval {
		units := amount / prices[i]
		totalUnits += units
		totalInvested += amount
	}

	finalValue := totalUnits * prices[len(prices)-1]
	avgCost := totalInvested / totalUnits
	returnPct := ((finalValue - totalInvested) / totalInvested) * 100

	fmt.Printf("Standard DCA (fixed $%.0f):\n", amount)
	fmt.Printf("  Total Invested: $%.2f\n", totalInvested)
	fmt.Printf("  Final Value:    $%.2f\n", finalValue)
	fmt.Printf("  Avg Cost Basis: $%.2f\n", avgCost)
	fmt.Printf("  Total Return:   %.2f%%\n", returnPct)
	fmt.Printf("  Profit/Loss:    $%.2f\n", finalValue-totalInvested)
}
