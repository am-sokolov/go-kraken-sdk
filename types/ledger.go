package types

import "github.com/shopspring/decimal"

// LedgerType represents the type of ledger entry.
type LedgerType string

const (
	LedgerTypeDeposit    LedgerType = "deposit"
	LedgerTypeWithdrawal LedgerType = "withdrawal"
	LedgerTypeTrade      LedgerType = "trade"
	LedgerTypeMargin     LedgerType = "margin"
	LedgerTypeRollover   LedgerType = "rollover"
	LedgerTypeSpend      LedgerType = "spend"
	LedgerTypeReceive    LedgerType = "receive"
	LedgerTypeSettled    LedgerType = "settled"
	LedgerTypeCredit     LedgerType = "credit"
	LedgerTypeTransfer   LedgerType = "transfer"
	LedgerTypeAdjustment LedgerType = "adjustment"
	LedgerTypeStaking    LedgerType = "staking"
	LedgerTypeReward     LedgerType = "reward"
	LedgerTypeDividend   LedgerType = "dividend"
	LedgerTypeSale       LedgerType = "sale"
	LedgerTypeConversion LedgerType = "conversion"
	LedgerTypeNftRebate  LedgerType = "nfttrade"
	LedgerTypeNftTrade   LedgerType = "nftcredit"
	LedgerTypeCustody    LedgerType = "custodytransfer"
)

// Ledger represents a ledger entry.
type Ledger struct {
	// RefID is the reference ID.
	RefID string `json:"refid"`

	// Time is the entry timestamp.
	Time float64 `json:"time"`

	// Type is the entry type.
	Type LedgerType `json:"type"`

	// SubType is the entry subtype.
	SubType string `json:"subtype"`

	// AClass is the asset class.
	AClass string `json:"aclass"`

	// Asset is the asset name.
	Asset string `json:"asset"`

	// Amount is the transaction amount.
	Amount decimal.Decimal `json:"amount"`

	// Fee is the transaction fee.
	Fee decimal.Decimal `json:"fee"`

	// Balance is the resulting balance.
	Balance decimal.Decimal `json:"balance"`
}

// LedgersResult contains the result of a ledgers query.
type LedgersResult struct {
	// Ledgers is a map of ledger ID to ledger entry.
	Ledgers map[string]Ledger `json:"ledger"`

	// Count is the total number of entries.
	Count int `json:"count"`
}

// TradeVolume represents the trade volume and fee structure.
type TradeVolume struct {
	// Currency is the volume currency.
	Currency string `json:"currency"`

	// Volume is the 30-day volume.
	Volume decimal.Decimal `json:"volume"`

	// Fees contains fee info per pair.
	Fees map[string]FeeInfo `json:"fees,omitempty"`

	// FeesMaker contains maker fee info per pair.
	FeesMaker map[string]FeeInfo `json:"fees_maker,omitempty"`
}

// FeeInfo contains fee information for a pair.
type FeeInfo struct {
	// Fee is the current fee percentage.
	Fee decimal.Decimal `json:"fee"`

	// MinFee is the minimum fee.
	MinFee decimal.Decimal `json:"minfee"`

	// MaxFee is the maximum fee.
	MaxFee decimal.Decimal `json:"maxfee"`

	// NextFee is the fee at next volume tier.
	NextFee decimal.Decimal `json:"nextfee,omitempty"`

	// TierVolume is the volume required for next tier.
	TierVolume decimal.Decimal `json:"tiervolume,omitempty"`

	// NextVolume is the volume required for next tier.
	NextVolume decimal.Decimal `json:"nextvolume,omitempty"`
}
