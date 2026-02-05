package rest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/am-sokolov/go-kraken-sdk/types"
	"github.com/shopspring/decimal"
)

// PublicService provides access to public REST API endpoints.
// These endpoints do not require authentication.
type PublicService struct {
	client *Client
}

// NewPublicService creates a new public service.
func NewPublicService(client *Client) *PublicService {
	return &PublicService{client: client}
}

// GetServerTime returns the server's time.
//
// API: GET /0/public/Time
// Docs: https://docs.kraken.com/api/docs/rest-api/get-server-time
func (s *PublicService) GetServerTime(ctx context.Context) (*types.ServerTime, error) {
	resp, err := s.client.DoPublic(ctx, "/0/public/Time", nil)
	if err != nil {
		return nil, err
	}

	var result types.ServerTime
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetSystemStatus returns the current system status.
//
// API: GET /0/public/SystemStatus
// Docs: https://docs.kraken.com/api/docs/rest-api/get-system-status
func (s *PublicService) GetSystemStatus(ctx context.Context) (*types.SystemStatus, error) {
	resp, err := s.client.DoPublic(ctx, "/0/public/SystemStatus", nil)
	if err != nil {
		return nil, err
	}

	var result types.SystemStatus
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetAssetsOptions contains options for GetAssets.
type GetAssetsOptions struct {
	// Assets is a comma-separated list of assets to query.
	// If empty, all assets are returned.
	Assets []string

	// AClass is the asset class to query (e.g., "currency").
	AClass string
}

// GetAssets returns information about available assets.
//
// API: GET /0/public/Assets
// Docs: https://docs.kraken.com/api/docs/rest-api/get-asset-info
func (s *PublicService) GetAssets(ctx context.Context, opts *GetAssetsOptions) (map[string]types.Asset, error) {
	params := url.Values{}

	if opts != nil {
		if len(opts.Assets) > 0 {
			params.Set("asset", strings.Join(opts.Assets, ","))
		}
		if opts.AClass != "" {
			params.Set("aclass", opts.AClass)
		}
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/Assets", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.Asset
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetAssetPairsOptions contains options for GetAssetPairs.
type GetAssetPairsOptions struct {
	// Pairs is a list of asset pairs to query.
	// If empty, all pairs are returned.
	Pairs []string

	// Info specifies what info to return: "info", "leverage", "fees", "margin".
	Info string

	// AClassBase is the asset class for base component.
	AClassBase string

	// CountryCode filters pairs available in the specified country.
	CountryCode string
}

// GetAssetPairs returns tradable asset pairs.
//
// API: GET /0/public/AssetPairs
// Docs: https://docs.kraken.com/api/docs/rest-api/get-tradable-asset-pairs
func (s *PublicService) GetAssetPairs(ctx context.Context, opts *GetAssetPairsOptions) (map[string]types.AssetPair, error) {
	params := url.Values{}

	if opts != nil {
		if len(opts.Pairs) > 0 {
			params.Set("pair", strings.Join(opts.Pairs, ","))
		}
		if opts.Info != "" {
			params.Set("info", opts.Info)
		}
		if opts.AClassBase != "" {
			params.Set("aclass_base", opts.AClassBase)
		}
		if opts.CountryCode != "" {
			params.Set("country_code", opts.CountryCode)
		}
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/AssetPairs", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.AssetPair
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetTicker returns ticker information for the specified pairs.
//
// API: GET /0/public/Ticker
// Docs: https://docs.kraken.com/api/docs/rest-api/get-ticker-information
func (s *PublicService) GetTicker(ctx context.Context, pairs []string) (map[string]types.Ticker, error) {
	params := url.Values{}

	if len(pairs) > 0 {
		params.Set("pair", strings.Join(pairs, ","))
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/Ticker", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.Ticker
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetOHLCOptions contains options for GetOHLC.
type GetOHLCOptions struct {
	// Interval is the time frame interval in minutes.
	// Allowed values: 1, 5, 15, 30, 60, 240, 1440, 10080, 21600
	Interval types.OHLCInterval

	// Since is the Unix timestamp to return data since.
	Since int64
}

// GetOHLC returns OHLC data for the specified pair.
//
// API: GET /0/public/OHLC
// Docs: https://docs.kraken.com/api/docs/rest-api/get-ohlc-data
func (s *PublicService) GetOHLC(ctx context.Context, pair string, opts *GetOHLCOptions) (*types.OHLCResult, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if opts != nil {
		if opts.Interval != 0 {
			params.Set("interval", strconv.Itoa(int(opts.Interval)))
		}
		if opts.Since > 0 {
			params.Set("since", strconv.FormatInt(opts.Since, 10))
		}
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/OHLC", params)
	if err != nil {
		return nil, err
	}

	// The response has a special format: {"XXBTZUSD": [...], "last": 123}
	var rawResult map[string]interface{}
	if err := resp.Decode(&rawResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &types.OHLCResult{
		Data: make(map[string][]types.OHLC),
	}

	for key, value := range rawResult {
		if key == "last" {
			if v, ok := value.(float64); ok {
				result.Last = int64(v)
			}
			continue
		}

		// Parse OHLC array
		if arr, ok := value.([]interface{}); ok {
			ohlcData := make([]types.OHLC, 0, len(arr))
			for _, item := range arr {
				if ohlcArr, ok := item.([]interface{}); ok && len(ohlcArr) >= 8 {
					ohlc := types.OHLC{
						Time:   int64(ohlcArr[0].(float64)),
						Open:   parseDecimal(ohlcArr[1]),
						High:   parseDecimal(ohlcArr[2]),
						Low:    parseDecimal(ohlcArr[3]),
						Close:  parseDecimal(ohlcArr[4]),
						VWAP:   parseDecimal(ohlcArr[5]),
						Volume: parseDecimal(ohlcArr[6]),
						Count:  int(ohlcArr[7].(float64)),
					}
					ohlcData = append(ohlcData, ohlc)
				}
			}
			result.Data[key] = ohlcData
		}
	}

	return result, nil
}

// GetOrderBookOptions contains options for GetOrderBook.
type GetOrderBookOptions struct {
	// Count is the maximum number of asks/bids to return.
	// Default is 100, maximum is 500.
	Count int
}

// GetOrderBook returns the order book for the specified pair.
//
// API: GET /0/public/Depth
// Docs: https://docs.kraken.com/api/docs/rest-api/get-order-book
func (s *PublicService) GetOrderBook(ctx context.Context, pair string, opts *GetOrderBookOptions) (map[string]types.OrderBook, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if opts != nil && opts.Count > 0 {
		params.Set("count", strconv.Itoa(opts.Count))
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/Depth", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.OrderBook
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetRecentTradesOptions contains options for GetRecentTrades.
type GetRecentTradesOptions struct {
	// Since is the Unix timestamp to return trades since.
	Since int64

	// Count is the maximum number of trades to return.
	// Default is 1000.
	Count int
}

// GetRecentTrades returns recent trades for the specified pair.
//
// API: GET /0/public/Trades
// Docs: https://docs.kraken.com/api/docs/rest-api/get-recent-trades
func (s *PublicService) GetRecentTrades(ctx context.Context, pair string, opts *GetRecentTradesOptions) (*types.PublicTradesResult, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if opts != nil {
		if opts.Since > 0 {
			params.Set("since", strconv.FormatInt(opts.Since, 10))
		}
		if opts.Count > 0 {
			params.Set("count", strconv.Itoa(opts.Count))
		}
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/Trades", params)
	if err != nil {
		return nil, err
	}

	// The response has a special format: {"XXBTZUSD": [...], "last": "123"}
	var rawResult map[string]interface{}
	if err := resp.Decode(&rawResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &types.PublicTradesResult{
		Trades: make([]types.PublicTrade, 0),
	}

	for key, value := range rawResult {
		if key == "last" {
			if v, ok := value.(string); ok {
				result.Last = v
			}
			continue
		}

		// Parse trades array
		if arr, ok := value.([]interface{}); ok {
			for _, item := range arr {
				if tradeArr, ok := item.([]interface{}); ok && len(tradeArr) >= 7 {
					trade := types.PublicTrade{
						Price:     parseDecimal(tradeArr[0]),
						Volume:    parseDecimal(tradeArr[1]),
						Time:      tradeArr[2].(float64),
						Side:      tradeArr[3].(string),
						OrderType: tradeArr[4].(string),
						Misc:      tradeArr[5].(string),
					}
					if len(tradeArr) >= 7 {
						trade.TradeID = int64(tradeArr[6].(float64))
					}
					result.Trades = append(result.Trades, trade)
				}
			}
		}
	}

	return result, nil
}

// GetRecentSpreadsOptions contains options for GetRecentSpreads.
type GetRecentSpreadsOptions struct {
	// Since is the Unix timestamp to return spreads since.
	Since int64
}

// GetRecentSpreads returns recent spread data for the specified pair.
//
// API: GET /0/public/Spread
// Docs: https://docs.kraken.com/api/docs/rest-api/get-recent-spreads
func (s *PublicService) GetRecentSpreads(ctx context.Context, pair string, opts *GetRecentSpreadsOptions) (*types.SpreadResult, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if opts != nil && opts.Since > 0 {
		params.Set("since", strconv.FormatInt(opts.Since, 10))
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/Spread", params)
	if err != nil {
		return nil, err
	}

	// The response has a special format: {"XXBTZUSD": [...], "last": 123}
	var rawResult map[string]interface{}
	if err := resp.Decode(&rawResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &types.SpreadResult{
		Data: make(map[string][]types.Spread),
	}

	for key, value := range rawResult {
		if key == "last" {
			if v, ok := value.(float64); ok {
				result.Last = int64(v)
			}
			continue
		}

		// Parse spreads array
		if arr, ok := value.([]interface{}); ok {
			spreads := make([]types.Spread, 0, len(arr))
			for _, item := range arr {
				if spreadArr, ok := item.([]interface{}); ok && len(spreadArr) >= 3 {
					spread := types.Spread{
						Time: int64(spreadArr[0].(float64)),
						Bid:  parseDecimal(spreadArr[1]),
						Ask:  parseDecimal(spreadArr[2]),
					}
					spreads = append(spreads, spread)
				}
			}
			result.Data[key] = spreads
		}
	}

	return result, nil
}

// parseDecimal parses a decimal value from various input types.
func parseDecimal(v interface{}) decimal.Decimal {
	switch val := v.(type) {
	case string:
		d, _ := decimal.NewFromString(val)
		return d
	case float64:
		return decimal.NewFromFloat(val)
	default:
		return decimal.Zero
	}
}

// GetGroupedOrderBookOptions contains options for GetGroupedOrderBook.
type GetGroupedOrderBookOptions struct {
	// Depth is the number of price levels to return per side.
	// Valid values: 10, 25, 100, 250, 1000. Default is 10.
	Depth int

	// Grouping specifies how many tick levels per price level.
	// Valid values: 1, 5, 10, 25, 50, 100, 250, 500, 1000. Default is 1.
	Grouping int
}

// GetGroupedOrderBook returns the grouped order book for the specified pair.
// The GroupedBook endpoint aggregates volume over a specified tick range.
//
// API: GET /0/public/GroupedBook
// Docs: https://docs.kraken.com/api/docs/rest-api/get-grouped-order-book
func (s *PublicService) GetGroupedOrderBook(ctx context.Context, pair string, opts *GetGroupedOrderBookOptions) (*types.GroupedOrderBook, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if opts != nil {
		if opts.Depth > 0 {
			params.Set("depth", strconv.Itoa(opts.Depth))
		}
		if opts.Grouping > 0 {
			params.Set("grouping", strconv.Itoa(opts.Grouping))
		}
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/GroupedBook", params)
	if err != nil {
		return nil, err
	}

	var result types.GroupedOrderBook
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetPreTradeData returns pre-trade transparency data for the specified symbol.
// Returns the top 10 levels of the aggregated order book.
//
// API: GET /0/public/PreTrade
// Docs: https://docs.kraken.com/api/docs/rest-api/get-pre-trade
func (s *PublicService) GetPreTradeData(ctx context.Context, symbol string) (*types.PreTradeData, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	resp, err := s.client.DoPublic(ctx, "/0/public/PreTrade", params)
	if err != nil {
		return nil, err
	}

	var result types.PreTradeData
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetPostTradeOptions contains options for GetPostTradeData.
type GetPostTradeOptions struct {
	// Symbol filters results to a specific currency pair.
	Symbol string

	// FromTS filters results to trades after this ISO 8601 timestamp.
	FromTS string

	// ToTS filters results to trades before or at this ISO 8601 timestamp.
	ToTS string

	// Count is the maximum number of trades to return (1-1000, default 1000).
	Count int
}

// GetPostTradeData returns post-trade transparency data.
// If no filter parameters are specified, returns the last 1000 trades for all pairs.
//
// API: GET /0/public/PostTrade
// Docs: https://docs.kraken.com/api/docs/rest-api/get-post-trade
func (s *PublicService) GetPostTradeData(ctx context.Context, opts *GetPostTradeOptions) (*types.PostTradeResult, error) {
	params := url.Values{}

	if opts != nil {
		if opts.Symbol != "" {
			params.Set("symbol", opts.Symbol)
		}
		if opts.FromTS != "" {
			params.Set("from_ts", opts.FromTS)
		}
		if opts.ToTS != "" {
			params.Set("to_ts", opts.ToTS)
		}
		if opts.Count > 0 {
			params.Set("count", strconv.Itoa(opts.Count))
		}
	}

	resp, err := s.client.DoPublic(ctx, "/0/public/PostTrade", params)
	if err != nil {
		return nil, err
	}

	var result types.PostTradeResult
	if err := resp.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
