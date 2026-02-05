package kraken

import (
	"context"
	"net/url"
	"strconv"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

// GetDepositMethodsOptions contains options for GetDepositMethods.
type GetDepositMethodsOptions struct {
	// Asset is the asset to get deposit methods for.
	Asset string
	// AClass is the asset class (optional, defaults to "currency").
	AClass string
}

// GetDepositMethods retrieves deposit methods for a given asset.
//
// API: POST /0/private/DepositMethods
// Docs: https://docs.kraken.com/api/docs/rest-api/get-deposit-methods
func (s *FundingService) GetDepositMethods(ctx context.Context, opts *GetDepositMethodsOptions) ([]types.DepositMethod, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.AClass != "" {
			params.Set("aclass", opts.AClass)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/DepositMethods", params)
	if err != nil {
		return nil, err
	}

	var result []types.DepositMethod
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetDepositAddressesOptions contains options for GetDepositAddresses.
type GetDepositAddressesOptions struct {
	// Asset is the asset to get deposit addresses for.
	Asset string
	// Method is the deposit method name.
	Method string
	// New if true, generates a new address.
	New bool
	// Amount is the amount to deposit (for some methods).
	Amount string
}

// GetDepositAddresses retrieves deposit addresses for a given asset and method.
//
// API: POST /0/private/DepositAddresses
// Docs: https://docs.kraken.com/api/docs/rest-api/get-deposit-addresses
func (s *FundingService) GetDepositAddresses(ctx context.Context, opts *GetDepositAddressesOptions) ([]types.DepositAddress, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.Method != "" {
			params.Set("method", opts.Method)
		}
		if opts.New {
			params.Set("new", "true")
		}
		if opts.Amount != "" {
			params.Set("amount", opts.Amount)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/DepositAddresses", params)
	if err != nil {
		return nil, err
	}

	var result []types.DepositAddress
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetDepositStatusOptions contains options for GetDepositStatus.
type GetDepositStatusOptions struct {
	// Asset is the asset to query status for.
	Asset string
	// Method is the deposit method name.
	Method string
	// Start is the start timestamp.
	Start int64
	// End is the end timestamp.
	End int64
	// Cursor is the pagination cursor.
	Cursor string
	// Limit is the maximum number of results.
	Limit int
}

// GetDepositStatus retrieves the status of recent deposits.
//
// API: POST /0/private/DepositStatus
// Docs: https://docs.kraken.com/api/docs/rest-api/get-status-of-recent-deposits
func (s *FundingService) GetDepositStatus(ctx context.Context, opts *GetDepositStatusOptions) ([]types.DepositStatus, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.Method != "" {
			params.Set("method", opts.Method)
		}
		if opts.Start != 0 {
			params.Set("start", strconv.FormatInt(opts.Start, 10))
		}
		if opts.End != 0 {
			params.Set("end", strconv.FormatInt(opts.End, 10))
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/DepositStatus", params)
	if err != nil {
		return nil, err
	}

	var result []types.DepositStatus
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetWithdrawalMethodsOptions contains options for GetWithdrawalMethods.
type GetWithdrawalMethodsOptions struct {
	// Asset is the asset to get withdrawal methods for.
	Asset string
	// AClass is the asset class (optional, defaults to "currency").
	AClass string
	// Network is the network name.
	Network string
}

// GetWithdrawalMethods retrieves withdrawal methods for a given asset.
//
// API: POST /0/private/WithdrawMethods
// Docs: https://docs.kraken.com/api/docs/rest-api/get-withdrawal-methods
func (s *FundingService) GetWithdrawalMethods(ctx context.Context, opts *GetWithdrawalMethodsOptions) ([]types.WithdrawalMethod, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.AClass != "" {
			params.Set("aclass", opts.AClass)
		}
		if opts.Network != "" {
			params.Set("network", opts.Network)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/WithdrawMethods", params)
	if err != nil {
		return nil, err
	}

	var result []types.WithdrawalMethod
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetWithdrawalAddressesOptions contains options for GetWithdrawalAddresses.
type GetWithdrawalAddressesOptions struct {
	// Asset is the asset to get withdrawal addresses for.
	Asset string
	// AClass is the asset class (optional, defaults to "currency").
	AClass string
	// Method is the withdrawal method name.
	Method string
	// Key is the address key/name.
	Key string
	// Verified if true, only returns verified addresses.
	Verified bool
}

// GetWithdrawalAddresses retrieves saved withdrawal addresses.
//
// API: POST /0/private/WithdrawAddresses
// Docs: https://docs.kraken.com/api/docs/rest-api/get-withdrawal-addresses
func (s *FundingService) GetWithdrawalAddresses(ctx context.Context, opts *GetWithdrawalAddressesOptions) ([]types.WithdrawalAddress, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.AClass != "" {
			params.Set("aclass", opts.AClass)
		}
		if opts.Method != "" {
			params.Set("method", opts.Method)
		}
		if opts.Key != "" {
			params.Set("key", opts.Key)
		}
		if opts.Verified {
			params.Set("verified", "true")
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/WithdrawAddresses", params)
	if err != nil {
		return nil, err
	}

	var result []types.WithdrawalAddress
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetWithdrawalInfoRequest contains parameters for GetWithdrawalInfo.
type GetWithdrawalInfoRequest struct {
	// Asset is the asset to withdraw.
	Asset string
	// Key is the withdrawal address key/name.
	Key string
	// Amount is the amount to withdraw.
	Amount string
}

// GetWithdrawalInfo retrieves withdrawal information for a specific withdrawal.
//
// API: POST /0/private/WithdrawInfo
// Docs: https://docs.kraken.com/api/docs/rest-api/get-withdrawal-info
func (s *FundingService) GetWithdrawalInfo(ctx context.Context, req *GetWithdrawalInfoRequest) (*types.WithdrawalInfo, error) {
	params := url.Values{}
	if req != nil {
		if req.Asset != "" {
			params.Set("asset", req.Asset)
		}
		if req.Key != "" {
			params.Set("key", req.Key)
		}
		if req.Amount != "" {
			params.Set("amount", req.Amount)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/WithdrawInfo", params)
	if err != nil {
		return nil, err
	}

	var result types.WithdrawalInfo
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// WithdrawRequest contains parameters for Withdraw.
type WithdrawRequest struct {
	// Asset is the asset to withdraw.
	Asset string
	// Key is the withdrawal address key/name.
	Key string
	// Amount is the amount to withdraw.
	Amount string
	// Address is the withdrawal address (optional, if not using saved address).
	Address string
	// MaxFee is the maximum fee to allow (optional).
	MaxFee string
}

// Withdraw initiates a withdrawal.
//
// API: POST /0/private/Withdraw
// Docs: https://docs.kraken.com/api/docs/rest-api/withdraw-funds
func (s *FundingService) Withdraw(ctx context.Context, req *WithdrawRequest) (*types.WithdrawalResult, error) {
	params := url.Values{}
	if req != nil {
		if req.Asset != "" {
			params.Set("asset", req.Asset)
		}
		if req.Key != "" {
			params.Set("key", req.Key)
		}
		if req.Amount != "" {
			params.Set("amount", req.Amount)
		}
		if req.Address != "" {
			params.Set("address", req.Address)
		}
		if req.MaxFee != "" {
			params.Set("max_fee", req.MaxFee)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Withdraw", params)
	if err != nil {
		return nil, err
	}

	var result types.WithdrawalResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWithdrawalStatusOptions contains options for GetWithdrawalStatus.
type GetWithdrawalStatusOptions struct {
	// Asset is the asset to query status for.
	Asset string
	// Method is the withdrawal method name.
	Method string
	// Start is the start timestamp.
	Start int64
	// End is the end timestamp.
	End int64
	// Cursor is the pagination cursor.
	Cursor string
	// Limit is the maximum number of results.
	Limit int
}

// GetWithdrawalStatus retrieves the status of recent withdrawals.
//
// API: POST /0/private/WithdrawStatus
// Docs: https://docs.kraken.com/api/docs/rest-api/get-status-of-recent-withdrawals
func (s *FundingService) GetWithdrawalStatus(ctx context.Context, opts *GetWithdrawalStatusOptions) ([]types.WithdrawalStatus, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.Method != "" {
			params.Set("method", opts.Method)
		}
		if opts.Start != 0 {
			params.Set("start", strconv.FormatInt(opts.Start, 10))
		}
		if opts.End != 0 {
			params.Set("end", strconv.FormatInt(opts.End, 10))
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/WithdrawStatus", params)
	if err != nil {
		return nil, err
	}

	var result []types.WithdrawalStatus
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// CancelWithdrawal cancels a pending withdrawal.
//
// API: POST /0/private/WithdrawCancel
// Docs: https://docs.kraken.com/api/docs/rest-api/cancel-withdrawal
func (s *FundingService) CancelWithdrawal(ctx context.Context, asset, refID string) (bool, error) {
	params := url.Values{}
	params.Set("asset", asset)
	params.Set("refid", refID)

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/WithdrawCancel", params)
	if err != nil {
		return false, err
	}

	var result bool
	if err := resp.Decode(&result); err != nil {
		return false, err
	}

	return result, nil
}

// WalletTransferRequest contains parameters for WalletTransfer.
type WalletTransferRequest struct {
	// Asset is the asset to transfer.
	Asset string
	// From is the source wallet.
	From string
	// To is the destination wallet.
	To string
	// Amount is the amount to transfer.
	Amount string
}

// WalletTransfer transfers funds between Kraken wallets.
//
// API: POST /0/private/WalletTransfer
// Docs: https://docs.kraken.com/api/docs/rest-api/wallet-transfer
func (s *FundingService) WalletTransfer(ctx context.Context, req *WalletTransferRequest) (*types.TransferResult, error) {
	params := url.Values{}
	if req != nil {
		if req.Asset != "" {
			params.Set("asset", req.Asset)
		}
		if req.From != "" {
			params.Set("from", req.From)
		}
		if req.To != "" {
			params.Set("to", req.To)
		}
		if req.Amount != "" {
			params.Set("amount", req.Amount)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/WalletTransfer", params)
	if err != nil {
		return nil, err
	}

	var result types.TransferResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
