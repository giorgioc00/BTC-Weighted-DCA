package indicators

// MA200 calculates the 200-day Simple Moving Average
// Returns a score from 0-100 based on price position relative to MA200:
// - Price > MA200: bullish signal (higher score)
// - Price < MA200: bearish signal (lower score)
func MA200(prices []float64) []float64 {
	period := 200
	result := newNeutralSlice(len(prices))

	// Return neutral values for insufficient data
	if len(prices) < period {
		return result
	}

	// Calculate initial sum for first window
	sum := 0.0
	for j := 0; j < period; j++ {
		sum += prices[j]
	}

	// Sliding window calculation - O(n) instead of O(n*period)
	for i := period - 1; i < len(prices); i++ {
		// Update sum with sliding window
		if i >= period {
			sum = sum - prices[i-period] + prices[i]
		}

		ma200 := sum / float64(period)

		// Handle division by zero when MA200 is 0
		if ma200 == 0 {
			result[i] = NeutralScore
			continue
		}

		// Convert to score: percentage distance from MA200
		// Price above MA200 = bullish (>50), below = bearish (<50)
		percentDiff := ((prices[i] - ma200) / ma200) * 100.0

		// Map to 0-100 scale: ±20% from MA200 maps to 0-100 range
		// +20% or more = 100, -20% or less = 0  (20% → 100, -20% → 0)
		score := 50.0 + (percentDiff * 2.5) // 2.5 = 50/20

		result[i] = clamp(score, 0, 100)
	}

	return result
}
