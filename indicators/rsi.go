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

	// STEP 1: Calculate price changes
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

	// STEP 2: Calculate initial average gain and loss (Simple Moving Average)
	var avgPeriodGain, avgPeriodLoss float64 // for the last period
	for i := 1; i <= period; i++ {
		avgPeriodGain += gains[i]
		avgPeriodLoss += losses[i]
	}
	avgPeriodGain /= float64(period)
	avgPeriodLoss /= float64(period)

	// STEP 3 & 4: Calculate RSI using smoothed moving average
	for i := period; i < len(prices); i++ {
		// STEP 4: Apply smoothing for subsequent values (i > period)
		if i > period {
			avgPeriodGain = (avgPeriodGain*float64(period-1) + gains[i]) / float64(period)
			avgPeriodLoss = (avgPeriodLoss*float64(period-1) + losses[i]) / float64(period)
		}
		// STEP 3: Calculate RSI (first at i=period uses initial averages)

		if avgPeriodLoss == 0 {
			rsi[i] = 100
		} else {
			relativeStrength := avgPeriodGain / avgPeriodLoss
			rsi[i] = 100 - (100 / (1 + relativeStrength))
		}
	}
	return rsi
}
