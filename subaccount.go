package kraken

import (
	"context"
	"net/url"
)

// CreateSubaccountRequest contains parameters for CreateSubaccount.
type CreateSubaccountRequest struct {
	// Username is the username for the new subaccount.
	Username string
	// Email is the email for the new subaccount.
	Email string
}

// CreateSubaccountResult contains the result of creating a subaccount.
type CreateSubaccountResult struct {
	// Result indicates if the creation was successful.
	Result bool `json:"result"`
}

// CreateSubaccount creates a new subaccount.
//
// API: POST /0/private/CreateSubaccount
// Docs: https://docs.kraken.com/api/docs/rest-api/create-sub-account
func (s *SubaccountService) CreateSubaccount(ctx context.Context, req *CreateSubaccountRequest) (*CreateSubaccountResult, error) {
	params := url.Values{}
	if req != nil {
		if req.Username != "" {
			params.Set("username", req.Username)
		}
		if req.Email != "" {
			params.Set("email", req.Email)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CreateSubaccount", params)
	if err != nil {
		return nil, err
	}

	var result CreateSubaccountResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AccountTransferRequest contains parameters for AccountTransfer.
type AccountTransferRequest struct {
	// Asset is the asset to transfer.
	Asset string
	// Amount is the amount to transfer.
	Amount string
	// From is the source account (subaccount ID or "main").
	From string
	// To is the destination account (subaccount ID or "main").
	To string
}

// AccountTransferResult contains the result of an account transfer.
type AccountTransferResult struct {
	// TransferID is the unique ID for this transfer.
	TransferID string `json:"transfer_id"`
	// Status is the status of the transfer.
	Status string `json:"status"`
}

// AccountTransfer transfers funds between accounts (main and subaccounts).
//
// API: POST /0/private/AccountTransfer
// Docs: https://docs.kraken.com/api/docs/rest-api/account-transfer
func (s *SubaccountService) AccountTransfer(ctx context.Context, req *AccountTransferRequest) (*AccountTransferResult, error) {
	params := url.Values{}
	if req != nil {
		if req.Asset != "" {
			params.Set("asset", req.Asset)
		}
		if req.Amount != "" {
			params.Set("amount", req.Amount)
		}
		if req.From != "" {
			params.Set("from", req.From)
		}
		if req.To != "" {
			params.Set("to", req.To)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/AccountTransfer", params)
	if err != nil {
		return nil, err
	}

	var result AccountTransferResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
