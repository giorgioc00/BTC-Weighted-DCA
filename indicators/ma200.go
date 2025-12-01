package indicators

// MA200 calculates the 200-day Simple Moving Average
// Returns a score from 0-100 based on price position relative to MA200:
// - Price > MA200: bullish signal (higher score)
// - Price < MA200: bearish signal (lower score)
func MA200(prices []float64) []float64 {
	period := 200
	result := make([]float64, len(prices))

	// Return neutral value (50) for insufficient data
	if len(prices) < period {
		for i := range result {
			result[i] = 50.0
		}
		return result
	}

	// Fill first period-1 values with neutral
	for i := 0; i < period-1; i++ {
		result[i] = 50.0
	}

	// Calculate MA200 and convert to 0-100 score
	for i := period - 1; i < len(prices); i++ {
		// Calculate simple moving average
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i-j]
		}
		ma200 := sum / float64(period)

		// Convert to score: percentage distance from MA200
		// Price above MA200 = bullish (>50), below = bearish (<50)
		percentDiff := ((prices[i] - ma200) / ma200) * 100.0

		// Map to 0-100 scale: ±20% from MA200 maps to 0-100 range
		// +20% or more = 100, -20% or less = 0  (20% → 100, -20% → 0)
		score := 50.0 + (percentDiff * 2.5) // 2.5 = 50/20

		// Clamp to 0-100
		if score > 100 {
			score = 100
		} else if score < 0 {
			score = 0
		}

		result[i] = score
	}

	return result
}
