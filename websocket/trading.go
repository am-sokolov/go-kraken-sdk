package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AddOrder places a new order via WebSocket.
func (c *Client) AddOrder(ctx context.Context, params AddOrderParams) (*AddOrderResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	// Required per Kraken WS v2 add_order docs.
	if strings.TrimSpace(params.Symbol) == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	if strings.TrimSpace(params.Side) == "" {
		return nil, fmt.Errorf("side is required")
	}
	if strings.TrimSpace(params.OrderType) == "" {
		return nil, fmt.Errorf("order_type is required")
	}
	if strings.TrimSpace(params.OrderQty) == "" {
		return nil, fmt.Errorf("order_qty is required")
	}

	params.Token = c.token
	reqID := c.nextReqID()

	req := struct {
		Method MessageType    `json:"method"`
		ReqID  int64          `json:"req_id"`
		Params AddOrderParams `json:"params"`
	}{
		Method: MsgTypeAddOrder,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("add_order error: %s", resp.Error)
	}

	var result AddOrderResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &result, nil
}

// EditOrder edits an existing order via WebSocket.
func (c *Client) EditOrder(ctx context.Context, params EditOrderParams) (*EditOrderResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	params.Token = c.token
	reqID := c.nextReqID()

	req := struct {
		Method MessageType     `json:"method"`
		ReqID  int64           `json:"req_id"`
		Params EditOrderParams `json:"params"`
	}{
		Method: MsgTypeEditOrder,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("edit_order error: %s", resp.Error)
	}

	var result EditOrderResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &result, nil
}

// AmendOrder amends an existing order in-place via WebSocket.
func (c *Client) AmendOrder(ctx context.Context, params AmendOrderParams) (*AmendOrderResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	params.Token = c.token
	reqID := c.nextReqID()

	req := struct {
		Method MessageType      `json:"method"`
		ReqID  int64            `json:"req_id"`
		Params AmendOrderParams `json:"params"`
	}{
		Method: MsgTypeAmendOrder,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("amend_order error: %s", resp.Error)
	}

	var result AmendOrderResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &result, nil
}

// CancelOrder cancels one or more orders via WebSocket.
func (c *Client) CancelOrder(ctx context.Context, params CancelOrderParams) ([]CancelOrderResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	params.Token = c.token
	reqID := c.nextReqID()

	req := struct {
		Method MessageType       `json:"method"`
		ReqID  int64             `json:"req_id"`
		Params CancelOrderParams `json:"params"`
	}{
		Method: MsgTypeCancelOrder,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("cancel_order error: %s", resp.Error)
	}

	var result []CancelOrderResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return result, nil
}

// CancelAll cancels all open orders via WebSocket.
func (c *Client) CancelAll(ctx context.Context) (*CancelAllResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	reqID := c.nextReqID()

	req := struct {
		Method MessageType     `json:"method"`
		ReqID  int64           `json:"req_id"`
		Params CancelAllParams `json:"params"`
	}{
		Method: MsgTypeCancelAll,
		ReqID:  reqID,
		Params: CancelAllParams{Token: c.token},
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("cancel_all error: %s", resp.Error)
	}

	var result CancelAllResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &result, nil
}

// BatchAdd places multiple orders in a single request via WebSocket.
func (c *Client) BatchAdd(ctx context.Context, params BatchAddParams) ([]AddOrderResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	params.Token = c.token
	reqID := c.nextReqID()

	req := struct {
		Method MessageType    `json:"method"`
		ReqID  int64          `json:"req_id"`
		Params BatchAddParams `json:"params"`
	}{
		Method: MsgTypeBatchAdd,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("batch_add error: %s", resp.Error)
	}

	var result []AddOrderResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return result, nil
}

// BatchCancel cancels multiple orders in a single request via WebSocket.
func (c *Client) BatchCancel(ctx context.Context, params BatchCancelParams) ([]CancelOrderResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	params.Token = c.token
	reqID := c.nextReqID()

	req := struct {
		Method MessageType       `json:"method"`
		ReqID  int64             `json:"req_id"`
		Params BatchCancelParams `json:"params"`
	}{
		Method: MsgTypeBatchCancel,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("batch_cancel error: %s", resp.Error)
	}

	var result []CancelOrderResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return result, nil
}

// CancelAllOrdersAfter sets a dead man's switch to cancel all orders after a timeout.
// Set timeout to 0 to disable the dead man's switch.
func (c *Client) CancelAllOrdersAfter(ctx context.Context, timeoutSeconds int) (*CancelOnDisconnectResult, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token required")
	}

	reqID := c.nextReqID()

	req := struct {
		Method MessageType              `json:"method"`
		ReqID  int64                    `json:"req_id"`
		Params CancelOnDisconnectParams `json:"params"`
	}{
		Method: MsgTypeCancelOnDisconnect,
		ReqID:  reqID,
		Params: CancelOnDisconnectParams{
			Timeout: timeoutSeconds,
			Token:   c.token,
		},
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("cancel_all_orders_after error: %s", resp.Error)
	}

	var result CancelOnDisconnectResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}

	return &result, nil
}
