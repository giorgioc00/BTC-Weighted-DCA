package indicators

// BollingerPercentB calculates Bollinger %B score
// Pure function: price at lower band = 100, upper band = 0
func BollingerPercentB(prices []float64, period int, multiplier float64) []float64 {
	scores := newNeutralSlice(len(prices))

	if len(prices) < period || period <= 1 {
		return scores
	}

	var sum float64
	for i := 0; i < period; i++ {
		sum += prices[i]
	}

	window := prices[:period]
	std := stdDev(window, sum)

	for i := 0; i < len(prices); i++ {
		if i < period-1 {
			scores[i] = NeutralScore
			continue
		}

		if i >= period {
			sum += prices[i] - prices[i-period]
			std = stdDev(prices[i-period+1:i+1], sum)
		}

		sma := sum / float64(period)
		lower := sma - multiplier*std
		upper := sma + multiplier*std
		width := upper - lower

		if width == 0 {
			scores[i] = NeutralScore
			continue
		}

		percentB := (prices[i] - lower) / width
		percentB = clamp(percentB, 0, 1)

		// Invert: lower band = high score
		score := (1 - percentB) * 100
		scores[i] = clamp(score, 0, 100)
	}

	return scores
}
