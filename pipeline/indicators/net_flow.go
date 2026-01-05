package indicators

// ExchangeNetFlowScore computes net flow z-score
// Pure function: outflows (negative) = high score (accumulation)
func ExchangeNetFlowScore(flowIn, flowOut []float64, period int, zCap float64) []float64 {
	length := min(len(flowIn), len(flowOut))

	if length == 0 || period <= 1 {
		return newNeutralSlice(length)
	}

	net := make([]float64, length)
	for i := 0; i < length; i++ {
		net[i] = flowIn[i] - flowOut[i]
	}

	return normalizeZScore(net, period, zCap, true)
}
