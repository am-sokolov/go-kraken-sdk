package types

import "github.com/shopspring/decimal"

// EarnStrategy represents an earn strategy.
type EarnStrategy struct {
	// ID is the strategy identifier.
	ID string `json:"id"`

	// Asset is the asset for this strategy.
	Asset string `json:"asset"`

	// LockType is the lock type.
	LockType EarnLockType `json:"lock_type"`

	// APREstimate is the estimated annual percentage rate.
	APREstimate *APREstimate `json:"apr_estimate,omitempty"`

	// UserMinAllocation is the minimum allocation.
	UserMinAllocation decimal.Decimal `json:"user_min_allocation,omitempty"`

	// AllocationFee is the allocation fee.
	AllocationFee decimal.Decimal `json:"allocation_fee,omitempty"`

	// DeallocationFee is the deallocation fee.
	DeallocationFee decimal.Decimal `json:"deallocation_fee,omitempty"`

	// AutoCompound contains auto-compound settings.
	AutoCompound *AutoCompoundConfig `json:"auto_compound,omitempty"`

	// CanAllocate indicates if allocations are possible.
	CanAllocate bool `json:"can_allocate"`

	// CanDeallocate indicates if deallocations are possible.
	CanDeallocate bool `json:"can_deallocate"`

	// AllocationRestrictionInfo contains allocation restrictions.
	AllocationRestrictionInfo []string `json:"allocation_restriction_info,omitempty"`

	// YieldSource is the source of yield.
	YieldSource *YieldSource `json:"yield_source,omitempty"`
}

// EarnLockType represents the type of lock for an earn strategy.
type EarnLockType struct {
	// Type is the lock type (flex, bonded, instant, etc.).
	Type string `json:"type"`

	// PayoutFrequency is how often payouts occur.
	PayoutFrequency int `json:"payout_frequency,omitempty"`

	// BondingPeriod is the bonding period in seconds.
	BondingPeriod int64 `json:"bonding_period,omitempty"`

	// BondingPeriodVariable indicates if bonding period varies.
	BondingPeriodVariable bool `json:"bonding_period_variable,omitempty"`

	// BondingRewards indicates if rewards during bonding.
	BondingRewards bool `json:"bonding_rewards,omitempty"`

	// ExitQueuePeriod is the exit queue period in seconds.
	ExitQueuePeriod int64 `json:"exit_queue_period,omitempty"`

	// UnbondingPeriod is the unbonding period in seconds.
	UnbondingPeriod int64 `json:"unbonding_period,omitempty"`

	// UnbondingPeriodVariable indicates if unbonding period varies.
	UnbondingPeriodVariable bool `json:"unbonding_period_variable,omitempty"`

	// UnbondingRewards indicates if rewards during unbonding.
	UnbondingRewards bool `json:"unbonding_rewards,omitempty"`
}

// APREstimate represents an APR estimate range.
type APREstimate struct {
	// Low is the low estimate.
	Low decimal.Decimal `json:"low"`

	// High is the high estimate.
	High decimal.Decimal `json:"high"`
}

// AutoCompoundConfig represents auto-compound configuration.
type AutoCompoundConfig struct {
	// Default indicates the default auto-compound setting.
	Default bool `json:"default"`

	// Type indicates the type of auto-compound (e.g., "enabled", "disabled", "optional").
	Type string `json:"type"`
}

// YieldSource represents the source of yield.
type YieldSource struct {
	// Type is the yield source type.
	Type string `json:"type"`
}

// EarnAllocation represents an allocation to an earn strategy.
type EarnAllocation struct {
	// StrategyID is the strategy identifier.
	StrategyID string `json:"strategy_id"`

	// NativeAsset is the native asset name.
	NativeAsset string `json:"native_asset"`

	// AmountAllocated is the allocated amount.
	AmountAllocated AllocationAmount `json:"amount_allocated"`

	// TotalRewarded is the total rewards earned.
	TotalRewarded AllocationAmount `json:"total_rewarded"`

	// Payout is the payout information.
	Payout *EarnPayout `json:"payout,omitempty"`
}

// AllocationAmount represents an amount with bonding details.
// Each field contains native and converted values.
type AllocationAmount struct {
	// Bonding is the amount bonding.
	Bonding *ConvertedAmount `json:"bonding,omitempty"`

	// ExitQueue is the amount in exit queue.
	ExitQueue *ConvertedAmount `json:"exit_queue,omitempty"`

	// Pending is the pending amount.
	Pending *ConvertedAmount `json:"pending,omitempty"`

	// Total is the total amount.
	Total *ConvertedAmount `json:"total,omitempty"`

	// Unbonding is the amount unbonding.
	Unbonding *ConvertedAmount `json:"unbonding,omitempty"`

	// Allocated is the allocated amount.
	Allocated *ConvertedAmount `json:"allocated,omitempty"`
}

// EarnPayout represents payout information.
type EarnPayout struct {
	// AccumulatedReward is the accumulated reward.
	AccumulatedReward *ConvertedAmount `json:"accumulated_reward,omitempty"`

	// EstimatedReward is the estimated reward.
	EstimatedReward *ConvertedAmount `json:"estimated_reward,omitempty"`

	// PeriodStart is the period start timestamp.
	PeriodStart string `json:"period_start"`

	// PeriodEnd is the period end timestamp.
	PeriodEnd string `json:"period_end"`
}

// AllocateResult represents the result of an allocation request.
type AllocateResult struct {
	// Pending indicates if the allocation is pending.
	Pending bool `json:"pending"`
}

// DeallocateResult represents the result of a deallocation request.
type DeallocateResult struct {
	// Pending indicates if the deallocation is pending.
	Pending bool `json:"pending"`
}

// AllocationStatus represents the status of an allocation operation.
type AllocationStatus struct {
	// Pending indicates if there's a pending operation.
	Pending bool `json:"pending"`
}

// StrategiesResult contains the result of a strategies query.
type StrategiesResult struct {
	// Items is the list of strategies.
	Items []EarnStrategy `json:"items"`

	// NextCursor is the cursor for pagination.
	NextCursor string `json:"next_cursor,omitempty"`
}

// ConvertedAmount represents an amount with native and converted values.
type ConvertedAmount struct {
	// Native is the amount in the native asset.
	Native decimal.Decimal `json:"native"`

	// Converted is the amount in the converted asset.
	Converted decimal.Decimal `json:"converted,omitempty"`
}

// AllocationsResult contains the result of an allocations query.
type AllocationsResult struct {
	// ConvertedAsset is the conversion asset.
	ConvertedAsset string `json:"converted_asset,omitempty"`

	// TotalAllocated is the total allocated amount (as simple decimal string).
	TotalAllocated decimal.Decimal `json:"total_allocated"`

	// TotalRewarded is the total rewarded amount (as simple decimal string).
	TotalRewarded decimal.Decimal `json:"total_rewarded"`

	// NextCursor is the cursor for pagination.
	NextCursor string `json:"next_cursor,omitempty"`

	// Items is the list of allocations.
	Items []EarnAllocation `json:"items"`
}
