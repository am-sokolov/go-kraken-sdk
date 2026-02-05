package types

import "github.com/shopspring/decimal"

// DepositMethod represents a deposit method.
type DepositMethod struct {
	// Method is the method name.
	Method string `json:"method"`

	// Limit is the maximum deposit amount (false if no limit).
	Limit interface{} `json:"limit"`

	// Fee is the deposit fee.
	Fee decimal.Decimal `json:"fee,omitempty"`

	// AddressSetupFee is the fee for setting up a new address.
	AddressSetupFee decimal.Decimal `json:"address-setup-fee,omitempty"`

	// GenAddress indicates if addresses can be generated.
	GenAddress bool `json:"gen-address,omitempty"`

	// Minimum is the minimum deposit amount.
	Minimum decimal.Decimal `json:"minimum,omitempty"`
}

// DepositAddress represents a deposit address.
type DepositAddress struct {
	// Address is the deposit address.
	Address string `json:"address"`

	// ExpiretM is the expiration timestamp.
	ExpireTM string `json:"expiretm,omitempty"`

	// New indicates if this is a newly generated address.
	New bool `json:"new,omitempty"`

	// Tag is the destination tag/memo (for some assets).
	Tag string `json:"tag,omitempty"`
}

// DepositStatus represents a deposit status entry.
type DepositStatus struct {
	// Method is the deposit method.
	Method string `json:"method"`

	// AClass is the asset class.
	AClass string `json:"aclass"`

	// Asset is the asset name.
	Asset string `json:"asset"`

	// RefID is the reference ID.
	RefID string `json:"refid"`

	// TxID is the blockchain transaction ID.
	TxID string `json:"txid"`

	// Info is additional info.
	Info string `json:"info"`

	// Amount is the deposit amount.
	Amount decimal.Decimal `json:"amount"`

	// Fee is the deposit fee.
	Fee decimal.Decimal `json:"fee"`

	// Time is the deposit timestamp.
	Time int64 `json:"time"`

	// Status is the deposit status.
	Status string `json:"status"`

	// StatusProp is additional status properties.
	StatusProp string `json:"status-prop,omitempty"`

	// Originators is the list of originating addresses.
	Originators []string `json:"originators,omitempty"`
}

// WithdrawalMethod represents a withdrawal method.
type WithdrawalMethod struct {
	// Asset is the asset name.
	Asset string `json:"asset"`

	// Method is the method name.
	Method string `json:"method"`

	// Network is the network name.
	Network string `json:"network,omitempty"`

	// Minimum is the minimum withdrawal amount.
	Minimum decimal.Decimal `json:"minimum"`

	// Maximum is the maximum withdrawal amount.
	Maximum decimal.Decimal `json:"maximum"`

	// Fee is the withdrawal fee.
	Fee decimal.Decimal `json:"fee"`
}

// WithdrawalAddress represents a saved withdrawal address.
type WithdrawalAddress struct {
	// Address is the withdrawal address.
	Address string `json:"address"`

	// Asset is the asset name.
	Asset string `json:"asset"`

	// Method is the withdrawal method.
	Method string `json:"method"`

	// Key is the address key/name.
	Key string `json:"key"`

	// Memo is the destination tag/memo.
	Memo string `json:"memo,omitempty"`

	// Verified indicates if the address is verified.
	Verified bool `json:"verified"`
}

// WithdrawalInfo represents withdrawal information.
type WithdrawalInfo struct {
	// Method is the withdrawal method.
	Method string `json:"method"`

	// Limit is the withdrawal limit.
	Limit decimal.Decimal `json:"limit"`

	// Amount is the amount to withdraw.
	Amount decimal.Decimal `json:"amount"`

	// Fee is the withdrawal fee.
	Fee decimal.Decimal `json:"fee"`
}

// WithdrawalStatus represents a withdrawal status entry.
type WithdrawalStatus struct {
	// Method is the withdrawal method.
	Method string `json:"method"`

	// AClass is the asset class.
	AClass string `json:"aclass"`

	// Asset is the asset name.
	Asset string `json:"asset"`

	// RefID is the reference ID.
	RefID string `json:"refid"`

	// TxID is the blockchain transaction ID.
	TxID string `json:"txid"`

	// Info is additional info.
	Info string `json:"info"`

	// Amount is the withdrawal amount.
	Amount decimal.Decimal `json:"amount"`

	// Fee is the withdrawal fee.
	Fee decimal.Decimal `json:"fee"`

	// Time is the withdrawal timestamp.
	Time int64 `json:"time"`

	// Status is the withdrawal status.
	Status string `json:"status"`

	// StatusProp is additional status properties.
	StatusProp string `json:"status-prop,omitempty"`

	// Key is the address key/name.
	Key string `json:"key,omitempty"`
}

// WithdrawalResult represents the result of a withdrawal request.
type WithdrawalResult struct {
	// RefID is the reference ID.
	RefID string `json:"refid"`
}

// TransferResult represents the result of an account transfer.
type TransferResult struct {
	// RefID is the reference ID.
	RefID string `json:"refid"`
}
