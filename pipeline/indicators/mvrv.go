package indicators

// MVRV converts raw MVRV ratio values to a 0-100 score
// Pure function: low MVRV = high score (undervalued)
func MVRV(mvrvRatios []float64) []float64 {
	result := make([]float64, len(mvrvRatios))

	for i, ratio := range mvrvRatios {
		var score float64
		switch {
		case ratio <= 0.5:
			score = 100
		case ratio >= 5.0:
			score = 0
		case ratio <= 1.0:
			score = lerpWithClamp(ratio, 0.5, 1.0, 100, 75)
		case ratio <= 2.0:
			score = lerpWithClamp(ratio, 1.0, 2.0, 75, 40)
		case ratio <= 3.5:
			score = lerpWithClamp(ratio, 2.0, 3.5, 40, 10)
		default:
			score = lerpWithClamp(ratio, 3.5, 5.0, 10, 0)
		}

		result[i] = clamp(score, 0, 100)
	}

	return result
}
