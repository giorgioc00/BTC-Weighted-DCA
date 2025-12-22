package indicators

// VolumeFlowScore computes signed volume z-score
// Pure function: positive volume with positive returns = high score
func VolumeFlowScore(close []float64, volume []float64, period int, zCap float64) []float64 {
	length := min(len(close), len(volume))
	scores := newNeutralSlice(length)

	if length == 0 || period <= 1 {
		return scores
	}

	signed := make([]float64, length)
	for i := 1; i < length; i++ {
		switch {
		case close[i] > close[i-1]:
			signed[i] = volume[i]
		case close[i] < close[i-1]:
			signed[i] = -volume[i]
		default:
			signed[i] = 0
		}
	}

	for i := 0; i < length; i++ {
		if i < period {
			scores[i] = NeutralScore
			continue
		}
		window := signed[i-period : i]
		mean, std := meanStd(window)
		if std == 0 {
			scores[i] = NeutralScore
			continue
		}
		z := (signed[i] - mean) / std
		z = clamp(z, -zCap, zCap)

		// Map [-zCap, zCap] to [0,100]
		pos := (z + zCap) / (2 * zCap)
		scores[i] = pos * 100
	}
	return scores
}
