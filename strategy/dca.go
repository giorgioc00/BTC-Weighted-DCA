package strategy

import (
	"backtester/indicators"
	"fmt"
	"reflect"
)

type DCAConfig struct {
	RSI               RSIConfig
	Score             ScoreConfig
	Weights           WeightsConfig
	ExecutionInterval int // Interval for operation execution in days
}

type RSIConfig struct {
	BaseAmount    float64 // Base amount to invest per period
	RSIPeriod     int     // RSI calculation period (typically 14)
	RSIOversold   float64 // RSI level considered oversold (e.g., 30)
	RSIOverbought float64 // RSI level considered overbought (e.g., 70)
	MaxMultiplier float64 // Maximum multiplier for base amount
	MinMultiplier float64 // Minimum multiplier for base amount
}

type ScoreConfig struct {
	BaseAmount      float64 // Base amount to invest per period for score-based strategy
	ScoreOversold   float64 // Score threshold for "oversold" (e.g., 30)
	ScoreOverbought float64 // Score threshold for "overbought" (e.g., 70)
	MaxMultiplier   float64 // Maximum multiplier for base amount
	MinMultiplier   float64 // Minimum multiplier for base amount
}
type WeightsConfig struct {
	RSI          int
	FearAndGreed int
	MVRM         int
	Ma200        int
}

// Validate checks that the weights sum to 100
func (w *WeightsConfig) Validate() error {
	v := reflect.ValueOf(*w)
	total := 0
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Int {
			total += int(field.Int())
		}
	}
	if total != 100 {
		return fmt.Errorf("weights must sum to 100, got %d", total)
	}
	return nil
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() DCAConfig {
	return DCAConfig{
		RSI: RSIConfig{
			BaseAmount:    100.0,
			RSIPeriod:     14,
			RSIOversold:   25.0,
			RSIOverbought: 75.0,
			MaxMultiplier: 3.0,
			MinMultiplier: 0.5,
		},
		Score: ScoreConfig{
			BaseAmount:      100.0,
			ScoreOversold:   30.0,
			ScoreOverbought: 70.0,
			MaxMultiplier:   3.0,
			MinMultiplier:   0.5,
		},
		Weights: WeightsConfig{
			RSI:          40,
			FearAndGreed: 30,
			MVRM:         20,
			Ma200:        10,
		},
		ExecutionInterval: 30,
	}
}

// IndicatorSet holds all indicator values for a single point in time
type IndicatorSet struct {
	RSI          float64
	FearAndGreed float64
	MVRM         float64
	MA200        float64
}

// DynamicDCA implements a dynamic DCA strategy based on RSI
type DynamicDCA struct {
	Config   DCAConfig
	MVRVData []float64 // Pre-loaded MVRV ratios from CSV
}

// NewDynamicDCA creates a new dynamic DCA strategy
func NewDynamicDCA(config DCAConfig) *DynamicDCA {
	return &DynamicDCA{Config: config}
}

// calculateAllIndicators calculates all indicators for the entire price series
func (d *DynamicDCA) calculateAllIndicators(prices []float64) []IndicatorSet {
	rsiValues := indicators.RSI(prices, d.Config.RSI.RSIPeriod)
	fearAndGreedValues := indicators.FearAndGreed(prices)
	mvrmValues := indicators.MVRM(prices)
	ma200Values := indicators.MA200(prices)

	sets := make([]IndicatorSet, len(prices))
	for i := 0; i < len(prices); i++ {
		sets[i] = IndicatorSet{
			RSI:          rsiValues[i],
			FearAndGreed: fearAndGreedValues[i],
			MVRM:         mvrvValues[i],
			MA200:        ma200Values[i],
		}
	}
	return sets
}

// calculateScore takes an IndicatorSet and creates a weighted score
func (d *DynamicDCA) calculateScore(indicators IndicatorSet) float64 {
	score := indicators.RSI * float64(d.Config.Weights.RSI) / 100.0
	score += indicators.FearAndGreed * float64(d.Config.Weights.FearAndGreed) / 100.0
	score += indicators.MVRM * float64(d.Config.Weights.MVRM) / 100.0
	score += indicators.MA200 * float64(d.Config.Weights.Ma200) / 100.0

	return score
}

// calculateInvestmentAmount determines how much to invest based on the weighted score
func (d *DynamicDCA) calculateInvestmentAmount(score float64) float64 {
	if score <= d.Config.Score.ScoreOversold {
		// Score is oversold - invest maximum
		return d.Config.Score.BaseAmount * d.Config.Score.MaxMultiplier
	} else if score >= d.Config.Score.ScoreOverbought {
		// Score is overbought - invest minimum
		return d.Config.Score.BaseAmount * d.Config.Score.MinMultiplier
	}

	// Linear interpolation between oversold and overbought
	// Lower score = higher multiplier | Inverse linear proportion
	scoreRange := d.Config.Score.ScoreOverbought - d.Config.Score.ScoreOversold
	multiplierRange := d.Config.Score.MaxMultiplier - d.Config.Score.MinMultiplier
	scorePosition := (score - d.Config.Score.ScoreOversold) / scoreRange

	multiplier := d.Config.Score.MaxMultiplier - (scorePosition * multiplierRange)
	return d.Config.Score.BaseAmount * multiplier
}

// calculateInvestmentAmountRSI determines how much to invest based on RSI
// When RSI is low (oversold), invest more. When RSI is high (overbought), invest less.
func (d *DynamicDCA) calculateInvestmentAmountRSI(rsi float64) float64 {
	if rsi <= d.Config.RSI.RSIOversold {
		// RSI is oversold - invest maximum
		return d.Config.RSI.BaseAmount * d.Config.RSI.MaxMultiplier
	} else if rsi >= d.Config.RSI.RSIOverbought {
		// RSI is overbought - invest minimum
		return d.Config.RSI.BaseAmount * d.Config.RSI.MinMultiplier
	}

	// Linear interpolation between oversold and overbought
	// Lower RSI = higher multiplier
	rsiRange := d.Config.RSI.RSIOverbought - d.Config.RSI.RSIOversold
	multiplierRange := d.Config.RSI.MaxMultiplier - d.Config.RSI.MinMultiplier
	rsiPosition := (rsi - d.Config.RSI.RSIOversold) / rsiRange

	multiplier := d.Config.RSI.MaxMultiplier - (rsiPosition * multiplierRange)
	return d.Config.RSI.BaseAmount * multiplier
}

// GenerateSignalsRSI generates buy signals with amounts for the entire price series
func (d *DynamicDCA) GenerateSignalsRSI(prices []float64) []Signal {
	rsiValues := indicators.RSI(prices, d.Config.RSI.RSIPeriod)
	if rsiValues == nil {
		return nil
	}

	signals := make([]Signal, len(prices))
	for i := 0; i < len(prices); i++ {
		signals[i] = Signal{
			Index:  i,
			Price:  prices[i],
			RSI:    rsiValues[i],
			Amount: d.calculateInvestmentAmountRSI(rsiValues[i]),
			Action: Buy, // DCA always buys
		}
	}

	return signals
}

// GenerateSignals generates buy signals using weighted multi-indicator score
func (d *DynamicDCA) GenerateSignals(prices []float64) []Signal {
	indicatorSets := d.calculateAllIndicators(prices)

	signals := make([]Signal, len(prices))
	for i := 0; i < len(prices); i++ {
		score := d.calculateScore(indicatorSets[i])
		signals[i] = Signal{
			Index:  i,
			Price:  prices[i],
			RSI:    indicatorSets[i].RSI, // TODO - Maybe to remove
			Score:  score,
			Amount: d.calculateInvestmentAmount(score),
			Action: Buy, // DCA always buys
		}
	}

	return signals
}

// Signal represents a trading signal
type Signal struct {
	Index  int
	Price  float64
	RSI    float64
	Score  float64
	Amount float64
	Action Action
}
type Score struct {
}

// Action represents a trading action
type Action int

const (
	Hold Action = iota
	Buy
	Sell
)

func (a Action) String() string {
	switch a {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return "HOLD"
	}
}
