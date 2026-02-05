package types

import (
	"time"

	"github.com/shopspring/decimal"
)

// Trade represents a trade from the trades endpoint.
type Trade struct {
	// OrderTxID is the order transaction ID.
	OrderTxID string `json:"ordertxid"`

	// PosTxID is the position transaction ID.
	PosTxID string `json:"postxid"`

	// Pair is the asset pair.
	Pair string `json:"pair"`

	// Time is the trade timestamp.
	Time float64 `json:"time"`

	// Type is the trade direction (buy/sell).
	Type string `json:"type"`

	// OrderType is the order type.
	OrderType string `json:"ordertype"`

	// Price is the trade price.
	Price decimal.Decimal `json:"price"`

	// Cost is the trade cost.
	Cost decimal.Decimal `json:"cost"`

	// Fee is the trade fee.
	Fee decimal.Decimal `json:"fee"`

	// Vol is the trade volume.
	Vol decimal.Decimal `json:"vol"`

	// Margin is the margin used.
	Margin decimal.Decimal `json:"margin"`

	// Leverage is the leverage used.
	Leverage string `json:"leverage,omitempty"`

	// Misc contains miscellaneous info.
	Misc string `json:"misc"`

	// Trade ID as provided by trades endpoint.
	TradeID int64 `json:"trade_id,omitempty"`

	// Maker indicates if this was a maker trade.
	Maker bool `json:"maker,omitempty"`
}

// PublicTrade represents a trade from the public trades endpoint.
type PublicTrade struct {
	// Price is the trade price.
	Price decimal.Decimal

	// Volume is the trade volume.
	Volume decimal.Decimal

	// Time is the trade timestamp.
	Time float64

	// Side is the trade side (buy/sell).
	Side string

	// OrderType is the order type (market/limit).
	OrderType string

	// Misc contains miscellaneous info.
	Misc string

	// TradeID is the trade identifier.
	TradeID int64
}

// TradeData represents real-time trade data from WebSocket.
type TradeData struct {
	// Symbol is the trading pair.
	Symbol string `json:"symbol"`

	// Side is the trade side (buy/sell).
	Side string `json:"side"`

	// Price is the trade price.
	Price decimal.Decimal `json:"price"`

	// Qty is the trade quantity.
	Qty decimal.Decimal `json:"qty"`

	// OrderType is the order type that caused the trade.
	OrderType string `json:"ord_type"`

	// TradeID is the trade identifier.
	TradeID int64 `json:"trade_id"`

	// Timestamp is the trade timestamp.
	Timestamp time.Time `json:"timestamp"`
}

// TradesResult contains the result of a trades query.
type TradesResult struct {
	// Trades is a map of trade ID to trade.
	Trades map[string]Trade `json:"trades"`

	// Count is the total number of trades.
	Count int `json:"count"`
}

// PublicTradesResult contains the result of a public trades query.
type PublicTradesResult struct {
	// Trades is the list of trades.
	Trades []PublicTrade

	// Last is the ID for pagination.
	Last string
}
