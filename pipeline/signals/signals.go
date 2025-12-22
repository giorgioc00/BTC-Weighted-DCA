package signals

import (
	"log"
	"time"

	"backtester/pipeline/indicators"
	"backtester/util"
)

// Signal represents a trading signal at a point in time
type Signal struct {
	Index  int
	Price  float64
	Score  float64
	Amount float64
	// Individual indicator values
	RSI          float64
	FearAndGreed float64
	MVRV         float64
	MA200        float64
	Bollinger    float64
	ATR          float64
	NetFlow      float64
	VolumeFlow   float64
}

// MarketData holds all input data for signal generation
type MarketData struct {
	Prices         []float64
	Highs          []float64
	Lows           []float64
	Volumes        []float64
	MVRVData       []float64
	FlowInUSD      []float64
	FlowOutUSD     []float64
	FirstPriceDate time.Time
}

// IndicatorSet holds all indicator values at a point in time
type IndicatorSet struct {
	RSI          float64
	FearAndGreed float64
	MVRV         float64
	MA200        float64
	Bollinger    float64
	ATR          float64
	NetFlow      float64
	VolumeFlow   float64
}

// GenerateSignals creates trading signals using multi-indicator scoring
// This is a pure pipeline: data → indicators → scores → signals
func GenerateSignals(data MarketData, config Config) []Signal {
	indicators := calculateAllIndicators(data, config)

	signals := make([]Signal, len(data.Prices))
	for i := 0; i < len(data.Prices); i++ {
		score := calculateScore(indicators[i], config)
		signals[i] = Signal{
			Index:        i,
			Price:        data.Prices[i],
			Score:        score,
			Amount:       calculateInvestmentAmount(score, config),
			RSI:          indicators[i].RSI,
			FearAndGreed: indicators[i].FearAndGreed,
			MVRV:         indicators[i].MVRV,
			MA200:        indicators[i].MA200,
			Bollinger:    indicators[i].Bollinger,
			ATR:          indicators[i].ATR,
			NetFlow:      indicators[i].NetFlow,
			VolumeFlow:   indicators[i].VolumeFlow,
		}
	}

	return signals
}

// GenerateSignalsRSI creates signals using RSI-only mode (legacy)
func GenerateSignalsRSI(data MarketData, config Config) []Signal {
	rsiValues := indicators.RSI(data.Prices, config.RSI.RSIPeriod)
	if rsiValues == nil {
		return nil
	}

	signals := make([]Signal, len(data.Prices))
	for i := 0; i < len(data.Prices); i++ {
		signals[i] = Signal{
			Index:        i,
			Price:        data.Prices[i],
			Score:        rsiValues[i],
			Amount:       calculateInvestmentAmountRSI(rsiValues[i], config),
			RSI:          rsiValues[i],
			FearAndGreed: 50.0,
			MVRV:         50.0,
			MA200:        50.0,
			Bollinger:    50.0,
			ATR:          50.0,
			NetFlow:      50.0,
			VolumeFlow:   50.0,
		}
	}

	return signals
}

// calculateAllIndicators computes all indicator values
func calculateAllIndicators(data MarketData, config Config) []IndicatorSet {
	prices := data.Prices
	highs := alignSeries(data.Highs, len(prices), 0)
	lows := alignSeries(data.Lows, len(prices), 0)
	volumes := alignSeries(data.Volumes, len(prices), 0)
	flowsIn := alignSeries(data.FlowInUSD, len(prices), 0)
	flowsOut := alignSeries(data.FlowOutUSD, len(prices), 0)

	rsiValues := indicators.RSI(prices, config.RSI.RSIPeriod)
	fearAndGreedValues := indicators.FearAndGreed(prices)
	bollingerValues := indicators.BollingerPercentB(prices, config.Bollinger.Period, config.Bollinger.Multiplier)
	atrValues := indicators.ATRVolatilityScore(prices, highs, lows, config.ATR.Period, config.ATR.LowVolPct, config.ATR.HighVolPct)
	volumeFlowValues := indicators.VolumeFlowScore(prices, volumes, config.VolumeFlow.Period, config.VolumeFlow.ZCap)
	netFlowValues := indicators.ExchangeNetFlowScore(flowsIn, flowsOut, config.NetFlow.Period, config.NetFlow.ZCap)

	// Handle MVRV data with validation
	var mvrvValues []float64
	if len(data.MVRVData) >= len(prices) && len(data.MVRVData) > 0 {
		aligned := alignSeries(data.MVRVData, len(prices), 50.0)
		if len(data.MVRVData) == len(prices) {
			firstMVRVDate, err := util.LoadLastCSVColumnRow("data/btc_mvrv.csv", "Date")
			if err != nil {
				log.Printf("Warning: Could not validate MVRV date: %v\n", err)
				mvrvValues = indicators.MVRV(aligned)
			} else {
				mvrvDate, err := util.ParseDateFlexible(firstMVRVDate)
				if err != nil || data.FirstPriceDate.IsZero() || !mvrvDate.Equal(data.FirstPriceDate) {
					log.Printf("Warning: MVRV data date mismatch (expected %s, got %s) - using neutral values\n",
						data.FirstPriceDate.Format("2006-01-02"), firstMVRVDate)
					mvrvValues = make([]float64, len(prices))
					for i := range mvrvValues {
						mvrvValues[i] = 50.0
					}
				} else {
					mvrvValues = indicators.MVRV(aligned)
				}
			}
		} else {
			mvrvValues = indicators.MVRV(aligned)
		}
	} else {
		mvrvValues = make([]float64, len(prices))
		for i := range mvrvValues {
			mvrvValues[i] = 50.0
		}
	}

	ma200Values := indicators.MA200(prices)

	sets := make([]IndicatorSet, len(prices))
	for i := 0; i < len(prices); i++ {
		sets[i] = IndicatorSet{
			RSI:          rsiValues[i],
			FearAndGreed: fearAndGreedValues[i],
			MVRV:         mvrvValues[i],
			MA200:        ma200Values[i],
			Bollinger:    bollingerValues[i],
			ATR:          atrValues[i],
			NetFlow:      netFlowValues[i],
			VolumeFlow:   volumeFlowValues[i],
		}
	}
	return sets
}

// calculateScore computes weighted composite score
func calculateScore(set IndicatorSet, config Config) float64 {
	score := set.RSI * float64(config.Weights.RSI) / 100.0
	score += set.FearAndGreed * float64(config.Weights.FearAndGreed) / 100.0
	score += set.MVRV * float64(config.Weights.MVRV) / 100.0
	score += set.MA200 * float64(config.Weights.MA200) / 100.0
	score += set.Bollinger * float64(config.Weights.Bollinger) / 100.0
	score += set.ATR * float64(config.Weights.ATR) / 100.0
	score += set.NetFlow * float64(config.Weights.NetFlow) / 100.0
	score += set.VolumeFlow * float64(config.Weights.VolumeFlow) / 100.0
	return score
}

// calculateInvestmentAmount determines investment based on composite score
func calculateInvestmentAmount(score float64, config Config) float64 {
	if score <= config.Score.ScoreOversold {
		return config.Score.BaseAmount * config.Score.MaxMultiplier
	} else if score >= config.Score.ScoreOverbought {
		return config.Score.BaseAmount * config.Score.MinMultiplier
	}

	// Linear interpolation
	scoreRange := config.Score.ScoreOverbought - config.Score.ScoreOversold
	multiplierRange := config.Score.MaxMultiplier - config.Score.MinMultiplier
	scorePosition := (score - config.Score.ScoreOversold) / scoreRange

	multiplier := config.Score.MaxMultiplier - (scorePosition * multiplierRange)
	return config.Score.BaseAmount * multiplier
}

// calculateInvestmentAmountRSI determines investment based on RSI only
func calculateInvestmentAmountRSI(rsi float64, config Config) float64 {
	if rsi <= config.RSI.RSIOversold {
		return config.RSI.BaseAmount * config.RSI.MaxMultiplier
	} else if rsi >= config.RSI.RSIOverbought {
		return config.RSI.BaseAmount * config.RSI.MinMultiplier
	}

	rsiRange := config.RSI.RSIOverbought - config.RSI.RSIOversold
	multiplierRange := config.RSI.MaxMultiplier - config.RSI.MinMultiplier
	rsiPosition := (rsi - config.RSI.RSIOversold) / rsiRange

	multiplier := config.RSI.MaxMultiplier - (rsiPosition * multiplierRange)
	return config.RSI.BaseAmount * multiplier
}

// alignSeries aligns a data series to target length
func alignSeries(values []float64, targetLen int, fill float64) []float64 {
	if targetLen <= 0 {
		return []float64{}
	}
	if len(values) >= targetLen {
		return values[len(values)-targetLen:]
	}
	aligned := make([]float64, targetLen)
	offset := targetLen - len(values)
	for i := 0; i < offset; i++ {
		aligned[i] = fill
	}
	copy(aligned[offset:], values)
	return aligned
}
