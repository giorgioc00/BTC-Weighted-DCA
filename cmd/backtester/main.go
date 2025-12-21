package main

import (
	"fmt"
	"log"

	"backtester/internal/infra/data"
	"backtester/internal/infra/util"
	"backtester/internal/usecase/engine"
	"backtester/internal/usecase/strategy"
)

func main() {
	// Load BTC price data from CSV
	prices, firstPriceDate, err := data.LoadPricesFromCSV("data/btc_price.csv")
	if err != nil {
		log.Fatal("Error loading CSV:", err)
	}

	highs, err := util.LoadCSVColumn("data/btc_price.csv", "High", true)
	if err != nil {
		log.Fatal("Error loading High prices:", err)
	}
	lows, err := util.LoadCSVColumn("data/btc_price.csv", "Low", true)
	if err != nil {
		log.Fatal("Error loading Low prices:", err)
	}
	volumes, err := util.LoadCSVVolumeColumn("data/btc_price.csv", "Vol.", true)
	if err != nil {
		log.Fatal("Error loading Volumes:", err)
	}

	// Load MVRV data from CSV
	mvrvData, err := util.LoadCSVColumn("data/btc_mvrv.csv", "CapMVRVCur", false)
	if err != nil {
		log.Printf("Warning: Could not load MVRV data: %v (using fallback)\n", err)
		mvrvData = nil
	}
	flowIn, err := util.LoadCSVColumnWithDefault("data/btc_mvrv.csv", "FlowInExUSD", false, 0)
	if err != nil {
		log.Printf("Warning: Could not load FlowInExUSD: %v (using zeros)\n", err)
		flowIn = nil
	}
	flowOut, err := util.LoadCSVColumnWithDefault("data/btc_mvrv.csv", "FlowOutExUSD", false, 0)
	if err != nil {
		log.Printf("Warning: Could not load FlowOutExUSD: %v (using zeros)\n", err)
		flowOut = nil
	}

	fmt.Println("Dynamic DCA Backtester with Multi-Indicator Score")
	fmt.Println("==================================================")
	fmt.Printf("Data points: %d\n", len(prices))
	fmt.Printf("Starting price: $%.2f\n", prices[0])
	fmt.Printf("Ending price: $%.2f\n", prices[len(prices)-1])
	if mvrvData != nil {
		fmt.Printf("MVRV data loaded: %d points\n", len(mvrvData))
	}

	// Configure the dynamic DCA strategy by reading the strategy_config.json
	var config strategy.DCAConfig
	err = util.LoadJSONConfig("config/strategy_config.json", &config)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	strat := strategy.NewDynamicDCA(config)
	strat.MVRVData = mvrvData
	strat.FirstPriceDate = firstPriceDate
	strat.Highs = highs
	strat.Lows = lows
	strat.Volumes = volumes
	strat.FlowInUSD = flowIn
	strat.FlowOutUSD = flowOut

	// Create and run the backtester
	runEngine := engine.NewEngine(strat)
	result := runEngine.Run(prices)

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
