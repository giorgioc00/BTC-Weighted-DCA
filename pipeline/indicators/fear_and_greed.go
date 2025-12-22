package indicators

// FearAndGreed is a placeholder that returns neutral values
// Pure function: []float64 → []float64
func FearAndGreed(prices []float64) []float64 {
	return newNeutralSlice(len(prices))
}
