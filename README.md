# Backtester

CLI backtester for a dynamic DCA strategy combining multiple indicators (RSI, MA200, Bollinger %B, ATR volatility, MVRV, exchange net-flow, volume-flow, and a placeholder fear/greed score).

## Quick Start
- `go run ./cmd/backtester` — load bundled CSVs and config, run the backtest, print summary and standard-DCA comparison.
- `go test ./...` or `go test -race ./...` — run unit tests.
- Data lives in `data/`:
  - `btc_price.csv` with `Date`, `Price`, `High`, `Low`, `Vol.` (newest-first; loaders reverse to oldest-first).
  - `btc_mvrv.csv` with `CapMVRVCur`, `FlowInExUSD`, `FlowOutExUSD` (oldest-first).

## Architecture
- `cmd/backtester/`: CLI entrypoint. Wires config, data loaders, and the pipeline.
- `pipeline/indicators/`: Pure indicator math (`[]float64 -> []float64`) for RSI, MA200, Bollinger %B, ATR, MVRV mapping, net-flow z-score, volume-flow, and helpers.
- `pipeline/signals/`: Aggregates market data, computes indicator scores, and produces `[]Signal` from `MarketData + Config`.
- `pipeline/execution/`: Executes signals at configured intervals, producing `[]Trade` and computing results.
- `pipeline/metrics/`: Holds result types; `report/` owns summary/trade printers.
- `data/`: Price loader with flexible date parsing.
- `util/`: CSV/JSON loaders, volume parsing (K/M/B), date parsing, config helpers.
- `tests/`: External-package tests for indicators and util.

## Config
- `config/strategy_config.json` holds nested sections:
  - `RSI`, `Score` (investment sizing), `Bollinger`, `ATR`, `NetFlow`, `VolumeFlow`, and `Weights`.
  - Adjust weights to sum to 100; `WeightsConfig.Validate` enforces this.
- `config/weights.json` mirrors weights if you keep a separate file (optional; `strategy_config.json` is authoritative).

## Data Handling
- Price CSV is parsed with flexible dates, commas stripped, volumes parsed with K/M/B suffixes, and reversed to chronological order when requested.
- MVRV/flow series are tail-aligned to the price series length with neutral fills when data is short; exact date matching uses flexible parsing when lengths match.
- Missing or invalid numeric values default to neutral scores (e.g., 50) to avoid crashes.

## Indicators (0-100 Scores)
- RSI (period, oversold/overbought, min/max multipliers).
- MA200 distance.
- Bollinger %B (period, std multiplier; lower band → high score).
- ATR volatility (period; low-vol/high-vol thresholds map to 100→0).
- MVRV mapping (piecewise 0.5→5.0 to 100→0).
- Exchange Net Flow z-score (FlowInExUSD - FlowOutExUSD; outflows → higher score).
- Volume Flow z-score (signed by price change; rising price + volume → higher score).
- Fear & Greed placeholder returns neutrals.

## Development Notes
- Format with `gofmt -w .`; vet with `go vet ./...`.
- Add new indicators under `pipeline/indicators` and wire into `pipeline/signals.calculateAllIndicators` and `calculateScore`.
- Place new tests under `tests/` mirroring package names; use `testdata/` for fixtures that differ from production.

## Limitations / Next Steps
- Data alignment is tail-based; for stricter correctness, add date-indexed merging for MVRV/flows vs. price.
- Fear & Greed is a stub; implement or drop its weight to 0.
- Consider command-line flags for data/config paths to avoid cwd assumptions.
