package types

import (
	"time"

	"github.com/shopspring/decimal"
)

// OrderBook represents the order book for an asset pair.
type OrderBook struct {
	// Asks is the list of ask orders [price, volume, timestamp].
	Asks [][]interface{} `json:"asks"`

	// Bids is the list of bid orders [price, volume, timestamp].
	Bids [][]interface{} `json:"bids"`
}

// BookLevel represents a price level in the order book (WebSocket).
type BookLevel struct {
	// Price is the price at this level.
	Price decimal.Decimal `json:"price"`

	// Qty is the total quantity at this level.
	Qty decimal.Decimal `json:"qty"`
}

// BookData represents order book data from WebSocket.
type BookData struct {
	// Symbol is the trading pair.
	Symbol string `json:"symbol"`

	// Bids is the list of bid levels.
	Bids []BookLevel `json:"bids"`

	// Asks is the list of ask levels.
	Asks []BookLevel `json:"asks"`

	// Checksum is the CRC32 checksum for top 10 levels.
	Checksum uint32 `json:"checksum"`

	// Timestamp is the book update timestamp.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// BookSnapshot represents the initial order book snapshot.
type BookSnapshot struct {
	BookData

	// Type indicates this is a snapshot.
	Type string `json:"type"` // "snapshot"
}

// BookUpdate represents an incremental order book update.
type BookUpdate struct {
	BookData

	// Type indicates this is an update.
	Type string `json:"type"` // "update"
}

// Level3Order represents a single order in the Level 3 book.
type Level3Order struct {
	// OrderID is the unique order identifier.
	OrderID string `json:"order_id"`

	// LimitPrice is the order's limit price.
	LimitPrice decimal.Decimal `json:"limit_price"`

	// OrderQty is the order quantity.
	OrderQty decimal.Decimal `json:"order_qty"`

	// Timestamp is when the order was placed.
	Timestamp time.Time `json:"timestamp"`
}

// Level3Data represents Level 3 order book data from WebSocket.
type Level3Data struct {
	// Symbol is the trading pair.
	Symbol string `json:"symbol"`

	// Bids is the list of bid orders.
	Bids []Level3Order `json:"bids,omitempty"`

	// Asks is the list of ask orders.
	Asks []Level3Order `json:"asks,omitempty"`

	// Checksum is the CRC32 checksum.
	Checksum uint32 `json:"checksum,omitempty"`
}

// L3BookOrder represents a single order in the REST Level 3 book response.
type L3BookOrder struct {
	// Price is the order price.
	Price string `json:"price"`

	// Qty is the order quantity.
	Qty string `json:"qty"`

	// OrderID is the unique order identifier.
	OrderID string `json:"order_id"`

	// Timestamp is the order timestamp in nanoseconds.
	Timestamp int64 `json:"timestamp"`
}

// L3OrderBook represents the REST Level 3 order book response.
type L3OrderBook struct {
	// Pair is the asset pair.
	Pair string `json:"pair"`

	// Bids is the list of bid orders.
	Bids []L3BookOrder `json:"bids"`

	// Asks is the list of ask orders.
	Asks []L3BookOrder `json:"asks"`
}

// GroupedBookLevel represents a price level in the grouped order book.
type GroupedBookLevel struct {
	// Price is the grouped price level.
	Price string `json:"price"`

	// Qty is the aggregated quantity at this level.
	Qty string `json:"qty"`
}

// GroupedOrderBook represents the grouped order book response.
type GroupedOrderBook struct {
	// Pair is the asset pair.
	Pair string `json:"pair"`

	// Grouping is the grouping value used.
	Grouping int `json:"grouping"`

	// Bids is the list of aggregated bid levels.
	Bids []GroupedBookLevel `json:"bids"`

	// Asks is the list of aggregated ask levels.
	Asks []GroupedBookLevel `json:"asks"`
}

// PreTradeLevel represents a price level in the pre-trade transparency data.
type PreTradeLevel struct {
	// Side indicates BUY or SELL.
	Side string `json:"side"`

	// Price is the price level.
	Price string `json:"price"`

	// Qty is the aggregated quantity at this level.
	Qty string `json:"qty"`

	// Count is the number of orders at this level.
	Count int `json:"count"`

	// PublicationTS is the timestamp when this level was published.
	PublicationTS string `json:"publication_ts"`
}

// PreTradeData represents pre-trade transparency data.
type PreTradeData struct {
	// Symbol is the trading pair symbol.
	Symbol string `json:"symbol"`

	// Description is the full description of the pair.
	Description string `json:"description"`

	// BaseAsset is the base currency code.
	BaseAsset string `json:"base_asset"`

	// BaseNotation indicates the notation type (NOML).
	BaseNotation string `json:"base_notation"`

	// QuoteAsset is the quote currency code.
	QuoteAsset string `json:"quote_asset"`

	// QuoteNotation indicates the notation type (MONE).
	QuoteNotation string `json:"quote_notation"`

	// Venue is the MIC of the trading platform.
	Venue string `json:"venue"`

	// System indicates CLOB (Central Limit Order Book).
	System string `json:"system"`

	// Bids is the list of bid levels.
	Bids []PreTradeLevel `json:"bids"`

	// Asks is the list of ask levels.
	Asks []PreTradeLevel `json:"asks"`
}

// PostTrade represents a single trade in the post-trade transparency data.
type PostTrade struct {
	// TradeID is the Kraken unique trade identifier.
	TradeID string `json:"trade_id"`

	// Price is the trade price excluding fees.
	Price string `json:"price"`

	// Quantity is the trade quantity.
	Quantity string `json:"quantity"`

	// Symbol is the trading pair symbol.
	Symbol string `json:"symbol"`

	// Description is the full description of the pair.
	Description string `json:"description"`

	// BaseAsset is the base currency code.
	BaseAsset string `json:"base_asset"`

	// BaseNotation indicates the notation type (UNIT).
	BaseNotation string `json:"base_notation"`

	// QuoteAsset is the quote currency code.
	QuoteAsset string `json:"quote_asset"`

	// QuoteNotation indicates the notation type (MONE).
	QuoteNotation string `json:"quote_notation"`

	// TradeVenue is the MIC where the trade was executed.
	TradeVenue string `json:"trade_venue"`

	// TradeTS is the trade timestamp.
	TradeTS string `json:"trade_ts"`

	// PublicationVenue is the MIC where the trade was published.
	PublicationVenue string `json:"publication_venue"`

	// PublicationTS is the publication timestamp.
	PublicationTS string `json:"publication_ts"`
}

// PostTradeResult represents the post-trade transparency data response.
type PostTradeResult struct {
	// LastTS is the timestamp of the latest trade.
	LastTS string `json:"last_ts"`

	// Count is the number of trades returned.
	Count int `json:"count"`

	// Trades is the list of trades.
	Trades []PostTrade `json:"trades"`
}
