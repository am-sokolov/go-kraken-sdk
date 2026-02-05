package websocket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func rawJSONString(value string) json.RawMessage {
	b, _ := json.Marshal(value)
	return b
}

func rawJSONBoolTrue() json.RawMessage {
	return json.RawMessage("true")
}

func rawJSONNumberFromString(value string) (json.RawMessage, error) {
	trim := strings.TrimSpace(value)
	if trim == "" {
		return nil, fmt.Errorf("empty value")
	}

	// Validate that the value is a valid JSON number representation.
	var tmp float64
	if err := json.Unmarshal([]byte(trim), &tmp); err != nil {
		return nil, err
	}

	return json.RawMessage(trim), nil
}

// MarshalJSON ensures numeric fields are encoded as JSON numbers per Kraken WebSocket v2 docs.
func (p AddOrderParams) MarshalJSON() ([]byte, error) {
	orderType := strings.TrimSpace(p.OrderType)
	if orderType == "" {
		return nil, fmt.Errorf("order_type is required")
	}
	side := strings.TrimSpace(p.Side)
	if side == "" {
		return nil, fmt.Errorf("side is required")
	}

	orderQty, err := rawJSONNumberFromString(p.OrderQty)
	if err != nil {
		return nil, fmt.Errorf("order_qty: %w", err)
	}

	obj := map[string]json.RawMessage{
		"order_type": rawJSONString(orderType),
		"side":       rawJSONString(side),
		"order_qty":  orderQty,
	}

	if symbol := strings.TrimSpace(p.Symbol); symbol != "" {
		obj["symbol"] = rawJSONString(symbol)
	}

	if token := strings.TrimSpace(p.Token); token != "" {
		obj["token"] = rawJSONString(token)
	}

	if v := strings.TrimSpace(p.LimitPrice); v != "" {
		raw, err := rawJSONNumberFromString(v)
		if err != nil {
			return nil, fmt.Errorf("limit_price: %w", err)
		}
		obj["limit_price"] = raw
	}
	if v := strings.TrimSpace(p.LimitPriceType); v != "" {
		obj["limit_price_type"] = rawJSONString(v)
	}

	if v := strings.TrimSpace(p.TimeInForce); v != "" {
		// WS v2 docs use lowercase time-in-force values (gtc/gtd/ioc).
		if v == "GTC" || v == "GTD" || v == "IOC" {
			v = strings.ToLower(v)
		}
		obj["time_in_force"] = rawJSONString(v)
	}

	if p.Margin {
		obj["margin"] = rawJSONBoolTrue()
	}
	if p.PostOnly {
		obj["post_only"] = rawJSONBoolTrue()
	}
	if p.ReduceOnly {
		obj["reduce_only"] = rawJSONBoolTrue()
	}

	if v := strings.TrimSpace(p.EffectiveTime); v != "" {
		obj["effective_time"] = rawJSONString(v)
	}
	if v := strings.TrimSpace(p.ExpireTime); v != "" {
		obj["expire_time"] = rawJSONString(v)
	}
	if v := strings.TrimSpace(p.Deadline); v != "" {
		obj["deadline"] = rawJSONString(v)
	}

	if v := strings.TrimSpace(p.ClientOrderID); v != "" {
		obj["cl_ord_id"] = rawJSONString(v)
	}
	if p.OrderUserRef != 0 {
		obj["order_userref"] = json.RawMessage(strconv.FormatInt(p.OrderUserRef, 10))
	}

	if v := strings.TrimSpace(p.DisplayQty); v != "" {
		raw, err := rawJSONNumberFromString(v)
		if err != nil {
			return nil, fmt.Errorf("display_qty: %w", err)
		}
		obj["display_qty"] = raw
	}

	if v := strings.TrimSpace(p.FeePreference); v != "" {
		obj["fee_preference"] = rawJSONString(v)
	}
	if p.NoMPP {
		obj["no_mpp"] = rawJSONBoolTrue()
	}
	if v := strings.TrimSpace(p.STPType); v != "" {
		// WS v2 docs use underscore separated stp_type values (cancel_newest, ...).
		obj["stp_type"] = rawJSONString(strings.ReplaceAll(v, "-", "_"))
	}

	if v := strings.TrimSpace(p.CashOrderQty); v != "" {
		raw, err := rawJSONNumberFromString(v)
		if err != nil {
			return nil, fmt.Errorf("cash_order_qty: %w", err)
		}
		obj["cash_order_qty"] = raw
	}

	if v := strings.TrimSpace(p.SenderSubID); v != "" {
		obj["sender_sub_id"] = rawJSONString(v)
	}

	if p.Triggers != nil {
		raw, err := json.Marshal(p.Triggers)
		if err != nil {
			return nil, fmt.Errorf("triggers: %w", err)
		}
		obj["triggers"] = raw
	}
	if p.Conditional != nil {
		raw, err := json.Marshal(p.Conditional)
		if err != nil {
			return nil, fmt.Errorf("conditional: %w", err)
		}
		obj["conditional"] = raw
	}

	if p.Validate {
		obj["validate"] = rawJSONBoolTrue()
	}

	return json.Marshal(obj)
}

// MarshalJSON ensures numeric fields are encoded as JSON numbers per Kraken WebSocket v2 docs.
func (p TriggerParams) MarshalJSON() ([]byte, error) {
	price, err := rawJSONNumberFromString(p.Price)
	if err != nil {
		return nil, fmt.Errorf("price: %w", err)
	}

	obj := map[string]json.RawMessage{
		"price": price,
	}

	if v := strings.TrimSpace(p.Reference); v != "" {
		obj["reference"] = rawJSONString(v)
	}
	if v := strings.TrimSpace(p.PriceType); v != "" {
		obj["price_type"] = rawJSONString(v)
	}

	return json.Marshal(obj)
}

// MarshalJSON ensures numeric fields are encoded as JSON numbers per Kraken WebSocket v2 docs.
func (p ConditionalParams) MarshalJSON() ([]byte, error) {
	orderType := strings.TrimSpace(p.OrderType)
	if orderType == "" {
		return nil, fmt.Errorf("order_type is required")
	}

	obj := map[string]json.RawMessage{
		"order_type": rawJSONString(orderType),
	}

	if v := strings.TrimSpace(p.LimitPrice); v != "" {
		raw, err := rawJSONNumberFromString(v)
		if err != nil {
			return nil, fmt.Errorf("limit_price: %w", err)
		}
		obj["limit_price"] = raw
	}
	if v := strings.TrimSpace(p.LimitPriceType); v != "" {
		obj["limit_price_type"] = rawJSONString(v)
	}
	if v := strings.TrimSpace(p.TriggerPrice); v != "" {
		raw, err := rawJSONNumberFromString(v)
		if err != nil {
			return nil, fmt.Errorf("trigger_price: %w", err)
		}
		obj["trigger_price"] = raw
	}
	if v := strings.TrimSpace(p.TriggerPriceType); v != "" {
		obj["trigger_price_type"] = rawJSONString(v)
	}

	return json.Marshal(obj)
}
