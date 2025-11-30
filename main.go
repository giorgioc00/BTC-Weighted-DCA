package main

import (
	"backtester/data"
	"backtester/engine"
	"backtester/strategy"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// Load BTC price data from CSV
	prices, err := data.LoadPricesFromCSV("data/btc.csv")
	if err != nil {
		log.Fatal("Error loading CSV:", err)
	}

	fmt.Println("Dynamic DCA Backtester with RSI")
	fmt.Println("================================")
	fmt.Printf("Data points: %d\n", len(prices))
	fmt.Printf("Starting price: $%.2f\n", prices[0])
	fmt.Printf("Ending price: $%.2f\n", prices[len(prices)-1])

	// Configure the dynamic DCA strategy by reading the rsi_config.json
	content, err := ioutil.ReadFile("./rsi_config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var config strategy.DCAConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	strat := strategy.NewDynamicDCA(config)

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
	compareWithStandardDCA(prices, config.BaseAmount, 7)
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
}
