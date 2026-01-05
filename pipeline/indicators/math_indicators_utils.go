package indicators

import "math"

// NeutralScore is the default score returned when insufficient data exists
const NeutralScore = 50.0

// clamp restricts a value to the range [min, max]
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// lerpWithClamp performs linear interpolation between two points
// Maps x from range [x1,x2] to range [y1,y2]
func lerpWithClamp(
	currentX float64,
	startX float64,
	endX float64,
	startY float64,
	endY float64,
) float64 {
	if currentX <= startX {
		return startY
	}
	if currentX >= endX {
		return endY
	}
	progressRatio := (currentX - startX) / (endX - startX)
	yRange := endY - startY
	interpolatedY := startY + progressRatio*yRange

	return interpolatedY
}

// fillSlice fills a slice with a constant value from start to end (exclusive)
func fillSlice(slice []float64, start, end int, value float64) {
	for i := start; i < end && i < len(slice); i++ {
		slice[i] = value
	}
}

// newNeutralSlice creates a slice filled with NeutralScore
func newNeutralSlice(length int) []float64 {
	result := make([]float64, length)
	fillSlice(result, 0, length, NeutralScore)
	return result
}

// meanStd calculates mean and standard deviation of a slice
func meanStd(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	return mean, math.Sqrt(variance)
}

// stdDev calculates standard deviation given precomputed sum
// This is an optimization for sliding windows where sum is already known
func stdDev(window []float64, sum float64) float64 {
	if len(window) == 0 {
		return 0
	}
	mean := sum / float64(len(window))
	var variance float64
	for _, v := range window {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(window))
	return math.Sqrt(variance)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Exported wrappers for testing from external packages.
func Clamp(value, min, max float64) float64 { return clamp(value, min, max) }
func LinearInterp(x, x1, x2, y1, y2 float64) float64 {
	return lerpWithClamp(x, x1, x2, y1, y2)
}
func FillSlice(slice []float64, start, end int, value float64) { fillSlice(slice, start, end, value) }
func NewNeutralSlice(length int) []float64                     { return newNeutralSlice(length) }
func MeanStd(values []float64) (float64, float64)              { return meanStd(values) }
func StdDev(window []float64, sum float64) float64             { return stdDev(window, sum) }
