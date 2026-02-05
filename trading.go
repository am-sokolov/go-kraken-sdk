package kraken

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

// AddOrderRequest contains parameters for placing a new order.
type AddOrderRequest struct {
	// Required fields
	OrderType types.OrderType `json:"ordertype"`
	Side      types.Side      `json:"type"`
	Volume    string          `json:"volume"`
	Pair      string          `json:"pair"`

	// Optional price fields
	Price  string `json:"price,omitempty"`
	Price2 string `json:"price2,omitempty"`

	// Optional order options
	Trigger     string `json:"trigger,omitempty"`
	Leverage    string `json:"leverage,omitempty"`
	ReduceOnly  bool   `json:"reduce_only,omitempty"`
	STPType     string `json:"stptype,omitempty"`
	OFlags      string `json:"oflags,omitempty"`
	TimeInForce string `json:"timeinforce,omitempty"`
	StartTime   string `json:"starttm,omitempty"`
	ExpireTime  string `json:"expiretm,omitempty"`
	Deadline    string `json:"deadline,omitempty"`

	// Asset class (required for xstocks trading)
	AssetClass string `json:"asset_class,omitempty"`

	// Client identifiers
	UserRef       int64  `json:"userref,omitempty"`
	ClientOrderID string `json:"cl_ord_id,omitempty"`

	// Conditional close order
	CloseOrderType string `json:"close[ordertype],omitempty"`
	ClosePrice     string `json:"close[price],omitempty"`
	ClosePrice2    string `json:"close[price2],omitempty"`

	// Iceberg orders
	DisplayVolume string `json:"displayvol,omitempty"`

	// Validation only
	Validate bool `json:"validate,omitempty"`
}

// AddOrderResult contains the result of placing an order.
type AddOrderResult struct {
	Description struct {
		Order string `json:"order"`
		Close string `json:"close,omitempty"`
	} `json:"descr"`
	TransactionIDs []string `json:"txid"`
}

// AddOrder places a new order.
//
// API: POST /0/private/AddOrder
// Docs: https://docs.kraken.com/api/docs/rest-api/add-order
func (s *TradingService) AddOrder(ctx context.Context, req *AddOrderRequest) (*AddOrderResult, error) {
	if req == nil {
		return nil, fmt.Errorf("add order request is nil")
	}

	// Required parameters per Kraken docs.
	if req.OrderType == "" {
		return nil, fmt.Errorf("ordertype is required")
	}
	if req.Side == "" {
		return nil, fmt.Errorf("type is required")
	}
	if req.Volume == "" {
		return nil, fmt.Errorf("volume is required")
	}
	if req.Pair == "" {
		return nil, fmt.Errorf("pair is required")
	}

	// Kraken REST AddOrder docs: userref and cl_ord_id are mutually exclusive.
	if req.UserRef != 0 && req.ClientOrderID != "" {
		return nil, fmt.Errorf("userref and cl_ord_id are mutually exclusive")
	}

	resp, err := s.client.restClient.DoPrivateJSON(ctx, "/0/private/AddOrder", req)
	if err != nil {
		return nil, err
	}

	var result AddOrderResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddOrderBatchRequest contains a batch of orders to place.
type AddOrderBatchRequest struct {
	// Pair is the asset pair for all orders in the batch.
	Pair string
	// Orders is the list of orders to place (max 15).
	Orders []AddOrderRequest
	// Deadline is the RFC3339 deadline for the batch.
	Deadline string
	// Validate only validates orders without placing them.
	Validate bool
}

// AddOrderBatchResult contains the result of a batch order.
type AddOrderBatchResult struct {
	Orders []struct {
		Description struct {
			Order string `json:"order"`
			Close string `json:"close,omitempty"`
		} `json:"descr"`
		TransactionID string `json:"txid,omitempty"`
		Error         string `json:"error,omitempty"`
	} `json:"orders"`
}

// AddOrderBatch places multiple orders in a single request.
//
// API: POST /0/private/AddOrderBatch
// Docs: https://docs.kraken.com/api/docs/rest-api/add-order-batch
func (s *TradingService) AddOrderBatch(ctx context.Context, req *AddOrderBatchRequest) (*AddOrderBatchResult, error) {
	params := url.Values{}
	params.Set("pair", req.Pair)

	for i, order := range req.Orders {
		prefix := fmt.Sprintf("orders[%d]", i)
		params.Set(prefix+"[ordertype]", string(order.OrderType))
		params.Set(prefix+"[type]", string(order.Side))
		params.Set(prefix+"[volume]", order.Volume)
		if order.Price != "" {
			params.Set(prefix+"[price]", order.Price)
		}
		if order.Price2 != "" {
			params.Set(prefix+"[price2]", order.Price2)
		}
		if order.UserRef != 0 {
			params.Set(prefix+"[userref]", strconv.FormatInt(order.UserRef, 10))
		}
		if order.OFlags != "" {
			params.Set(prefix+"[oflags]", order.OFlags)
		}
	}

	if req.Deadline != "" {
		params.Set("deadline", req.Deadline)
	}
	if req.Validate {
		params.Set("validate", "true")
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/AddOrderBatch", params)
	if err != nil {
		return nil, err
	}

	var result AddOrderBatchResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditOrderRequest contains parameters for editing an existing order.
type EditOrderRequest struct {
	// TxID is the transaction ID of the order to edit.
	TxID string
	// Pair is the asset pair.
	Pair string
	// Optional new values
	Volume         string
	DisplayVolume  string
	Price          string
	Price2         string
	OFlags         string
	Deadline       string
	CancelResponse bool
	Validate       bool
	UserRef        int64
}

// EditOrderResult contains the result of editing an order.
type EditOrderResult struct {
	Description struct {
		Order string `json:"order"`
	} `json:"descr"`
	NewTransactionID      string `json:"txid"`
	OriginalTransactionID string `json:"originaltxid"`
	Volume                string `json:"volume"`
	Price                 string `json:"price"`
	Price2                string `json:"price2"`
	OrdersCancelled       int    `json:"orders_cancelled"`
	Status                string `json:"status"`
}

// EditOrder edits an existing order (cancel and replace).
//
// API: POST /0/private/EditOrder
// Docs: https://docs.kraken.com/api/docs/rest-api/edit-order
func (s *TradingService) EditOrder(ctx context.Context, req *EditOrderRequest) (*EditOrderResult, error) {
	params := url.Values{}
	params.Set("txid", req.TxID)
	params.Set("pair", req.Pair)

	if req.Volume != "" {
		params.Set("volume", req.Volume)
	}
	if req.DisplayVolume != "" {
		params.Set("displayvol", req.DisplayVolume)
	}
	if req.Price != "" {
		params.Set("price", req.Price)
	}
	if req.Price2 != "" {
		params.Set("price2", req.Price2)
	}
	if req.OFlags != "" {
		params.Set("oflags", req.OFlags)
	}
	if req.Deadline != "" {
		params.Set("deadline", req.Deadline)
	}
	if req.CancelResponse {
		params.Set("cancel_response", "true")
	}
	if req.Validate {
		params.Set("validate", "true")
	}
	if req.UserRef != 0 {
		params.Set("userref", strconv.FormatInt(req.UserRef, 10))
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/EditOrder", params)
	if err != nil {
		return nil, err
	}

	var result EditOrderResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AmendOrderRequest contains parameters for amending an order in-place.
type AmendOrderRequest struct {
	// TxID is the transaction ID of the order to amend.
	TxID string
	// Optional new values
	OrderQty     string
	DisplayQty   string
	LimitPrice   string
	TriggerPrice string
	PostOnly     *bool
	Deadline     string
}

// AmendOrderResult contains the result of amending an order.
type AmendOrderResult struct {
	AmendID string `json:"amend_id"`
}

// AmendOrder amends an order in-place without losing queue priority.
//
// API: POST /0/private/AmendOrder
// Docs: https://docs.kraken.com/api/docs/rest-api/amend-order
func (s *TradingService) AmendOrder(ctx context.Context, req *AmendOrderRequest) (*AmendOrderResult, error) {
	params := url.Values{}
	params.Set("txid", req.TxID)

	if req.OrderQty != "" {
		params.Set("order_qty", req.OrderQty)
	}
	if req.DisplayQty != "" {
		params.Set("display_qty", req.DisplayQty)
	}
	if req.LimitPrice != "" {
		params.Set("limit_price", req.LimitPrice)
	}
	if req.TriggerPrice != "" {
		params.Set("trigger_price", req.TriggerPrice)
	}
	if req.PostOnly != nil {
		params.Set("post_only", strconv.FormatBool(*req.PostOnly))
	}
	if req.Deadline != "" {
		params.Set("deadline", req.Deadline)
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/AmendOrder", params)
	if err != nil {
		return nil, err
	}

	var result AmendOrderResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelOrderResult contains the result of cancelling an order.
type CancelOrderResult struct {
	Count   int  `json:"count"`
	Pending bool `json:"pending,omitempty"`
}

// CancelOrder cancels an open order.
//
// API: POST /0/private/CancelOrder
// Docs: https://docs.kraken.com/api/docs/rest-api/cancel-order
func (s *TradingService) CancelOrder(ctx context.Context, txid string) (*CancelOrderResult, error) {
	params := url.Values{}
	params.Set("txid", txid)

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CancelOrder", params)
	if err != nil {
		return nil, err
	}

	var result CancelOrderResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelOrderByClientID cancels an order by client order ID.
//
// API: POST /0/private/CancelOrder
// Docs: https://docs.kraken.com/api/docs/rest-api/cancel-order
func (s *TradingService) CancelOrderByClientID(ctx context.Context, clientOrderID string) (*CancelOrderResult, error) {
	params := url.Values{}
	params.Set("cl_ord_id", clientOrderID)

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CancelOrder", params)
	if err != nil {
		return nil, err
	}

	var result CancelOrderResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelOrderBatch cancels multiple orders.
//
// API: POST /0/private/CancelOrderBatch
// Docs: https://docs.kraken.com/api/docs/rest-api/cancel-order-batch
func (s *TradingService) CancelOrderBatch(ctx context.Context, txids []string) (*CancelOrderResult, error) {
	params := url.Values{}
	for i, txid := range txids {
		params.Set(fmt.Sprintf("orders[%d]", i), txid)
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CancelOrderBatch", params)
	if err != nil {
		return nil, err
	}

	var result CancelOrderResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelAllOrdersResult contains the result of cancelling all orders.
type CancelAllOrdersResult struct {
	Count int `json:"count"`
}

// CancelAllOrders cancels all open orders.
//
// API: POST /0/private/CancelAll
// Docs: https://docs.kraken.com/api/docs/rest-api/cancel-all-orders
func (s *TradingService) CancelAllOrders(ctx context.Context) (*CancelAllOrdersResult, error) {
	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CancelAll", url.Values{})
	if err != nil {
		return nil, err
	}

	var result CancelAllOrdersResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelAllOrdersAfterResult contains the result of setting dead man's switch.
type CancelAllOrdersAfterResult struct {
	CurrentTime string `json:"currentTime"`
	TriggerTime string `json:"triggerTime"`
}

// CancelAllOrdersAfter sets a dead man's switch to cancel all orders after timeout.
// Set timeout to 0 to disable the dead man's switch.
//
// API: POST /0/private/CancelAllOrdersAfter
// Docs: https://docs.kraken.com/api/docs/rest-api/cancel-all-orders-after
func (s *TradingService) CancelAllOrdersAfter(ctx context.Context, timeoutSeconds int) (*CancelAllOrdersAfterResult, error) {
	params := url.Values{}
	params.Set("timeout", strconv.Itoa(timeoutSeconds))

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CancelAllOrdersAfter", params)
	if err != nil {
		return nil, err
	}

	var result CancelAllOrdersAfterResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrderAmendsOptions contains options for GetOrderAmends.
type GetOrderAmendsOptions struct {
	// OrderID is the order ID to query amends for.
	OrderID string
	// Start is the starting timestamp.
	Start int64
	// End is the ending timestamp.
	End int64
}

// OrderAmend represents an order amendment record.
type OrderAmend struct {
	AmendID       string `json:"amend_id"`
	OrderID       string `json:"order_id"`
	Timestamp     string `json:"timestamp"`
	AmendType     string `json:"amend_type"`
	OriginalQty   string `json:"original_qty,omitempty"`
	NewQty        string `json:"new_qty,omitempty"`
	OriginalPrice string `json:"original_price,omitempty"`
	NewPrice      string `json:"new_price,omitempty"`
}

// GetOrderAmends retrieves order amendment history.
//
// API: POST /0/private/OrderAmends
// Docs: https://docs.kraken.com/api/docs/rest-api/get-order-amends
func (s *TradingService) GetOrderAmends(ctx context.Context, opts *GetOrderAmendsOptions) ([]OrderAmend, error) {
	params := url.Values{}
	if opts != nil {
		if opts.OrderID != "" {
			params.Set("order_id", opts.OrderID)
		}
		if opts.Start != 0 {
			params.Set("start", strconv.FormatInt(opts.Start, 10))
		}
		if opts.End != 0 {
			params.Set("end", strconv.FormatInt(opts.End, 10))
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/OrderAmends", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Amends []OrderAmend `json:"amends"`
	}
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result.Amends, nil
}

// Ensure TradingService is used
var _ = (*TradingService)(nil)
