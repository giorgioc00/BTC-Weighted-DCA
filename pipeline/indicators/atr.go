package indicators

import "math"

// ATRVolatilityScore calculates ATR-based volatility score
// Pure function: low volatility = high score
func ATRVolatilityScore(close, high, low []float64, period int, lowPct, highPct float64) []float64 {
	length := len(close)
	scores := newNeutralSlice(length)

	if period <= 1 || len(high) != length || len(low) != length || length <= period {
		return scores
	}

	tr := make([]float64, length)
	for i := 1; i < length; i++ {
		h := high[i]
		l := low[i]
		cp := close[i-1]
		range1 := h - l
		range2 := math.Abs(h - cp)
		range3 := math.Abs(l - cp)
		tr[i] = math.Max(range1, math.Max(range2, range3))
	}

	// Wilder's smoothing
	var atr float64
	for i := 1; i <= period; i++ {
		atr += tr[i]
	}
	atr /= float64(period)

	for i := 0; i < length; i++ {
		if i < period {
			scores[i] = NeutralScore
			continue
		}
		if i > period {
			atr = (atr*float64(period-1) + tr[i]) / float64(period)
		}

		if close[i] == 0 {
			scores[i] = NeutralScore
			continue
		}
		atrPct := (atr / close[i]) * 100

		scores[i] = lerpWithClamp(atrPct, lowPct, highPct, 100, 0)
	}

	return scores
}
