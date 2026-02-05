package types

import "github.com/shopspring/decimal"

// Balance represents account balance for an asset.
type Balance struct {
	// Balance is the total balance.
	Balance decimal.Decimal `json:"balance"`

	// Credit is the credit amount.
	Credit decimal.Decimal `json:"credit,omitempty"`

	// CreditUsed is the amount of credit used.
	CreditUsed decimal.Decimal `json:"credit_used,omitempty"`

	// HoldTrade is the amount held for trading.
	HoldTrade decimal.Decimal `json:"hold_trade,omitempty"`
}

// TradeBalance represents the trade balance summary.
type TradeBalance struct {
	// EquivalentBalance is the equivalent balance in query currency.
	EquivalentBalance decimal.Decimal `json:"eb"`

	// TradeBalance is the trade balance in query currency.
	TradeBalance decimal.Decimal `json:"tb"`

	// MarginOpen is the margin used for open positions.
	MarginOpen decimal.Decimal `json:"m"`

	// UnrealizedPnL is the unrealized P&L of open positions.
	UnrealizedPnL decimal.Decimal `json:"n"`

	// CostBasis is the cost basis of open positions.
	CostBasis decimal.Decimal `json:"c"`

	// FloatingValuation is the floating valuation.
	FloatingValuation decimal.Decimal `json:"v"`

	// Equity is the trade balance + unrealized P&L.
	Equity decimal.Decimal `json:"e"`

	// FreeMargin is the available margin.
	FreeMargin decimal.Decimal `json:"mf"`

	// MarginLevel is the margin level percentage.
	MarginLevel decimal.Decimal `json:"ml,omitempty"`

	// UnexecutedValue is the value of unexecuted orders.
	UnexecutedValue decimal.Decimal `json:"uv,omitempty"`
}

// BalanceData represents real-time balance data from WebSocket.
type BalanceData struct {
	// Asset is the asset name.
	Asset string `json:"asset"`

	// Balance is the total balance.
	Balance decimal.Decimal `json:"balance"`

	// WalletID is the wallet identifier.
	WalletID string `json:"wallet_id,omitempty"`

	// WalletType is the wallet type.
	WalletType string `json:"wallet_type,omitempty"`
}

// LedgerEntry represents a balance update from WebSocket executions.
type LedgerEntry struct {
	// Asset is the asset name.
	Asset string `json:"asset"`

	// Amount is the change amount.
	Amount decimal.Decimal `json:"amount"`

	// Balance is the new balance after change.
	Balance decimal.Decimal `json:"balance"`

	// Fee is any fee associated with this entry.
	Fee decimal.Decimal `json:"fee,omitempty"`

	// RefID is the reference ID.
	RefID string `json:"refid,omitempty"`

	// Type is the ledger entry type.
	Type string `json:"type"`
}

// Position represents an open margin position.
type Position struct {
	// OrderTxID is the order transaction ID.
	OrderTxID string `json:"ordertxid"`

	// PosTxID is the position transaction ID.
	PosTxID string `json:"postxid"`

	// Pair is the asset pair.
	Pair string `json:"pair"`

	// Time is when the position was opened.
	Time float64 `json:"time"`

	// Type is the position direction (buy/sell).
	Type string `json:"type"`

	// OrderType is the order type.
	OrderType string `json:"ordertype"`

	// Cost is the position cost.
	Cost decimal.Decimal `json:"cost"`

	// Fee is the opening fee.
	Fee decimal.Decimal `json:"fee"`

	// Vol is the position size.
	Vol decimal.Decimal `json:"vol"`

	// VolClosed is the volume closed.
	VolClosed decimal.Decimal `json:"vol_closed"`

	// Margin is the margin used.
	Margin decimal.Decimal `json:"margin"`

	// Value is the current value.
	Value decimal.Decimal `json:"value,omitempty"`

	// Net is the net P&L.
	Net decimal.Decimal `json:"net,omitempty"`

	// Terms is the margin call terms.
	Terms string `json:"terms,omitempty"`

	// RolloverTime is the next rollover time.
	RolloverTime float64 `json:"rollovertm,omitempty"`

	// Misc contains miscellaneous info.
	Misc string `json:"misc,omitempty"`

	// OFlags contains order flags.
	OFlags string `json:"oflags,omitempty"`
}

// ExtendedBalance represents extended balance information for an asset.
type ExtendedBalance struct {
	// Balance is the total balance amount.
	Balance string `json:"balance"`

	// Credit is the total credit amount (if account has credit line).
	Credit string `json:"credit,omitempty"`

	// CreditUsed is the used credit amount.
	CreditUsed string `json:"credit_used,omitempty"`

	// HoldTrade is the total held amount for trading.
	HoldTrade string `json:"hold_trade,omitempty"`
}

// CreditLineAsset represents credit line details for a specific asset.
type CreditLineAsset struct {
	// Balance is the current balance for the asset.
	Balance string `json:"balance"`

	// CreditLimit is the credit limit for the asset.
	CreditLimit string `json:"credit_limit"`

	// CreditUsed is the currently used credit.
	CreditUsed string `json:"credit_used"`

	// AvailableCredit is the available credit.
	AvailableCredit string `json:"available_credit"`
}

// CreditLimitsMonitor contains credit monitoring metrics.
type CreditLimitsMonitor struct {
	// TotalCreditUSD is total credit across all assets in USD.
	TotalCreditUSD string `json:"total_credit_usd,omitempty"`

	// TotalCreditUsedUSD is total credit used in USD.
	TotalCreditUsedUSD string `json:"total_credit_used_usd,omitempty"`

	// TotalCollateralValueUSD is sum of asset balance * collateral in USD.
	TotalCollateralValueUSD string `json:"total_collateral_value_usd,omitempty"`

	// EquityUSD is total collateral - total credit in USD.
	EquityUSD string `json:"equity_usd,omitempty"`

	// OngoingBalance is total collateral / total credit.
	OngoingBalance string `json:"ongoing_balance,omitempty"`

	// DebtToEquity is total credit used / equity.
	DebtToEquity string `json:"debt_to_equity,omitempty"`
}

// CreditLines represents credit line details for VIP accounts.
type CreditLines struct {
	// AssetDetails contains credit line details by asset.
	AssetDetails map[string]CreditLineAsset `json:"asset_details"`

	// LimitsMonitor contains credit monitoring metrics.
	LimitsMonitor CreditLimitsMonitor `json:"limits_monitor"`
}
