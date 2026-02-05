package types

import "github.com/shopspring/decimal"

// Asset represents a tradable asset on Kraken.
type Asset struct {
	// AClass is the asset class.
	AClass string `json:"aclass"`

	// Altname is the alternate name.
	Altname string `json:"altname"`

	// Decimals is the number of decimal places for display.
	Decimals int `json:"decimals"`

	// DisplayDecimals is the number of decimals for display.
	DisplayDecimals int `json:"display_decimals"`

	// CollateralValue is the collateral value for margin.
	CollateralValue *decimal.Decimal `json:"collateral_value,omitempty"`

	// Status is the asset status.
	Status string `json:"status,omitempty"`
}

// AssetPair represents a tradable asset pair.
type AssetPair struct {
	// Altname is the alternate pair name.
	Altname string `json:"altname"`

	// WSName is the WebSocket pair name.
	WSName string `json:"wsname"`

	// AClassBase is the asset class of base component.
	AClassBase string `json:"aclass_base"`

	// Base is the asset ID of base component.
	Base string `json:"base"`

	// AClassQuote is the asset class of quote component.
	AClassQuote string `json:"aclass_quote"`

	// Quote is the asset ID of quote component.
	Quote string `json:"quote"`

	// PairDecimals is the number of decimal places for price.
	PairDecimals int `json:"pair_decimals"`

	// CostDecimals is the number of decimals for cost.
	CostDecimals int `json:"cost_decimals"`

	// LotDecimals is the number of decimals for volume.
	LotDecimals int `json:"lot_decimals"`

	// LotMultiplier is the amount to multiply lot volume.
	LotMultiplier int `json:"lot_multiplier"`

	// LeverageBuy is available leverage for buying.
	LeverageBuy []int `json:"leverage_buy"`

	// LeverageSell is available leverage for selling.
	LeverageSell []int `json:"leverage_sell"`

	// Fees is the fee schedule [volume, percent_fee].
	Fees [][]decimal.Decimal `json:"fees"`

	// FeesMaker is the maker fee schedule.
	FeesMaker [][]decimal.Decimal `json:"fees_maker"`

	// FeeVolumeCurrency is the volume discount currency.
	FeeVolumeCurrency string `json:"fee_volume_currency"`

	// MarginCall is the margin call level.
	MarginCall int `json:"margin_call"`

	// MarginStop is the stop-out/liquidation level.
	MarginStop int `json:"margin_stop"`

	// OrderMin is the minimum order size in base currency.
	OrderMin decimal.Decimal `json:"ordermin"`

	// CostMin is the minimum order cost in quote currency.
	CostMin decimal.Decimal `json:"costmin"`

	// TickSize is the minimum price increment.
	TickSize decimal.Decimal `json:"tick_size"`

	// Status is the pair status (online, cancel_only, etc.).
	Status string `json:"status"`

	// LongPositionLimit is the max long margin position.
	LongPositionLimit *int64 `json:"long_position_limit,omitempty"`

	// ShortPositionLimit is the max short margin position.
	ShortPositionLimit *int64 `json:"short_position_limit,omitempty"`
}

// Instrument represents reference data for an asset pair (WebSocket).
type Instrument struct {
	// Symbol is the trading pair symbol.
	Symbol string `json:"symbol"`

	// Base is the base asset.
	Base string `json:"base"`

	// Quote is the quote asset.
	Quote string `json:"quote"`

	// Status is the pair status.
	Status string `json:"status"`

	// HasIndex indicates if index price is available.
	HasIndex bool `json:"has_index"`

	// Marginable indicates if margin trading is available.
	Marginable bool `json:"marginable"`

	// MarginInitial is the initial margin requirement.
	MarginInitial *decimal.Decimal `json:"margin_initial,omitempty"`

	// PriceIncrement is the minimum price increment.
	PriceIncrement decimal.Decimal `json:"price_increment"`

	// PricePrecision is the price decimal places.
	PricePrecision int `json:"price_precision"`

	// QtyIncrement is the minimum quantity increment.
	QtyIncrement decimal.Decimal `json:"qty_increment"`

	// QtyMin is the minimum order quantity.
	QtyMin decimal.Decimal `json:"qty_min"`

	// QtyPrecision is the quantity decimal places.
	QtyPrecision int `json:"qty_precision"`

	// CostMin is the minimum order cost.
	CostMin *decimal.Decimal `json:"cost_min,omitempty"`

	// CostPrecision is the cost decimal places.
	CostPrecision int `json:"cost_precision,omitempty"`
}
