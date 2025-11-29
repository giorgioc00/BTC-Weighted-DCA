package backtester

import (
	"backtester/strategy"
	"fmt"
)

// Trade represents a single trade execution
type Trade struct {
	Index       int
	Price       float64
	Amount      float64 // Dollar amount invested
	Units       float64 // Units/shares purchased
	RSI         float64
	TotalUnits  float64 // Running total of units held
	TotalInvest float64 // Running total invested
}

// Result holds the backtesting results
type Result struct {
	Trades           []Trade
	TotalInvested    float64
	TotalUnits       float64
	FinalPrice       float64
	FinalValue       float64
	TotalReturn      float64 // Percentage return
	AverageCostBasis float64
	NumTrades        int
}

// Engine is the backtesting engine
type Engine struct {
	Strategy    *strategy.DynamicDCA
	InitialCash float64
}

// NewEngine creates a new backtesting engine
func NewEngine(strat *strategy.DynamicDCA) *Engine {
	return &Engine{
		Strategy: strat,
	}
}

// Run executes the backtest on the given price series
func (e *Engine) Run(prices []float64) *Result {
	signals := e.Strategy.GenerateSignals(prices)
	if signals == nil {
		return nil
	}

	result := &Result{
		Trades: make([]Trade, 0),
	}

	var totalUnits, totalInvested float64

	// Execute trades at specified intervals
	for i, signal := range signals {
		// Skip if not at interval (e.g., buy weekly, not daily)
		if i%e.Strategy.Config.ExecutionInterval != 0 {
			continue
		}

		// Skip initial period where RSI isn't reliable
		if i < e.Strategy.Config.RSIPeriod {
			continue
		}

		units := signal.Amount / signal.Price
		totalUnits += units
		totalInvested += signal.Amount

		trade := Trade{
			Index:       i,
			Price:       signal.Price,
			Amount:      signal.Amount,
			Units:       units,
			RSI:         signal.RSI,
			TotalUnits:  totalUnits,
			TotalInvest: totalInvested,
		}

		result.Trades = append(result.Trades, trade)
	}

	// Calculate final results
	result.TotalInvested = totalInvested
	result.TotalUnits = totalUnits
	result.FinalPrice = prices[len(prices)-1]
	result.FinalValue = totalUnits * result.FinalPrice
	result.NumTrades = len(result.Trades)

	if totalInvested > 0 {
		result.TotalReturn = ((result.FinalValue - totalInvested) / totalInvested) * 100
		result.AverageCostBasis = totalInvested / totalUnits
	}

	return result
}

// PrintSummary prints a summary of the backtest results
func (r *Result) PrintSummary() {
	fmt.Println("\n========== BACKTEST RESULTS ==========")
	fmt.Printf("Number of Trades:    %d\n", r.NumTrades)
	fmt.Printf("Total Invested:      $%.2f\n", r.TotalInvested)
	fmt.Printf("Total Units:         %.6f\n", r.TotalUnits)
	fmt.Printf("Average Cost Basis:  $%.2f\n", r.AverageCostBasis)
	fmt.Printf("Final Price:         $%.2f\n", r.FinalPrice)
	fmt.Printf("Final Value:         $%.2f\n", r.FinalValue)
	fmt.Printf("Total Return:        %.2f%%\n", r.TotalReturn)
	fmt.Printf("Profit/Loss:         $%.2f\n", r.FinalValue-r.TotalInvested)
	fmt.Println("=======================================")
}

// PrintTrades prints detailed trade information
func (r *Result) PrintTrades() {
	fmt.Println("\n========== TRADE HISTORY ==========")
	fmt.Printf("%-6s %-12s %-12s %-10s %-8s\n", "Index", "Price", "Amount", "Units", "RSI")
	fmt.Println("----------------------------------------------------")
	for _, t := range r.Trades {
		fmt.Printf("%-6d $%-11.2f $%-11.2f %-10.4f %-8.2f\n",
			t.Index, t.Price, t.Amount, t.Units, t.RSI)
	}
	fmt.Println("====================================")
}
