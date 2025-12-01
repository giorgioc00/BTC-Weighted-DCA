package indicators

// MA200 calculates the 200-day moving average indicator (stub implementation)
// TODO: Implement actual MA200 calculation
func MA200(prices []float64) []float64 {
	// Return neutral values (50) for now
	result := make([]float64, len(prices))
	for i := range result {
		result[i] = 50.0
	}
	return result
}
