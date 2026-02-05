package kraken

import (
	"context"
	"net/url"
	"strconv"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

// GetStrategiesOptions contains options for GetStrategies.
type GetStrategiesOptions struct {
	// Asset filters strategies by asset.
	Asset string
	// LockType filters strategies by lock type.
	LockType string
	// Cursor is the pagination cursor.
	Cursor string
	// Limit is the maximum number of results.
	Limit int
}

// GetStrategies retrieves available earn strategies.
//
// API: POST /0/private/Earn/Strategies
// Docs: https://docs.kraken.com/api/docs/rest-api/list-earn-strategies
func (s *EarnService) GetStrategies(ctx context.Context, opts *GetStrategiesOptions) (*types.StrategiesResult, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.LockType != "" {
			params.Set("lock_type", opts.LockType)
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Earn/Strategies", params)
	if err != nil {
		return nil, err
	}

	var result types.StrategiesResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllocationsOptions contains options for GetAllocations.
type GetAllocationsOptions struct {
	// ConvertedAsset is the asset to convert values to.
	ConvertedAsset string
	// HideZeroAllocations hides allocations with zero balance.
	HideZeroAllocations bool
	// AscendingOrder returns results in ascending order.
	AscendingOrder bool
}

// GetAllocations retrieves the user's earn allocations.
//
// API: POST /0/private/Earn/Allocations
// Docs: https://docs.kraken.com/api/docs/rest-api/list-earn-allocations
func (s *EarnService) GetAllocations(ctx context.Context, opts *GetAllocationsOptions) (*types.AllocationsResult, error) {
	params := url.Values{}
	if opts != nil {
		if opts.ConvertedAsset != "" {
			params.Set("converted_asset", opts.ConvertedAsset)
		}
		if opts.HideZeroAllocations {
			params.Set("hide_zero_allocations", "true")
		}
		if opts.AscendingOrder {
			params.Set("ascending", "true")
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Earn/Allocations", params)
	if err != nil {
		return nil, err
	}

	var result types.AllocationsResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AllocateRequest contains parameters for Allocate.
type AllocateRequest struct {
	// StrategyID is the strategy to allocate to.
	StrategyID string
	// Amount is the amount to allocate.
	Amount string
}

// Allocate allocates funds to an earn strategy.
//
// API: POST /0/private/Earn/Allocate
// Docs: https://docs.kraken.com/api/docs/rest-api/allocate-earn-funds
func (s *EarnService) Allocate(ctx context.Context, req *AllocateRequest) (*types.AllocateResult, error) {
	params := url.Values{}
	if req != nil {
		if req.StrategyID != "" {
			params.Set("strategy_id", req.StrategyID)
		}
		if req.Amount != "" {
			params.Set("amount", req.Amount)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Earn/Allocate", params)
	if err != nil {
		return nil, err
	}

	var result types.AllocateResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeallocateRequest contains parameters for Deallocate.
type DeallocateRequest struct {
	// StrategyID is the strategy to deallocate from.
	StrategyID string
	// Amount is the amount to deallocate.
	Amount string
}

// Deallocate deallocates funds from an earn strategy.
//
// API: POST /0/private/Earn/Deallocate
// Docs: https://docs.kraken.com/api/docs/rest-api/deallocate-earn-funds
func (s *EarnService) Deallocate(ctx context.Context, req *DeallocateRequest) (*types.DeallocateResult, error) {
	params := url.Values{}
	if req != nil {
		if req.StrategyID != "" {
			params.Set("strategy_id", req.StrategyID)
		}
		if req.Amount != "" {
			params.Set("amount", req.Amount)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Earn/Deallocate", params)
	if err != nil {
		return nil, err
	}

	var result types.DeallocateResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllocationStatus retrieves the status of the last allocation request.
//
// API: POST /0/private/Earn/AllocateStatus
// Docs: https://docs.kraken.com/api/docs/rest-api/get-allocation-status
func (s *EarnService) GetAllocationStatus(ctx context.Context, strategyID string) (*types.AllocationStatus, error) {
	params := url.Values{}
	params.Set("strategy_id", strategyID)

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Earn/AllocateStatus", params)
	if err != nil {
		return nil, err
	}

	var result types.AllocationStatus
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeallocationStatus retrieves the status of the last deallocation request.
//
// API: POST /0/private/Earn/DeallocateStatus
// Docs: https://docs.kraken.com/api/docs/rest-api/get-deallocation-status
func (s *EarnService) GetDeallocationStatus(ctx context.Context, strategyID string) (*types.AllocationStatus, error) {
	params := url.Values{}
	params.Set("strategy_id", strategyID)

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Earn/DeallocateStatus", params)
	if err != nil {
		return nil, err
	}

	var result types.AllocationStatus
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
