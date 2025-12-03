package indicators

// MVRV converts raw MVRV ratio values to a 0-100 score for strategy use
// MVRV (Market Value to Realized Value) interpretation:
// - Low MVRV (<1.0): Undervalued, higher score (buy signal)
// - High MVRV (>3.5): Overvalued, lower score (sell signal)
// Returns a score from 0-100 where higher = better buying opportunity
func MVRV(mvrvRatios []float64) []float64 {
	result := make([]float64, len(mvrvRatios))

	for i, ratio := range mvrvRatios {
		// Map MVRV ratio to 0-100 score (inverted: low ratio = high score)
		// MVRV <= 0.5: 100 (extreme undervaluation, strong buy)
		// MVRV = 1.0: 75 (fair value)
		// MVRV = 2.0: 40 (moderately overvalued)
		// MVRV = 3.5: 10 (highly overvalued)
		// MVRV >= 5.0: 0 (extreme overvaluation, strong sell)

		var score float64
		if ratio <= 0.5 {
			score = 100
		} else if ratio >= 5.0 {
			score = 0
		} else {
			// Piecewise linear mapping for better granularity
			if ratio <= 1.0 {
				// 0.5->100, 1.0->75
				score = 100 - ((ratio - 0.5) / 0.5 * 25)
			} else if ratio <= 2.0 {
				// 1.0->75, 2.0->40
				score = 75 - ((ratio - 1.0) / 1.0 * 35)
			} else if ratio <= 3.5 {
				// 2.0->40, 3.5->10
				score = 40 - ((ratio - 2.0) / 1.5 * 30)
			} else {
				// 3.5->10, 5.0->0
				score = 10 - ((ratio - 3.5) / 1.5 * 10)
			}
		}

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
