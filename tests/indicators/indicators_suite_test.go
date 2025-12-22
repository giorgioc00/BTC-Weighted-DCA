package indicators_test

import (
	"math"
	"testing"

	"backtester/pipeline/indicators"
)

func TestClampAndLinearInterp(t *testing.T) {
	if indicators.Clamp(150, 0, 100) != 100 {
		t.Fatalf("Clamp did not cap upper bound")
	}
	val := indicators.LinearInterp(0.75, 0.5, 1.0, 100, 75)
	if math.Abs(val-87.5) > 0.01 {
		t.Fatalf("LinearInterp unexpected result: %.2f", val)
	}
}

func TestMA200_NeutralOnShortSeries(t *testing.T) {
	prices := make([]float64, 50)
	for i := range prices {
		prices[i] = 100
	}
	scores := indicators.MA200(prices)
	for _, v := range scores {
		if v != indicators.NeutralScore {
			t.Fatalf("expected neutral score for insufficient data, got %.2f", v)
		}
	}
}

func TestBollingerPercentB_ScoreRange(t *testing.T) {
	prices := []float64{10, 10, 10, 10, 10, 9}
	scores := indicators.BollingerPercentB(prices, 5, 2.0)
	if scores[len(scores)-1] <= 50 {
		t.Fatalf("expected higher score near lower band, got %.2f", scores[len(scores)-1])
	}
}

func TestMVRV_MappingAnchors(t *testing.T) {
	input := []float64{0.5, 1.0, 2.0, 3.5, 5.0}
	want := []float64{100, 75, 40, 10, 0}
	got := indicators.MVRV(input)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d got %.2f want %.2f", i, got[i], want[i])
		}
	}
}

func TestNetFlow_OutflowsRaiseScore(t *testing.T) {
	in := []float64{100, 100, 100, 100, 90, 80, 70}
	out := []float64{50, 50, 50, 50, 110, 120, 130}
	scores := indicators.ExchangeNetFlowScore(in, out, 3, 3)
	if scores[len(scores)-1] <= 50 {
		t.Fatalf("expected high score for outflows, got %.2f", scores[len(scores)-1])
	}
}

func TestVolumeFlow_PositiveTrend(t *testing.T) {
	close := []float64{10, 11, 12, 13, 14, 15}
	vol := []float64{100, 200, 300, 400, 500, 600}
	scores := indicators.VolumeFlowScore(close, vol, 3, 3)
	if scores[len(scores)-1] <= 50 {
		t.Fatalf("expected high score for positive price/volume trend, got %.2f", scores[len(scores)-1])
	}
}

func TestRSI_Bounds(t *testing.T) {
	prices := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	rsi := indicators.RSI(prices, 14)
	if rsi == nil || len(rsi) != len(prices) {
		t.Fatalf("RSI returned invalid length")
	}
	if rsi[14] != 100 {
		t.Fatalf("expected RSI 100 on all gains, got %.2f", rsi[14])
	}
}
