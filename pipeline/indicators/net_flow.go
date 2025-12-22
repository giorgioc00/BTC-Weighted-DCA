package indicators

// ExchangeNetFlowScore computes net flow z-score
// Pure function: outflows (negative) = high score (accumulation)
func ExchangeNetFlowScore(flowIn, flowOut []float64, period int, zCap float64) []float64 {
	length := min(len(flowIn), len(flowOut))
	scores := newNeutralSlice(length)

	if length == 0 || period <= 1 {
		return scores
	}

	net := make([]float64, length)
	for i := 0; i < length; i++ {
		net[i] = flowIn[i] - flowOut[i]
	}

	for i := 0; i < length; i++ {
		if i < period {
			scores[i] = NeutralScore
			continue
		}
		window := net[i-period : i]
		mean, std := meanStd(window)
		if std == 0 {
			scores[i] = NeutralScore
			continue
		}
		z := (net[i] - mean) / std
		z = clamp(z, -zCap, zCap)

		// Map [-zCap, zCap] to [100,0]
		pos := (z + zCap) / (2 * zCap)
		scores[i] = 100 - (pos * 100)
	}
	return scores
}
