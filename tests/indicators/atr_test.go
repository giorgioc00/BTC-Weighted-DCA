package indicators_test

import (
	"testing"

	"backtester/pipeline/indicators"
)

func TestATRVolatilityScore_Mapping(t *testing.T) {
	close := []float64{10, 10, 10, 10, 10, 10, 10, 10, 10}
	high := []float64{10, 10.5, 10.2, 10.1, 10.3, 10.4, 10.2, 10.1, 10.05}
	low := []float64{10, 9.8, 9.9, 9.95, 9.9, 9.8, 9.9, 9.95, 9.9}

	scores := indicators.ATRVolatilityScore(close, high, low, 3, 1, 5)
	if len(scores) != len(close) {
		t.Fatalf("length mismatch: got %d want %d", len(scores), len(close))
	}
	if scores[0] != indicators.NeutralScore {
		t.Fatalf("expected neutral score with insufficient data, got %.2f", scores[0])
	}
	if scores[len(scores)-1] < 55 {
		t.Fatalf("expected higher score for low ATR, got %.2f", scores[len(scores)-1])
	}
}
