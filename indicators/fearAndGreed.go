package indicators

func FearAndGreed(prices []float64) []float64 {
	values := make([]float64, len(prices))

	for i := range values {
		values[i] = 50.0
	}

	return values
}
