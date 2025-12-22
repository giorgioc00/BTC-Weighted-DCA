package execution

import (
	"backtester/pipeline/metrics"
	"backtester/pipeline/signals"
)

// Trade represents a single executed trade
type Trade struct {
	Index       int
	Price       float64
	Amount      float64
	Units       float64
	TotalUnits  float64
	TotalInvest float64
	// Indicator values at trade time
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

// Execute runs the backtest by executing signals at intervals
// This is a pure function: signals + config → trades
func Execute(sigs []signals.Signal, config signals.Config) []Trade {
	if sigs == nil {
		return nil
	}

	trades := make([]Trade, 0)
	var totalUnits, totalInvested float64

	for i, signal := range sigs {
		// Execute only at specified intervals
		if i%config.ExecutionInterval != 0 {
			continue
		}

		// Skip warmup period
		if i < config.RSI.RSIPeriod {
			continue
		}

		units := signal.Amount / signal.Price
		totalUnits += units
		totalInvested += signal.Amount

		trade := Trade{
			Index:        i,
			Price:        signal.Price,
			Amount:       signal.Amount,
			Units:        units,
			TotalUnits:   totalUnits,
			TotalInvest:  totalInvested,
			Score:        signal.Score,
			RSI:          signal.RSI,
			FearAndGreed: signal.FearAndGreed,
			MVRV:         signal.MVRV,
			MA200:        signal.MA200,
			Bollinger:    signal.Bollinger,
			ATR:          signal.ATR,
			NetFlow:      signal.NetFlow,
			VolumeFlow:   signal.VolumeFlow,
		}

		trades = append(trades, trade)
	}

	return trades
}

// CalculateResult computes final metrics from trades
// This is a pure function: trades + final price → result
func CalculateResult(trades []Trade, finalPrice float64) *metrics.Result {
	if len(trades) == 0 {
		return &metrics.Result{}
	}

	lastTrade := trades[len(trades)-1]
	totalInvested := lastTrade.TotalInvest
	totalUnits := lastTrade.TotalUnits
	finalValue := totalUnits * finalPrice

	// Convert execution trades to metrics trades
	metricsTrades := make([]metrics.Trade, len(trades))
	for i, t := range trades {
		metricsTrades[i] = metrics.Trade{
			Index:        t.Index,
			Price:        t.Price,
			Amount:       t.Amount,
			Units:        t.Units,
			TotalUnits:   t.TotalUnits,
			TotalInvest:  t.TotalInvest,
			Score:        t.Score,
			RSI:          t.RSI,
			FearAndGreed: t.FearAndGreed,
			MVRV:         t.MVRV,
			MA200:        t.MA200,
			Bollinger:    t.Bollinger,
			ATR:          t.ATR,
			NetFlow:      t.NetFlow,
			VolumeFlow:   t.VolumeFlow,
		}
	}

	result := &metrics.Result{
		Trades:        metricsTrades,
		TotalInvested: totalInvested,
		TotalUnits:    totalUnits,
		FinalPrice:    finalPrice,
		FinalValue:    finalValue,
		NumTrades:     len(trades),
	}

	if totalInvested > 0 {
		result.TotalReturn = ((finalValue - totalInvested) / totalInvested) * 100
		result.AverageCostBasis = totalInvested / totalUnits
	}

	return result
}
