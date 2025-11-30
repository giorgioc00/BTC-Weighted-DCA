package strategy

import (
	"backtester/indicators"
	"fmt"
)

type DCAConfig struct {
	RSI               RSIConfig
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
type WeightsConfig struct {
	RSI          int
	FearAndGreed int
	MVRM         int
	Ma200        int
}

// Validate checks that the weights sum to 100
func (w *WeightsConfig) Validate() error {
	total := w.RSI + w.FearAndGreed + w.MVRM + w.Ma200
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
		Weights: WeightsConfig{
			RSI:          40,
			FearAndGreed: 30,
			MVRM:         20,
			Ma200:        10,
		},
		ExecutionInterval: 30,
	}
}

// DynamicDCA implements a dynamic DCA strategy based on RSI
type DynamicDCA struct {
	Config DCAConfig
}

// NewDynamicDCA creates a new dynamic DCA strategy
func NewDynamicDCA(config DCAConfig) *DynamicDCA {
	return &DynamicDCA{Config: config}
}

// CalculateScore takes all the different indicators and create a weighted score
func (d *DynamicDCA) CalculateScore(rsi float64) float64 {

}

// CalculateInvestmentAmount CalculateInvestmentAmountRSI determines how much to invest based on the weighted score
func (d *DynamicDCA) CalculateInvestmentAmount(rsi float64) float64 {
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

// CalculateInvestmentAmountRSI determines how much to invest based on RSI
// When RSI is low (oversold), invest more. When RSI is high (overbought), invest less.
func (d *DynamicDCA) CalculateInvestmentAmountRSI(rsi float64) float64 {
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

// GenerateSignals generates buy signals with amounts for the entire price series
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
			Amount: d.CalculateInvestmentAmountRSI(rsiValues[i]),
			Action: Buy, // DCA always buys
		}
	}

	return signals
}

// For the ongoing weighted score
func (d *DynamicDCA) GenerateSignals(prices []float64) []Signal {
	rsiValues := indicators.RSI(prices, d.Config.RSIPeriod)
	if rsiValues == nil {
		return nil
	}

	signals := make([]Signal, len(prices))
	for i := 0; i < len(prices); i++ {
		signals[i] = Signal{
			Index:  i,
			Price:  prices[i],
			RSI:    rsiValues[i],
			Score:  0,
			Amount: d.CalculateInvestmentAmount(rsiValues[i]),
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
