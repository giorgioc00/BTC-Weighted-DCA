package signals

import (
	"fmt"
	"reflect"
)

// Config holds all strategy configuration
type Config struct {
	RSI               RSIConfig       `json:"RSI"`
	Score             ScoreConfig     `json:"Score"`
	Bollinger         BollingerConfig `json:"Bollinger"`
	ATR               ATRConfig       `json:"ATR"`
	NetFlow           NetFlowConfig   `json:"NetFlow"`
	VolumeFlow        VolumeConfig    `json:"VolumeFlow"`
	Weights           WeightsConfig   `json:"Weights"`
	ExecutionInterval int             `json:"ExecutionInterval"`
}

type RSIConfig struct {
	BaseAmount    float64
	RSIPeriod     int
	RSIOversold   float64
	RSIOverbought float64
	MaxMultiplier float64
	MinMultiplier float64
}

type ScoreConfig struct {
	BaseAmount      float64
	ScoreOversold   float64
	ScoreOverbought float64
	MaxMultiplier   float64
	MinMultiplier   float64
}

type BollingerConfig struct {
	Period     int     `json:"Period"`
	Multiplier float64 `json:"Multiplier"`
}

type ATRConfig struct {
	Period     int     `json:"Period"`
	LowVolPct  float64 `json:"LowVolPct"`
	HighVolPct float64 `json:"HighVolPct"`
}

type NetFlowConfig struct {
	Period int     `json:"Period"`
	ZCap   float64 `json:"ZCap"`
}

type VolumeConfig struct {
	Period int     `json:"Period"`
	ZCap   float64 `json:"ZCap"`
}

type WeightsConfig struct {
	RSI          int `json:"RSI"`
	FearAndGreed int `json:"FearAndGreed"`
	MVRV         int `json:"MVRV"`
	MA200        int `json:"MA200"`
	Bollinger    int `json:"Bollinger"`
	ATR          int `json:"ATR"`
	NetFlow      int `json:"NetFlow"`
	VolumeFlow   int `json:"VolumeFlow"`
}

// Validation constants
const (
	MaxReasonableMultiplier = 10.0  // Maximum multiplier to prevent absurd values
	MaxPeriod               = 365   // Maximum period (1 year of daily data)
	MaxPercentage           = 100.0 // Maximum percentage value
)

// Validation helper functions

// validatePeriod checks that a period value is within valid bounds
func validatePeriod(period int, fieldName string) error {
	if period <= 1 {
		return fmt.Errorf("%s.Period must be > 1, got %d", fieldName, period)
	}
	if period > MaxPeriod {
		return fmt.Errorf("%s.Period too large (max %d), got %d", fieldName, MaxPeriod, period)
	}
	return nil
}

// validateZCap checks that a Z-score cap value is within valid bounds
func validateZCap(zCap float64, fieldName string) error {
	if zCap <= 0 {
		return fmt.Errorf("%s.ZCap must be > 0, got %.2f", fieldName, zCap)
	}
	if zCap > 10.0 {
		return fmt.Errorf("%s.ZCap too large (max 10.0), got %.2f", fieldName, zCap)
	}
	return nil
}

// validateOversoldOverbought checks oversold/overbought threshold values
func validateOversoldOverbought(oversold, overbought float64, fieldName string) error {
	if oversold < 0 || oversold > MaxPercentage {
		return fmt.Errorf("%s.Oversold must be between 0 and %.0f, got %.2f", fieldName, MaxPercentage, oversold)
	}
	if overbought < 0 || overbought > MaxPercentage {
		return fmt.Errorf("%s.Overbought must be between 0 and %.0f, got %.2f", fieldName, MaxPercentage, overbought)
	}
	if oversold >= overbought {
		return fmt.Errorf("%s.Oversold (%.2f) must be < %s.Overbought (%.2f)", fieldName, oversold, fieldName, overbought)
	}
	return nil
}

// validateMultipliers checks min/max multiplier values
func validateMultipliers(minMult, maxMult float64, fieldName string) error {
	if minMult < 0 {
		return fmt.Errorf("%s.MinMultiplier must be >= 0, got %.2f", fieldName, minMult)
	}
	if maxMult <= 0 {
		return fmt.Errorf("%s.MaxMultiplier must be > 0, got %.2f", fieldName, maxMult)
	}
	if maxMult > MaxReasonableMultiplier {
		return fmt.Errorf("%s.MaxMultiplier too large (max %.1f), got %.2f", fieldName, MaxReasonableMultiplier, maxMult)
	}
	if minMult > maxMult {
		return fmt.Errorf("%s.MinMultiplier (%.2f) must be <= %s.MaxMultiplier (%.2f)", fieldName, minMult, fieldName, maxMult)
	}
	return nil
}

// Validate checks config values for consistency.
func (c *Config) Validate() error {
	// Validate weights first
	if err := c.Weights.Validate(); err != nil {
		return err
	}

	// ExecutionInterval validation
	if c.ExecutionInterval <= 0 {
		return fmt.Errorf("ExecutionInterval must be > 0, got %d", c.ExecutionInterval)
	}

	// RSI configuration validation
	if c.RSI.BaseAmount <= 0 {
		return fmt.Errorf("RSI.BaseAmount must be > 0, got %.2f", c.RSI.BaseAmount)
	}
	if c.RSI.RSIPeriod <= 1 {
		return fmt.Errorf("RSI.RSIPeriod must be > 1, got %d", c.RSI.RSIPeriod)
	}
	if c.RSI.RSIPeriod > MaxPeriod {
		return fmt.Errorf("RSI.RSIPeriod too large (max %d), got %d", MaxPeriod, c.RSI.RSIPeriod)
	}
	if err := validateOversoldOverbought(c.RSI.RSIOversold, c.RSI.RSIOverbought, "RSI"); err != nil {
		return err
	}
	if err := validateMultipliers(c.RSI.MinMultiplier, c.RSI.MaxMultiplier, "RSI"); err != nil {
		return err
	}

	// Score configuration validation
	if c.Score.BaseAmount <= 0 {
		return fmt.Errorf("Score.BaseAmount must be > 0, got %.2f", c.Score.BaseAmount)
	}
	if err := validateOversoldOverbought(c.Score.ScoreOversold, c.Score.ScoreOverbought, "Score"); err != nil {
		return err
	}
	if err := validateMultipliers(c.Score.MinMultiplier, c.Score.MaxMultiplier, "Score"); err != nil {
		return err
	}

	// Bollinger configuration validation
	if err := validatePeriod(c.Bollinger.Period, "Bollinger"); err != nil {
		return err
	}
	if c.Bollinger.Multiplier <= 0 {
		return fmt.Errorf("Bollinger.Multiplier must be > 0, got %.2f", c.Bollinger.Multiplier)
	}
	if c.Bollinger.Multiplier > 5.0 {
		return fmt.Errorf("Bollinger.Multiplier too large (max 5.0), got %.2f", c.Bollinger.Multiplier)
	}

	// ATR configuration validation
	if err := validatePeriod(c.ATR.Period, "ATR"); err != nil {
		return err
	}
	if c.ATR.LowVolPct <= 0 || c.ATR.LowVolPct > MaxPercentage {
		return fmt.Errorf("ATR.LowVolPct must be between 0 and %.0f, got %.2f", MaxPercentage, c.ATR.LowVolPct)
	}
	if c.ATR.HighVolPct <= 0 || c.ATR.HighVolPct > MaxPercentage {
		return fmt.Errorf("ATR.HighVolPct must be between 0 and %.0f, got %.2f", MaxPercentage, c.ATR.HighVolPct)
	}
	if c.ATR.LowVolPct >= c.ATR.HighVolPct {
		return fmt.Errorf("ATR.LowVolPct (%.2f) must be < ATR.HighVolPct (%.2f)", c.ATR.LowVolPct, c.ATR.HighVolPct)
	}

	// NetFlow configuration validation
	if err := validatePeriod(c.NetFlow.Period, "NetFlow"); err != nil {
		return err
	}
	if err := validateZCap(c.NetFlow.ZCap, "NetFlow"); err != nil {
		return err
	}

	// VolumeFlow configuration validation
	if err := validatePeriod(c.VolumeFlow.Period, "VolumeFlow"); err != nil {
		return err
	}
	if err := validateZCap(c.VolumeFlow.ZCap, "VolumeFlow"); err != nil {
		return err
	}

	return nil
}

// Validate checks that weights sum to 100 and all are non-negative
func (w *WeightsConfig) Validate() error {
	v := reflect.ValueOf(*w)
	t := v.Type()
	total := 0

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Int {
			weight := int(field.Int())

			// Check for negative weights
			if weight < 0 {
				fieldName := t.Field(i).Name
				return fmt.Errorf("weight %s must be >= 0, got %d", fieldName, weight)
			}

			total += weight
		}
	}

	if total != 100 {
		return fmt.Errorf("weights must sum to 100, got %d", total)
	}

	return nil
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
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
		Bollinger: BollingerConfig{
			Period:     20,
			Multiplier: 2.0,
		},
		ATR: ATRConfig{
			Period:     14,
			LowVolPct:  2.0,
			HighVolPct: 10.0,
		},
		NetFlow: NetFlowConfig{
			Period: 30,
			ZCap:   3.0,
		},
		VolumeFlow: VolumeConfig{
			Period: 20,
			ZCap:   3.0,
		},
		Weights: WeightsConfig{
			RSI:          25,
			FearAndGreed: 15,
			MVRV:         15,
			MA200:        10,
			Bollinger:    10,
			ATR:          10,
			NetFlow:      10,
			VolumeFlow:   5,
		},
		ExecutionInterval: 30,
	}
}
