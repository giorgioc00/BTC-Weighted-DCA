# Backtester Architecture & Flow

This document provides visual diagrams showing how the backtester works.

## High-Level System Flow

```mermaid
graph TB
    Start([Start Application]) --> LoadConfig[Load Config<br/>strategy_config.json]
    LoadConfig --> LoadData[Load CSV Data<br/>Prices, MVRV, Flows]
    LoadData --> PrepareMarket[Prepare MarketData Struct]
    PrepareMarket --> GenSignals[Generate Signals<br/>signals.GenerateSignals]
    GenSignals --> Execute[Execute Trades<br/>execution.Execute]
    Execute --> CalcResults[Calculate Results<br/>execution.CalculateResult]
    CalcResults --> PrintReport[Print Summary<br/>report.PrintSummary]
    PrintReport --> CompareDCA[Compare with Standard DCA]
    CompareDCA --> End([End])

    style Start stroke:#00ff00,stroke-width:3px
    style End stroke:#ff4444,stroke-width:3px
    style GenSignals stroke:#00bfff,stroke-width:2px
    style Execute stroke:#00bfff,stroke-width:2px
    style CalcResults stroke:#00bfff,stroke-width:2px
```

## Detailed Data Transformation Pipeline

```mermaid
graph LR
    subgraph "1. Data Loading"
        CSV1[btc_price.csv] --> Prices[Prices Array]
        CSV1 --> Highs[Highs Array]
        CSV1 --> Lows[Lows Array]
        CSV1 --> Volumes[Volumes Array]
        CSV2[btc_mvrv.csv] --> MVRV[MVRV Data]
        CSV2 --> FlowIn[Flow In USD]
        CSV2 --> FlowOut[Flow Out USD]
    end

    subgraph "2. Signal Generation"
        Prices --> Indicators[Calculate All Indicators<br/>RSI, MA200, Bollinger, etc.]
        Highs --> Indicators
        Lows --> Indicators
        Volumes --> Indicators
        MVRV --> Indicators
        FlowIn --> Indicators
        FlowOut --> Indicators
        Indicators --> Score[Calculate Weighted Score]
        Score --> Amount[Calculate Investment Amount]
        Amount --> Signals[Signal Array]
    end

    subgraph "3. Execution"
        Signals --> Filter[Filter by Interval<br/>Skip Warmup Period]
        Filter --> Trades[Execute Trades<br/>Calculate Units]
    end

    subgraph "4. Results"
        Trades --> Metrics[Calculate Metrics<br/>Returns, Cost Basis, etc.]
        Metrics --> Report[Generate Report]
    end

    style Indicators stroke:#ffaa00,stroke-width:2px
    style Score stroke:#ffaa00,stroke-width:2px
    style Amount stroke:#ffaa00,stroke-width:2px
```

## Indicator Calculation Flow

```mermaid
graph TB
    subgraph "Input Data"
        Prices[Price Array]
        Highs[High Array]
        Lows[Low Array]
        Volumes[Volume Array]
        MVRV[MVRV Array]
        FlowIn[Flow In Array]
        FlowOut[Flow Out Array]
    end

    subgraph "Indicator Calculations"
        Prices --> RSI[RSI<br/>14 period]
        Prices --> FG[Fear & Greed<br/>Placeholder]
        MVRV --> MVRVCalc[MVRV Score<br/>0-100 mapping]
        Prices --> MA200[MA200 Distance<br/>% from 200-day MA]
        Prices --> BB[Bollinger %B<br/>Position in bands]
        Prices --> ATR1[ATR Calculation]
        Highs --> ATR1
        Lows --> ATR1
        ATR1 --> ATRScore[ATR Volatility Score<br/>Low vol = buy signal]
        FlowIn --> NetFlow[Net Flow Score<br/>Outflows = bullish]
        FlowOut --> NetFlow
        Prices --> VolFlow1[Volume*Price]
        Volumes --> VolFlow1
        VolFlow1 --> VolFlow[Volume Flow Score<br/>Z-score of flow]
    end

    subgraph "Weighted Combination"
        RSI --> |25%| ScoreCalc[Composite Score<br/>Weighted Average]
        FG --> |15%| ScoreCalc
        MVRVCalc --> |15%| ScoreCalc
        MA200 --> |10%| ScoreCalc
        BB --> |10%| ScoreCalc
        ATRScore --> |10%| ScoreCalc
        NetFlow --> |10%| ScoreCalc
        VolFlow --> |5%| ScoreCalc
    end

    ScoreCalc --> FinalScore[Final Score 0-100<br/>Low = Oversold = Buy More<br/>High = Overbought = Buy Less]

    style ScoreCalc stroke:#00ff00,stroke-width:2px
    style FinalScore stroke:#ffaa00,stroke-width:3px
```

## Investment Amount Calculation

```mermaid
graph TB
    Score[Composite Score] --> Check{Score Value?}

    Check -->|≤ 30<br/>Oversold| MaxMult["BaseAmount * 3.0<br/>Max Multiplier"]
    Check -->|≥ 70<br/>Overbought| MinMult["BaseAmount * 0.5<br/>Min Multiplier"]
    Check -->|30-70<br/>Neutral| Linear[Linear Interpolation<br/>Between 3.0x and 0.5x]

    MaxMult --> Investment[Investment Amount]
    MinMult --> Investment
    Linear --> Investment

    Investment --> Example1["Example: Score=20<br/>$100 * 3.0 = $300"]
    Investment --> Example2["Example: Score=50<br/>$100 * 1.75 = $175"]
    Investment --> Example3["Example: Score=80<br/>$100 * 0.5 = $50"]

    style MaxMult stroke:#00ff00,stroke-width:2px
    style MinMult stroke:#ff4444,stroke-width:2px
    style Linear stroke:#ffaa00,stroke-width:2px
```

## Trade Execution Flow

```mermaid
sequenceDiagram
    participant Main as main.go
    participant Signals as signals.GenerateSignals()
    participant Exec as execution.Execute()
    participant Calc as execution.CalculateResult()

    Main->>Signals: MarketData + Config

    loop For Each Price Point
        Signals->>Signals: Calculate Indicators
        Signals->>Signals: Calculate Score
        Signals->>Signals: Calculate Amount
        Signals->>Signals: Create Signal
    end

    Signals-->>Main: []Signal

    Main->>Exec: []Signal + Config

    loop For Each Signal
        Exec->>Exec: Check Interval (every 7 days)
        Exec->>Exec: Check Warmup (skip first 14)
        alt Should Execute
            Exec->>Exec: Calculate Units = Amount/Price
            Exec->>Exec: Update Totals
            Exec->>Exec: Create Trade Record
        end
    end

    Exec-->>Main: []Trade

    Main->>Calc: []Trade + FinalPrice
    Calc->>Calc: Calculate Total Returns
    Calc->>Calc: Calculate Avg Cost Basis
    Calc-->>Main: *Result

    Main->>Main: Print Summary
    Main->>Main: Compare with Standard DCA
```

## Package Dependencies

```mermaid
graph TB
    subgraph "Entry Point"
        Main["cmd/backtester/main.go"]
    end

    subgraph "Core Pipeline"
        Data["data/loader.go"]
        Signals["pipeline/signals/"]
        Execution["pipeline/execution/"]
        Metrics["pipeline/metrics/"]
    end

    subgraph "Supporting"
        Indicators["pipeline/indicators/"]
        Util["util/"]
        Report["report/"]
        Config["config/"]
    end

    Main --> Data
    Main --> Signals
    Main --> Execution
    Main --> Report
    Main --> Util

    Data --> Util
    Signals --> Indicators
    Signals --> Util
    Execution --> Signals
    Execution --> Metrics
    Report --> Metrics

    Config -.-> Main

    style Main stroke:#ffaa00,stroke-width:3px
    style Signals stroke:#00bfff,stroke-width:2px
    style Execution stroke:#00bfff,stroke-width:2px
    style Metrics stroke:#00bfff,stroke-width:2px
```

## Key Data Structures

```mermaid
classDiagram
    class MarketData {
        +[]float64 Prices
        +[]float64 Highs
        +[]float64 Lows
        +[]float64 Volumes
        +[]float64 MVRVData
        +[]float64 FlowInUSD
        +[]float64 FlowOutUSD
        +time.Time FirstPriceDate
    }

    class IndicatorSet {
        +float64 RSI
        +float64 FearAndGreed
        +float64 MVRV
        +float64 MA200
        +float64 Bollinger
        +float64 ATR
        +float64 NetFlow
        +float64 VolumeFlow
    }

    class Signal {
        +int Index
        +float64 Price
        +float64 Score
        +float64 Amount
        +IndicatorSet
    }

    class Trade {
        +int Index
        +float64 Price
        +float64 Amount
        +float64 Units
        +float64 TotalUnits
        +float64 TotalInvest
        +float64 Score
        +IndicatorSet
    }

    class Result {
        +[]Trade Trades
        +float64 TotalInvested
        +float64 TotalUnits
        +float64 FinalPrice
        +float64 FinalValue
        +float64 TotalReturn
        +float64 AverageCostBasis
        +int NumTrades
    }

    class Config {
        +int ExecutionInterval
        +RSIConfig RSI
        +ScoreConfig Score
        +WeightsConfig Weights
        +BollingerConfig Bollinger
        +ATRConfig ATR
        +NetFlowConfig NetFlow
        +VolumeFlowConfig VolumeFlow
    }

    MarketData --> Signal : generates
    Signal *-- IndicatorSet : embeds
    Signal --> Trade : executed as
    Trade *-- IndicatorSet : includes
    Trade --> Result : aggregated into
    Config --> Signal : configures
```

## Execution Timeline Example

```mermaid
gantt
    title Trade Execution Timeline (Weekly DCA, 14-day warmup)
    dateFormat X
    axisFormat Day %d

    section Warmup
    Skip trades (warmup period)     :done, warmup, 0, 14

    section Trades
    Trade 1 (Day 14)                :crit, t1, 14, 1
    Skip (interval)                 :done, skip1, 15, 6
    Trade 2 (Day 21)                :crit, t2, 21, 1
    Skip (interval)                 :done, skip2, 22, 6
    Trade 3 (Day 28)                :crit, t3, 28, 1
    Skip (interval)                 :done, skip3, 29, 6
    Trade 4 (Day 35)                :crit, t4, 35, 1
```

## Constants & Configuration

```mermaid
mindmap
    root((Constants))
        Signals Package
            NeutralIndicatorValue = 50.0
        Util Package
            DefaultIndicatorValue = 50.0
            DefaultVolumeValue = 0.0
        Config JSON
            ExecutionInterval = 7 days
            RSIPeriod = 14 days
            BaseAmount = $100
            MaxMultiplier = 3.0x
            MinMultiplier = 0.5x
            ScoreOversold = 30
            ScoreOverbought = 70
```

## Error Handling & Validation

```mermaid
graph TB
    Start[Load Data] --> CheckPrices{Prices Loaded?}
    CheckPrices -->|No| Error1[Fatal Error]
    CheckPrices -->|Yes| CheckMVRV{MVRV Loaded?}

    CheckMVRV -->|No| Warn1[Warning: Using Neutral Values]
    CheckMVRV -->|Yes| ValidateDate{Date Matches?}

    ValidateDate -->|No| Warn2[Warning: Date Mismatch<br/>Use Neutral Values]
    ValidateDate -->|Yes| UseMVRV[Use MVRV Data]

    Warn1 --> Continue[Continue Backtest]
    Warn2 --> Continue
    UseMVRV --> Continue

    Continue --> GenSignals[Generate Signals]
    GenSignals --> CheckSignals{Signals Valid?}
    CheckSignals -->|No| Error2[Return nil]
    CheckSignals -->|Yes| Execute[Execute Trades]

    style Error1 stroke:#ff0000,stroke-width:3px
    style Error2 stroke:#ff0000,stroke-width:3px
    style Warn1 stroke:#ff8800,stroke-width:2px
    style Warn2 stroke:#ff8800,stroke-width:2px
```

---

## Quick Reference

### Main Pipeline Functions

1. **`data.LoadPricesFromCSV()`** - Load price data from CSV
2. **`signals.GenerateSignals()`** - Calculate indicators and create signals
3. **`execution.Execute()`** - Execute trades based on signals
4. **`execution.CalculateResult()`** - Compute final metrics
5. **`report.PrintSummary()`** - Display results

### Key Configuration Parameters

- **ExecutionInterval**: Days between DCA purchases (default: 7)
- **RSIPeriod**: RSI calculation window (default: 14) - also sets warmup period
- **BaseAmount**: Base investment per interval (default: $100)
- **Weights**: Indicator importance (must sum to 100%)

### Indicator Score Meaning

- **0-30**: Oversold → Buy MORE (3.0x multiplier)
- **30-70**: Neutral → Linear scaling (3.0x → 0.5x)
- **70-100**: Overbought → Buy LESS (0.5x multiplier)
