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

// Validate checks that weights sum to 100
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
