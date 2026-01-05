package metrics

// Trade represents a single trade (re-exported for metrics)
type Trade struct {
	Index        int
	Price        float64
	Amount       float64
	Units        float64
	TotalUnits   float64
	TotalInvest  float64
	Score        float64
	RSI          float64
	FearAndGreed float64
	MVRV         float64
	MA200        float64
	Bollinger    float64
	ATR          float64
	NetFlow      float64
	VolumeFlow   float64
}

// Result holds backtest results
type Result struct {
	Trades           []Trade
	TotalInvested    float64
	TotalUnits       float64
	FinalPrice       float64
	FinalValue       float64
	TotalReturn      float64
	AverageCostBasis float64
	NumTrades        int
}
