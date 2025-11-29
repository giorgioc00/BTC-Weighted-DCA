package indicators

// RSI calculates the Relative Strength Index for a given price series
// period is typically 14
func RSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return nil
	}

	rsi := make([]float64, len(prices))

	// Initialize with NaN-like values for insufficient data points
	for i := 0; i < period; i++ {
		rsi[i] = 50 // neutral RSI for insufficient data
	}

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
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI using smoothed moving average
	for i := period; i < len(prices); i++ {
		if i > period {
			avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
		}

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi
}
