package websocket

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// MessageType represents the type of WebSocket message.
type MessageType string

const (
	// Request message types
	MsgTypeSubscribe   MessageType = "subscribe"
	MsgTypeUnsubscribe MessageType = "unsubscribe"
	MsgTypePing        MessageType = "ping"
	MsgTypePong        MessageType = "pong"

	// Trading request types
	MsgTypeAddOrder           MessageType = "add_order"
	MsgTypeEditOrder          MessageType = "edit_order"
	MsgTypeAmendOrder         MessageType = "amend_order"
	MsgTypeCancelOrder        MessageType = "cancel_order"
	MsgTypeCancelAll          MessageType = "cancel_all"
	MsgTypeBatchAdd           MessageType = "batch_add"
	MsgTypeBatchCancel        MessageType = "batch_cancel"
	MsgTypeCancelOnDisconnect MessageType = "cancel_all_orders_after"
)

// Channel represents a WebSocket channel name.
type Channel string

const (
	// Public channels
	ChannelTicker     Channel = "ticker"
	ChannelBook       Channel = "book"
	ChannelTrade      Channel = "trade"
	ChannelOHLC       Channel = "ohlc"
	ChannelInstrument Channel = "instrument"
	ChannelStatus     Channel = "status"
	ChannelHeartbeat  Channel = "heartbeat"

	// Private channels
	ChannelExecutions Channel = "executions"
	ChannelBalances   Channel = "balances"
)

// BaseMessage contains common fields for all messages.
type BaseMessage struct {
	// Channel is the channel name.
	Channel Channel `json:"channel,omitempty"`
	// Type is the message type.
	Type string `json:"type,omitempty"`
	// Data contains the message payload.
	Data json.RawMessage `json:"data,omitempty"`
	// ReqID is the request ID for request/response correlation.
	ReqID int64 `json:"req_id,omitempty"`
	// Method is the method name for trading requests.
	Method MessageType `json:"method,omitempty"`
	// Params contains request parameters.
	Params json.RawMessage `json:"params,omitempty"`
	// Result contains response data.
	Result json.RawMessage `json:"result,omitempty"`
	// Success indicates if the request was successful.
	Success bool `json:"success,omitempty"`
	// Error contains error information.
	Error string `json:"error,omitempty"`
	// TimeIn is when the message was received by Kraken.
	TimeIn string `json:"time_in,omitempty"`
	// TimeOut is when the message was sent by Kraken.
	TimeOut string `json:"time_out,omitempty"`
}

// SubscribeRequest represents a subscription request.
type SubscribeRequest struct {
	Method MessageType     `json:"method"`
	ReqID  int64           `json:"req_id,omitempty"`
	Params SubscribeParams `json:"params"`
}

// SubscribeParams contains subscription parameters.
type SubscribeParams struct {
	Channel  Channel  `json:"channel"`
	Symbol   []string `json:"symbol,omitempty"`
	Depth    int      `json:"depth,omitempty"`
	Interval int      `json:"interval,omitempty"`
	Snapshot bool     `json:"snapshot,omitempty"`
	Token    string   `json:"token,omitempty"`
}

// UnsubscribeRequest represents an unsubscription request.
type UnsubscribeRequest struct {
	Method MessageType       `json:"method"`
	ReqID  int64             `json:"req_id,omitempty"`
	Params UnsubscribeParams `json:"params"`
}

// UnsubscribeParams contains unsubscription parameters.
type UnsubscribeParams struct {
	Channel  Channel  `json:"channel"`
	Symbol   []string `json:"symbol,omitempty"`
	Depth    int      `json:"depth,omitempty"`
	Interval int      `json:"interval,omitempty"`
	Token    string   `json:"token,omitempty"`
}

// PingMessage represents a ping request.
type PingMessage struct {
	Method MessageType `json:"method"`
	ReqID  int64       `json:"req_id,omitempty"`
}

// TickerData represents ticker data from WebSocket.
type TickerData struct {
	Symbol    string          `json:"symbol"`
	Bid       decimal.Decimal `json:"bid"`
	BidQty    decimal.Decimal `json:"bid_qty"`
	Ask       decimal.Decimal `json:"ask"`
	AskQty    decimal.Decimal `json:"ask_qty"`
	Last      decimal.Decimal `json:"last"`
	Volume    decimal.Decimal `json:"volume"`
	VolumeWAP decimal.Decimal `json:"vwap"`
	Low       decimal.Decimal `json:"low"`
	High      decimal.Decimal `json:"high"`
	Change    decimal.Decimal `json:"change"`
	ChangePct decimal.Decimal `json:"change_pct"`
}

// BookEntry represents a single order book entry.
type BookEntry struct {
	Price decimal.Decimal `json:"price"`
	Qty   decimal.Decimal `json:"qty"`
}

// BookData represents order book data from WebSocket.
type BookData struct {
	Symbol    string      `json:"symbol"`
	Checksum  uint32      `json:"checksum"`
	Bids      []BookEntry `json:"bids"`
	Asks      []BookEntry `json:"asks"`
	Timestamp time.Time   `json:"timestamp"`
}

// TradeData represents trade data from WebSocket.
type TradeData struct {
	Symbol    string          `json:"symbol"`
	Side      string          `json:"side"`
	Price     decimal.Decimal `json:"price"`
	Qty       decimal.Decimal `json:"qty"`
	OrderType string          `json:"ord_type"`
	TradeID   int64           `json:"trade_id"`
	Timestamp time.Time       `json:"timestamp"`
}

// OHLCData represents OHLC data from WebSocket.
type OHLCData struct {
	Symbol    string          `json:"symbol"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	VolumeWAP decimal.Decimal `json:"vwap"`
	Volume    decimal.Decimal `json:"volume"`
	Count     int             `json:"count"`
	Interval  int             `json:"interval"`
	Timestamp time.Time       `json:"timestamp"`
	EndTime   time.Time       `json:"interval_begin"`
}

// InstrumentData represents instrument data from WebSocket.
type InstrumentData struct {
	Symbol         string          `json:"symbol"`
	Status         string          `json:"status"`
	StatusNote     string          `json:"status_note,omitempty"`
	PricePrecision int             `json:"price_precision"`
	QtyPrecision   int             `json:"qty_precision"`
	PriceIncrement decimal.Decimal `json:"price_increment"`
	QtyIncrement   decimal.Decimal `json:"qty_increment"`
	QtyMin         decimal.Decimal `json:"qty_min"`
	Marginable     bool            `json:"marginable"`
	HasIndex       bool            `json:"has_index"`
	Base           string          `json:"base"`
	Quote          string          `json:"quote"`
}

// StatusData represents system status data from WebSocket.
type StatusData struct {
	System       string `json:"system"`
	Version      string `json:"version"`
	APIVersion   string `json:"api_version"`
	ConnectionID int64  `json:"connection_id"`
}

// HeartbeatData represents heartbeat data from WebSocket.
type HeartbeatData struct {
	Timestamp time.Time `json:"timestamp"`
}

// ExecutionData represents execution data from WebSocket.
type ExecutionData struct {
	OrderID     string           `json:"order_id"`
	ExecID      string           `json:"exec_id"`
	TradeID     int64            `json:"trade_id,omitempty"`
	Symbol      string           `json:"symbol"`
	Side        string           `json:"side"`
	OrderType   string           `json:"order_type"`
	OrderQty    decimal.Decimal  `json:"order_qty"`
	LastQty     decimal.Decimal  `json:"last_qty,omitempty"`
	LastPrice   decimal.Decimal  `json:"last_price,omitempty"`
	CumQty      decimal.Decimal  `json:"cum_qty"`
	CumCost     decimal.Decimal  `json:"cum_cost"`
	AvgPrice    decimal.Decimal  `json:"avg_price"`
	OrderStatus string           `json:"order_status"`
	ExecType    string           `json:"exec_type"`
	Timestamp   time.Time        `json:"timestamp"`
	LimitPrice  *decimal.Decimal `json:"limit_price,omitempty"`
}

// BalanceData represents balance data from WebSocket.
type BalanceData struct {
	Asset      string          `json:"asset"`
	Balance    decimal.Decimal `json:"balance"`
	WalletID   string          `json:"wallet_id,omitempty"`
	WalletType string          `json:"wallet_type,omitempty"`
}

// AddOrderParams contains parameters for adding an order via WebSocket.
type AddOrderParams struct {
	// Required fields
	OrderType string `json:"order_type"`
	Side      string `json:"side"`
	Symbol    string `json:"symbol"`
	OrderQty  string `json:"order_qty"`
	Token     string `json:"token"`

	// Optional fields
	LimitPrice     string `json:"limit_price,omitempty"`
	LimitPriceType string `json:"limit_price_type,omitempty"`
	TimeInForce    string `json:"time_in_force,omitempty"`
	Margin         bool   `json:"margin,omitempty"`
	PostOnly       bool   `json:"post_only,omitempty"`
	ReduceOnly     bool   `json:"reduce_only,omitempty"`
	EffectiveTime  string `json:"effective_time,omitempty"`
	ExpireTime     string `json:"expire_time,omitempty"`
	Deadline       string `json:"deadline,omitempty"`
	ClientOrderID  string `json:"cl_ord_id,omitempty"`
	OrderUserRef   int64  `json:"order_userref,omitempty"`

	// Iceberg orders
	DisplayQty string `json:"display_qty,omitempty"`

	// Misc options
	FeePreference string `json:"fee_preference,omitempty"`
	NoMPP         bool   `json:"no_mpp,omitempty"`
	STPType       string `json:"stp_type,omitempty"`

	// Buy market orders (cash quantity in quote) without margin funding
	CashOrderQty string `json:"cash_order_qty,omitempty"`

	// Institutional accounts with enhanced STP
	SenderSubID string `json:"sender_sub_id,omitempty"`

	// Advanced order configs
	Triggers    *TriggerParams     `json:"triggers,omitempty"`
	Conditional *ConditionalParams `json:"conditional,omitempty"`

	// Validation only
	Validate bool `json:"validate,omitempty"`
}

// TriggerParams contains trigger parameters for stop/trailing orders.
type TriggerParams struct {
	Reference string `json:"reference,omitempty"`
	Price     string `json:"price"`
	PriceType string `json:"price_type,omitempty"`
}

// ConditionalParams contains conditional close parameters.
type ConditionalParams struct {
	OrderType        string `json:"order_type"`
	LimitPrice       string `json:"limit_price,omitempty"`
	LimitPriceType   string `json:"limit_price_type,omitempty"`
	TriggerPrice     string `json:"trigger_price,omitempty"`
	TriggerPriceType string `json:"trigger_price_type,omitempty"`
}

// EditOrderParams contains parameters for editing an order via WebSocket.
type EditOrderParams struct {
	OrderID    string `json:"order_id"`
	Symbol     string `json:"symbol"`
	OrderQty   string `json:"order_qty,omitempty"`
	LimitPrice string `json:"limit_price,omitempty"`
	Token      string `json:"token"`
}

// AmendOrderParams contains parameters for amending an order via WebSocket.
type AmendOrderParams struct {
	OrderID      string `json:"order_id"`
	OrderQty     string `json:"order_qty,omitempty"`
	LimitPrice   string `json:"limit_price,omitempty"`
	TriggerPrice string `json:"trigger_price,omitempty"`
	PostOnly     *bool  `json:"post_only,omitempty"`
	Token        string `json:"token"`
}

// CancelOrderParams contains parameters for canceling an order via WebSocket.
type CancelOrderParams struct {
	OrderID       []string `json:"order_id,omitempty"`
	ClientOrderID []string `json:"cl_ord_id,omitempty"`
	Token         string   `json:"token"`
}

// CancelAllParams contains parameters for canceling all orders via WebSocket.
type CancelAllParams struct {
	Token string `json:"token"`
}

// CancelOnDisconnectParams contains parameters for dead man's switch via WebSocket.
type CancelOnDisconnectParams struct {
	Timeout int    `json:"timeout"`
	Token   string `json:"token"`
}

// BatchAddParams contains parameters for batch adding orders via WebSocket.
type BatchAddParams struct {
	Symbol   string           `json:"symbol"`
	Orders   []AddOrderParams `json:"orders"`
	Token    string           `json:"token"`
	Deadline string           `json:"deadline,omitempty"`
	Validate bool             `json:"validate,omitempty"`
}

// BatchCancelParams contains parameters for batch canceling orders via WebSocket.
type BatchCancelParams struct {
	Orders []CancelOrderParams `json:"orders"`
	Token  string              `json:"token"`
}

// AddOrderResult contains the result of an add order request.
type AddOrderResult struct {
	OrderID       string   `json:"order_id"`
	ClientOrderID string   `json:"cl_ord_id,omitempty"`
	OrderUserRef  int64    `json:"order_userref,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
}

// EditOrderResult contains the result of an edit order request.
type EditOrderResult struct {
	OrderID         string `json:"order_id"`
	OriginalOrderID string `json:"original_order_id"`
}

// AmendOrderResult contains the result of an amend order request.
type AmendOrderResult struct {
	AmendID string `json:"amend_id"`
	OrderID string `json:"order_id"`
}

// CancelOrderResult contains the result of a cancel order request.
type CancelOrderResult struct {
	OrderID string `json:"order_id"`
}

// CancelAllResult contains the result of a cancel all request.
type CancelAllResult struct {
	Count int `json:"count"`
}

// CancelOnDisconnectResult contains the result of setting dead man's switch.
type CancelOnDisconnectResult struct {
	CurrentTime string `json:"currentTime"`
	TriggerTime string `json:"triggerTime"`
}
