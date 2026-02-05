package types

import "github.com/shopspring/decimal"

// Ticker represents ticker information for an asset pair.
type Ticker struct {
	// Ask is the ask price [price, whole_lot_volume, lot_volume].
	Ask []decimal.Decimal `json:"a"`

	// Bid is the bid price [price, whole_lot_volume, lot_volume].
	Bid []decimal.Decimal `json:"b"`

	// Close is the last trade closed [price, lot_volume].
	Close []decimal.Decimal `json:"c"`

	// Volume is the volume [today, last_24h].
	Volume []decimal.Decimal `json:"v"`

	// VWAP is the volume weighted average price [today, last_24h].
	VWAP []decimal.Decimal `json:"p"`

	// Trades is the number of trades [today, last_24h].
	Trades []int `json:"t"`

	// Low is the low price [today, last_24h].
	Low []decimal.Decimal `json:"l"`

	// High is the high price [today, last_24h].
	High []decimal.Decimal `json:"h"`

	// OpeningPrice is today's opening price.
	OpeningPrice decimal.Decimal `json:"o"`
}

// TickerData represents real-time ticker data from WebSocket.
type TickerData struct {
	// Symbol is the trading pair.
	Symbol string `json:"symbol"`

	// Bid is the best bid price.
	Bid decimal.Decimal `json:"bid"`

	// BidQty is the bid quantity.
	BidQty decimal.Decimal `json:"bid_qty"`

	// Ask is the best ask price.
	Ask decimal.Decimal `json:"ask"`

	// AskQty is the ask quantity.
	AskQty decimal.Decimal `json:"ask_qty"`

	// Last is the last trade price.
	Last decimal.Decimal `json:"last"`

	// Volume is the 24h volume.
	Volume decimal.Decimal `json:"volume"`

	// VWAP is the 24h volume weighted average price.
	VWAP decimal.Decimal `json:"vwap"`

	// Low is the 24h low price.
	Low decimal.Decimal `json:"low"`

	// High is the 24h high price.
	High decimal.Decimal `json:"high"`

	// Change is the 24h price change.
	Change decimal.Decimal `json:"change"`

	// ChangePct is the 24h price change percentage.
	ChangePct decimal.Decimal `json:"change_pct"`
}
