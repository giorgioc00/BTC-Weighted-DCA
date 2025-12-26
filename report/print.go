package report

import (
	"fmt"

	"backtester/pipeline/metrics"
)

// PrintSummary displays result summary.
func PrintSummary(r *metrics.Result) {
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

// PrintTrades displays trade history.
func PrintTrades(r *metrics.Result) {
	fmt.Println("\n========== TRADE HISTORY ==========")
	fmt.Printf("%-6s %-12s %-12s %-10s %-8s %-8s %-8s\n",
		"Index", "Price", "Amount", "Units", "Score", "RSI", "MVRV")
	fmt.Println("-----------------------------------------------------------------------")
	for _, t := range r.Trades {
		fmt.Printf("%-6d $%-11.2f $%-11.2f %-10.4f %-8.2f %-8.2f %-8.2f\n",
			t.Index, t.Price, t.Amount, t.Units, t.Score, t.RSI, t.MVRV)
	}
	fmt.Println("=======================================================================")
}

// PrintTradesDetailed displays detailed trade information.
func PrintTradesDetailed(r *metrics.Result) {
	fmt.Println("\n========== DETAILED TRADE HISTORY ==========")
	for i, t := range r.Trades {
		fmt.Printf("\nTrade #%d (Index %d)\n", i+1, t.Index)
		fmt.Printf("  Price:        $%.2f\n", t.Price)
		fmt.Printf("  Amount:       $%.2f\n", t.Amount)
		fmt.Printf("  Units:        %.6f\n", t.Units)
		fmt.Printf("  Score:        %.2f\n", t.Score)
		fmt.Println("  Indicators:")
		fmt.Printf("    RSI:          %.2f\n", t.RSI)
		fmt.Printf("    Fear&Greed:   %.2f\n", t.FearAndGreed)
		fmt.Printf("    MVRV:         %.2f\n", t.MVRV)
		fmt.Printf("    MA200:        %.2f\n", t.MA200)
		fmt.Printf("    Bollinger:    %.2f\n", t.Bollinger)
		fmt.Printf("    ATR:          %.2f\n", t.ATR)
		fmt.Printf("    NetFlow:      %.2f\n", t.NetFlow)
		fmt.Printf("    VolumeFlow:   %.2f\n", t.VolumeFlow)
		fmt.Printf("  Total Units:  %.6f\n", t.TotalUnits)
		fmt.Printf("  Total Invest: $%.2f\n", t.TotalInvest)
	}
	fmt.Println("\n============================================")
}
