package strategy

import "backtester/indicators"

// DCAConfig holds the configuration for the dynamic DCA strategy
type DCAConfig struct {
	ExecutionInterval int     // Interval for operation execution in days
	BaseAmount        float64 // Base amount to invest per period
	RSIPeriod         int     // RSI calculation period (typically 14)
	RSIOversold       float64 // RSI level considered oversold (e.g., 30)
	RSIOverbought     float64 // RSI level considered overbought (e.g., 70)
	MaxMultiplier     float64 // Maximum multiplier for base amount
	MinMultiplier     float64 // Minimum multiplier for base amount
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() DCAConfig {
	return DCAConfig{
		ExecutionInterval: 30,
		BaseAmount:        100.0,
		RSIPeriod:         14,
		RSIOversold:       25.0,
		RSIOverbought:     75.0,
		MaxMultiplier:     3.0,
		MinMultiplier:     0.5,
	}
}

// DynamicDCA implements a dynamic DCA strategy based on RSI
type DynamicDCA struct {
	Config DCAConfig
}

// NewDynamicDCA creates a new dynamic DCA strategy
func NewDynagmicDCA(config DCAConfig) *DynamicDCA {
	return &DynamicDCA{Config: config}
}

// CalculateInvestmentAmount determines how much to invest based on RSI
// When RSI is low (oversold), invest more. When RSI is high (overbought), invest less.
func (d *DynamicDCA) CalculateInvestmentAmount(rsi float64) float64 {
	if rsi <= d.Config.RSIOversold {
		// RSI is oversold - invest maximum
		return d.Config.BaseAmount * d.Config.MaxMultiplier
	} else if rsi >= d.Config.RSIOverbought {
		// RSI is overbought - invest minimum
		return d.Config.BaseAmount * d.Config.MinMultiplier
	}

	// Linear interpolation between oversold and overbought
	// Lower RSI = higher multiplier
	rsiRange := d.Config.RSIOverbought - d.Config.RSIOversold
	multiplierRange := d.Config.MaxMultiplier - d.Config.MinMultiplier
	rsiPosition := (rsi - d.Config.RSIOversold) / rsiRange

	multiplier := d.Config.MaxMultiplier - (rsiPosition * multiplierRange)
	return d.Config.BaseAmount * multiplier
}

// GenerateSignals generates buy signals with amounts for the entire price series
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
	Amount float64
	Action Action
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
