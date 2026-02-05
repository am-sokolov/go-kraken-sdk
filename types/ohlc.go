package types

import (
	"time"

	"github.com/shopspring/decimal"
)

// OHLCInterval represents the time interval for OHLC data.
type OHLCInterval int

const (
	Interval1Min   OHLCInterval = 1
	Interval5Min   OHLCInterval = 5
	Interval15Min  OHLCInterval = 15
	Interval30Min  OHLCInterval = 30
	Interval1Hour  OHLCInterval = 60
	Interval4Hour  OHLCInterval = 240
	Interval1Day   OHLCInterval = 1440
	Interval1Week  OHLCInterval = 10080
	Interval15Days OHLCInterval = 21600
)

// OHLC represents a single OHLC candle.
type OHLC struct {
	// Time is the candle start time.
	Time int64 `json:"time"`

	// Open is the opening price.
	Open decimal.Decimal `json:"open"`

	// High is the highest price.
	High decimal.Decimal `json:"high"`

	// Low is the lowest price.
	Low decimal.Decimal `json:"low"`

	// Close is the closing price.
	Close decimal.Decimal `json:"close"`

	// VWAP is the volume weighted average price.
	VWAP decimal.Decimal `json:"vwap"`

	// Volume is the candle volume.
	Volume decimal.Decimal `json:"volume"`

	// Count is the number of trades.
	Count int `json:"count"`
}

// OHLCResult contains the result of an OHLC query.
type OHLCResult struct {
	// Data is the OHLC data keyed by pair.
	Data map[string][]OHLC

	// Last is the ID for pagination.
	Last int64
}

// OHLCData represents real-time OHLC data from WebSocket.
type OHLCData struct {
	// Symbol is the trading pair.
	Symbol string `json:"symbol"`

	// Open is the opening price.
	Open decimal.Decimal `json:"open"`

	// High is the highest price.
	High decimal.Decimal `json:"high"`

	// Low is the lowest price.
	Low decimal.Decimal `json:"low"`

	// Close is the closing price.
	Close decimal.Decimal `json:"close"`

	// VWAP is the volume weighted average price.
	VWAP decimal.Decimal `json:"vwap"`

	// Volume is the candle volume.
	Volume decimal.Decimal `json:"volume"`

	// Trades is the number of trades.
	Trades int `json:"trades"`

	// IntervalBegin is the interval start time.
	IntervalBegin time.Time `json:"interval_begin"`

	// Interval is the interval in minutes.
	Interval int `json:"interval"`

	// Timestamp is the data timestamp.
	Timestamp time.Time `json:"timestamp"`
}

// Spread represents spread data from the spreads endpoint.
type Spread struct {
	// Time is the spread timestamp.
	Time int64

	// Bid is the bid price.
	Bid decimal.Decimal

	// Ask is the ask price.
	Ask decimal.Decimal
}

// SpreadResult contains the result of a spread query.
type SpreadResult struct {
	// Data is the spread data keyed by pair.
	Data map[string][]Spread

	// Last is the ID for pagination.
	Last int64
}
