package indicators

// MA200 calculates 200-day Moving Average score
// Pure function: price above MA200 = bullish (>50), below = bearish (<50)
func MA200(prices []float64) []float64 {
	period := 200
	result := newNeutralSlice(len(prices))

	if len(prices) < period {
		return result
	}

	// Calculate initial sum for first window
	sum := 0.0
	for j := 0; j < period; j++ {
		sum += prices[j]
	}

	// Sliding window calculation - O(n)
	for i := period - 1; i < len(prices); i++ {
		if i >= period {
			sum = sum - prices[i-period] + prices[i]
		}

		ma200 := sum / float64(period)

		if ma200 == 0 {
			result[i] = NeutralScore
			continue
		}

		// Percentage distance from MA200
		percentDiff := ((prices[i] - ma200) / ma200) * 100.0
		// Map ±20% to 0-100 scale
		score := 50.0 + (percentDiff * 2.5)

		result[i] = clamp(score, 0, 100)
	}

	return result
}
