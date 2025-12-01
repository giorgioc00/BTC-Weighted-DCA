package indicators

// MVRM calculates the MVRM indicator (stub implementation)
// TODO: Implement actual MVRM calculation
func MVRM(prices []float64) []float64 {
	// Return neutral values (50) for now
	result := make([]float64, len(prices))
	for i := range result {
		result[i] = 50.0
	}
	return result
}
