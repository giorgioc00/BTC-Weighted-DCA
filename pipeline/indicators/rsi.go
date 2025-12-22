package indicators

// RSI calculates the Relative Strength Index for a given price series
// Returns a 0-100 score. Pure function: []float64 → []float64
func RSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return nil
	}

	rsi := make([]float64, len(prices))
	fillSlice(rsi, 0, period, NeutralScore)

	// Calculate price changes
	gains := make([]float64, len(prices))
	losses := make([]float64, len(prices))

	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = -change
		}
	}

	// Calculate initial average gain and loss
	var avgPeriodGain, avgPeriodLoss float64
	for i := 1; i <= period; i++ {
		avgPeriodGain += gains[i]
		avgPeriodLoss += losses[i]
	}
	avgPeriodGain /= float64(period)
	avgPeriodLoss /= float64(period)

	// Calculate RSI using smoothed moving average
	for i := period; i < len(prices); i++ {
		if i > period {
			avgPeriodGain = (avgPeriodGain*float64(period-1) + gains[i]) / float64(period)
			avgPeriodLoss = (avgPeriodLoss*float64(period-1) + losses[i]) / float64(period)
		}

		if avgPeriodLoss == 0 {
			rsi[i] = 100
		} else {
			relativeStrength := avgPeriodGain / avgPeriodLoss
			rsi[i] = 100 - (100 / (1 + relativeStrength))
		}
	}
	return rsi
}
