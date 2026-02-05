package kraken

import (
	"context"
	"net/url"
	"strconv"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

// GetWebSocketsToken retrieves a WebSocket authentication token.
// The token should be used within 15 minutes of creation.
//
// API: POST /0/private/GetWebSocketsToken
// Docs: https://docs.kraken.com/api/docs/rest-api/get-websockets-token
func (s *AccountService) GetWebSocketsToken(ctx context.Context) (*types.WebSocketToken, error) {
	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/GetWebSocketsToken", url.Values{})
	if err != nil {
		return nil, err
	}

	var result types.WebSocketToken
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBalance retrieves all cash balances, net of pending withdrawals.
//
// API: POST /0/private/Balance
// Docs: https://docs.kraken.com/api/docs/rest-api/get-account-balance
func (s *AccountService) GetBalance(ctx context.Context) (map[string]string, error) {
	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Balance", url.Values{})
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTradeBalanceOptions contains options for GetTradeBalance.
type GetTradeBalanceOptions struct {
	// Asset is the base asset to determine balance.
	Asset string
}

// GetTradeBalance retrieves a summary of collateral balances and trading margins.
//
// API: POST /0/private/TradeBalance
// Docs: https://docs.kraken.com/api/docs/rest-api/get-trade-balance
func (s *AccountService) GetTradeBalance(ctx context.Context, opts *GetTradeBalanceOptions) (*types.TradeBalance, error) {
	params := url.Values{}
	if opts != nil && opts.Asset != "" {
		params.Set("asset", opts.Asset)
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/TradeBalance", params)
	if err != nil {
		return nil, err
	}

	var result types.TradeBalance
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOpenOrdersOptions contains options for GetOpenOrders.
type GetOpenOrdersOptions struct {
	// Trades includes trades in output.
	Trades bool
	// UserRef restricts results to given user reference.
	UserRef int64
	// ClientOrderID restricts results to given client order ID.
	ClientOrderID string
	// RebaseMultiplier is for viewing xstocks data ("rebased" or "base").
	RebaseMultiplier string
}

// GetOpenOrders retrieves information about currently open orders.
//
// API: POST /0/private/OpenOrders
// Docs: https://docs.kraken.com/api/docs/rest-api/get-open-orders
func (s *AccountService) GetOpenOrders(ctx context.Context, opts *GetOpenOrdersOptions) (map[string]types.Order, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Trades {
			params.Set("trades", "true")
		}
		if opts.UserRef != 0 {
			params.Set("userref", strconv.FormatInt(opts.UserRef, 10))
		}
		if opts.ClientOrderID != "" {
			params.Set("cl_ord_id", opts.ClientOrderID)
		}
		if opts.RebaseMultiplier != "" {
			params.Set("rebase_multiplier", opts.RebaseMultiplier)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/OpenOrders", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Open map[string]types.Order `json:"open"`
	}
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result.Open, nil
}

// GetClosedOrdersOptions contains options for GetClosedOrders.
type GetClosedOrdersOptions struct {
	// Trades includes trades in output.
	Trades bool
	// UserRef restricts results to given user reference.
	UserRef int64
	// ClientOrderID restricts results to given client order ID.
	ClientOrderID string
	// Start is the starting Unix timestamp or order tx ID.
	Start int64
	// End is the ending Unix timestamp or order tx ID.
	End int64
	// Offset is the result offset for pagination.
	Offset int
	// CloseTime is which time to use for filtering (open, close, both).
	CloseTime string
	// ConsolidateTaker consolidates trades by individual taker trades.
	ConsolidateTaker *bool
	// WithoutCount skips the count for faster performance.
	WithoutCount bool
	// RebaseMultiplier is for viewing xstocks data ("rebased" or "base").
	RebaseMultiplier string
}

// GetClosedOrders retrieves information about closed orders.
//
// API: POST /0/private/ClosedOrders
// Docs: https://docs.kraken.com/api/docs/rest-api/get-closed-orders
func (s *AccountService) GetClosedOrders(ctx context.Context, opts *GetClosedOrdersOptions) (map[string]types.Order, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Trades {
			params.Set("trades", "true")
		}
		if opts.Start != 0 {
			params.Set("start", strconv.FormatInt(opts.Start, 10))
		}
		if opts.End != 0 {
			params.Set("end", strconv.FormatInt(opts.End, 10))
		}
		if opts.Offset != 0 {
			params.Set("ofs", strconv.Itoa(opts.Offset))
		}
		if opts.UserRef != 0 {
			params.Set("userref", strconv.FormatInt(opts.UserRef, 10))
		}
		if opts.CloseTime != "" {
			params.Set("closetime", opts.CloseTime)
		}
		if opts.ClientOrderID != "" {
			params.Set("cl_ord_id", opts.ClientOrderID)
		}
		if opts.ConsolidateTaker != nil {
			if *opts.ConsolidateTaker {
				params.Set("consolidate_taker", "true")
			} else {
				params.Set("consolidate_taker", "false")
			}
		}
		if opts.WithoutCount {
			params.Set("without_count", "true")
		}
		if opts.RebaseMultiplier != "" {
			params.Set("rebase_multiplier", opts.RebaseMultiplier)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/ClosedOrders", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Closed map[string]types.Order `json:"closed"`
		Count  int                    `json:"count"`
	}
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result.Closed, nil
}

// QueryOrders queries information about specific orders.
//
// API: POST /0/private/QueryOrders
// Docs: https://docs.kraken.com/api/docs/rest-api/get-orders-info
func (s *AccountService) QueryOrders(ctx context.Context, txids []string, trades bool) (map[string]types.Order, error) {
	params := url.Values{}
	if len(txids) > 0 {
		for _, txid := range txids {
			params.Add("txid", txid)
		}
	}
	if trades {
		params.Set("trades", "true")
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/QueryOrders", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.Order
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTradesHistoryOptions contains options for GetTradesHistory.
type GetTradesHistoryOptions struct {
	// Type restricts results to given trade type.
	Type string
	// Trades includes trades related to position.
	Trades bool
	// Start is the starting Unix timestamp or trade tx ID.
	Start int64
	// End is the ending Unix timestamp or trade tx ID.
	End int64
	// Offset is the result offset for pagination.
	Offset int
	// WithoutCount skips the count for faster performance.
	WithoutCount bool
	// ConsolidateTaker consolidates trades by individual taker trades.
	ConsolidateTaker *bool
	// Ledgers includes related ledger IDs for each trade.
	Ledgers bool
	// RebaseMultiplier is for viewing xstocks data ("rebased" or "base").
	RebaseMultiplier string
}

// GetTradesHistory retrieves trade history.
//
// API: POST /0/private/TradesHistory
// Docs: https://docs.kraken.com/api/docs/rest-api/get-trade-history
func (s *AccountService) GetTradesHistory(ctx context.Context, opts *GetTradesHistoryOptions) (*types.TradesResult, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Type != "" {
			params.Set("type", opts.Type)
		}
		if opts.Trades {
			params.Set("trades", "true")
		}
		if opts.Start != 0 {
			params.Set("start", strconv.FormatInt(opts.Start, 10))
		}
		if opts.End != 0 {
			params.Set("end", strconv.FormatInt(opts.End, 10))
		}
		if opts.Offset != 0 {
			params.Set("ofs", strconv.Itoa(opts.Offset))
		}
		if opts.WithoutCount {
			params.Set("without_count", "true")
		}
		if opts.ConsolidateTaker != nil {
			if *opts.ConsolidateTaker {
				params.Set("consolidate_taker", "true")
			} else {
				params.Set("consolidate_taker", "false")
			}
		}
		if opts.Ledgers {
			params.Set("ledgers", "true")
		}
		if opts.RebaseMultiplier != "" {
			params.Set("rebase_multiplier", opts.RebaseMultiplier)
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/TradesHistory", params)
	if err != nil {
		return nil, err
	}

	var result types.TradesResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QueryTrades queries information about specific trades.
//
// API: POST /0/private/QueryTrades
// Docs: https://docs.kraken.com/api/docs/rest-api/get-trades-info
func (s *AccountService) QueryTrades(ctx context.Context, txids []string, trades bool) (map[string]types.Trade, error) {
	params := url.Values{}
	if len(txids) > 0 {
		for _, txid := range txids {
			params.Add("txid", txid)
		}
	}
	if trades {
		params.Set("trades", "true")
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/QueryTrades", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.Trade
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetOpenPositionsOptions contains options for GetOpenPositions.
type GetOpenPositionsOptions struct {
	// TxIDs restricts results to given transaction IDs.
	TxIDs []string
	// DoCalcs includes P&L calculations.
	DoCalcs bool
}

// GetOpenPositions retrieves information about open margin positions.
//
// API: POST /0/private/OpenPositions
// Docs: https://docs.kraken.com/api/docs/rest-api/get-open-positions
func (s *AccountService) GetOpenPositions(ctx context.Context, opts *GetOpenPositionsOptions) (map[string]types.Position, error) {
	params := url.Values{}
	if opts != nil {
		if len(opts.TxIDs) > 0 {
			for _, txid := range opts.TxIDs {
				params.Add("txid", txid)
			}
		}
		if opts.DoCalcs {
			params.Set("docalcs", "true")
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/OpenPositions", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.Position
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetLedgersOptions contains options for GetLedgers.
type GetLedgersOptions struct {
	// Asset restricts results to given asset.
	Asset string
	// AClass restricts results to given asset class.
	AClass string
	// Type restricts results to given type.
	Type string
	// Start is the starting Unix timestamp or ledger ID.
	Start int64
	// End is the ending Unix timestamp or ledger ID.
	End int64
	// Offset is the result offset for pagination.
	Offset int
}

// GetLedgers retrieves ledger entries.
//
// API: POST /0/private/Ledgers
// Docs: https://docs.kraken.com/api/docs/rest-api/get-ledgers
func (s *AccountService) GetLedgers(ctx context.Context, opts *GetLedgersOptions) (*types.LedgersResult, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Asset != "" {
			params.Set("asset", opts.Asset)
		}
		if opts.AClass != "" {
			params.Set("aclass", opts.AClass)
		}
		if opts.Type != "" {
			params.Set("type", opts.Type)
		}
		if opts.Start != 0 {
			params.Set("start", strconv.FormatInt(opts.Start, 10))
		}
		if opts.End != 0 {
			params.Set("end", strconv.FormatInt(opts.End, 10))
		}
		if opts.Offset != 0 {
			params.Set("ofs", strconv.Itoa(opts.Offset))
		}
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Ledgers", params)
	if err != nil {
		return nil, err
	}

	var result types.LedgersResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QueryLedgers queries information about specific ledger entries.
//
// API: POST /0/private/QueryLedgers
// Docs: https://docs.kraken.com/api/docs/rest-api/get-ledgers-info
func (s *AccountService) QueryLedgers(ctx context.Context, ids []string, trades bool) (map[string]types.Ledger, error) {
	params := url.Values{}
	if len(ids) > 0 {
		for _, id := range ids {
			params.Add("id", id)
		}
	}
	if trades {
		params.Set("trades", "true")
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/QueryLedgers", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.Ledger
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTradeVolumeOptions contains options for GetTradeVolume.
type GetTradeVolumeOptions struct {
	// Pair is a comma-separated list of asset pairs.
	Pair string
}

// GetTradeVolume retrieves trade volume and fee info.
//
// API: POST /0/private/TradeVolume
// Docs: https://docs.kraken.com/api/docs/rest-api/get-trade-volume
func (s *AccountService) GetTradeVolume(ctx context.Context, opts *GetTradeVolumeOptions) (*types.TradeVolume, error) {
	params := url.Values{}
	if opts != nil && opts.Pair != "" {
		params.Set("pair", opts.Pair)
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/TradeVolume", params)
	if err != nil {
		return nil, err
	}

	var result types.TradeVolume
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetL3OrderBookOptions contains options for GetL3OrderBook.
type GetL3OrderBookOptions struct {
	// Depth is the number of price levels per side.
	// Valid values: 0, 10, 25, 100, 250, 1000. Use 0 for full book.
	// Default is 100.
	Depth int
}

// GetL3OrderBook retrieves Level 3 order book data with individual order information.
// This endpoint requires authentication.
//
// API: POST /0/private/Level3
// Docs: https://docs.kraken.com/api/docs/rest-api/get-level-3-order-book
func (s *AccountService) GetL3OrderBook(ctx context.Context, pair string, opts *GetL3OrderBookOptions) (*types.L3OrderBook, error) {
	params := url.Values{}
	params.Set("pair", pair)

	if opts != nil {
		params.Set("depth", strconv.Itoa(opts.Depth))
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/Level3", params)
	if err != nil {
		return nil, err
	}

	var result types.L3OrderBook
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExtendedBalanceOptions contains options for GetExtendedBalance.
type GetExtendedBalanceOptions struct {
	// RebaseMultiplier is for viewing xstocks data ("rebased" or "base").
	RebaseMultiplier string
}

// GetExtendedBalance retrieves all extended account balances.
// Includes credits and held amounts. Available balance = balance + credit - credit_used - hold_trade
//
// API: POST /0/private/BalanceEx
// Docs: https://docs.kraken.com/api/docs/rest-api/get-extended-balance
func (s *AccountService) GetExtendedBalance(ctx context.Context, opts *GetExtendedBalanceOptions) (map[string]types.ExtendedBalance, error) {
	params := url.Values{}
	if opts != nil && opts.RebaseMultiplier != "" {
		params.Set("rebase_multiplier", opts.RebaseMultiplier)
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/BalanceEx", params)
	if err != nil {
		return nil, err
	}

	var result map[string]types.ExtendedBalance
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCreditLinesOptions contains options for GetCreditLines.
type GetCreditLinesOptions struct {
	// RebaseMultiplier is for viewing xstocks data ("rebased" or "base").
	RebaseMultiplier string
}

// GetCreditLines retrieves credit line details for VIP accounts.
//
// API: POST /0/private/CreditLines
// Docs: https://docs.kraken.com/api/docs/rest-api/get-credit-lines
func (s *AccountService) GetCreditLines(ctx context.Context, opts *GetCreditLinesOptions) (*types.CreditLines, error) {
	params := url.Values{}
	if opts != nil && opts.RebaseMultiplier != "" {
		params.Set("rebase_multiplier", opts.RebaseMultiplier)
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/CreditLines", params)
	if err != nil {
		return nil, err
	}

	var result types.CreditLines
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddExportRequest contains parameters for requesting an export report.
type AddExportRequest struct {
	// Report is the type of data to export ("trades" or "ledgers").
	Report types.ExportReportType

	// Description is a description for the export.
	Description string

	// Format is the file format ("CSV" or "TSV"). Default is CSV.
	Format types.ExportFormat

	// Fields is a comma-delimited list of fields to include.
	// For trades: ordertxid, time, ordertype, price, cost, fee, vol, margin, misc, ledgers
	// For ledgers: refid, time, type, subtype, aclass, asset, amount, fee, balance, wallet
	Fields string

	// StartTM is the UNIX timestamp for report start time.
	StartTM int64

	// EndTM is the UNIX timestamp for report end time.
	EndTM int64
}

// RequestExportReport requests an export of trades or ledgers.
//
// API: POST /0/private/AddExport
// Docs: https://docs.kraken.com/api/docs/rest-api/add-export
func (s *AccountService) RequestExportReport(ctx context.Context, req *AddExportRequest) (*types.ExportRequest, error) {
	params := url.Values{}
	params.Set("report", string(req.Report))
	params.Set("description", req.Description)

	if req.Format != "" {
		params.Set("format", string(req.Format))
	}
	if req.Fields != "" {
		params.Set("fields", req.Fields)
	}
	if req.StartTM > 0 {
		params.Set("starttm", strconv.FormatInt(req.StartTM, 10))
	}
	if req.EndTM > 0 {
		params.Set("endtm", strconv.FormatInt(req.EndTM, 10))
	}

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/AddExport", params)
	if err != nil {
		return nil, err
	}

	var result types.ExportRequest
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExportReportStatus retrieves the status of export reports.
//
// API: POST /0/private/ExportStatus
// Docs: https://docs.kraken.com/api/docs/rest-api/export-status
func (s *AccountService) GetExportReportStatus(ctx context.Context, reportType types.ExportReportType) ([]types.ExportReport, error) {
	params := url.Values{}
	params.Set("report", string(reportType))

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/ExportStatus", params)
	if err != nil {
		return nil, err
	}

	var result []types.ExportReport
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// RetrieveExportReport retrieves a processed data export as binary data.
// The returned data is a ZIP archive containing the report.
//
// API: POST /0/private/RetrieveExport
// Docs: https://docs.kraken.com/api/docs/rest-api/retrieve-export
func (s *AccountService) RetrieveExportReport(ctx context.Context, id string) ([]byte, error) {
	params := url.Values{}
	params.Set("id", id)

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/RetrieveExport", params)
	if err != nil {
		return nil, err
	}

	// The response is binary data (ZIP archive)
	return resp.RawData, nil
}

// RemoveExportType specifies whether to delete or cancel an export.
type RemoveExportType string

const (
	// RemoveExportDelete deletes a processed report.
	RemoveExportDelete RemoveExportType = "delete"

	// RemoveExportCancel cancels a queued or processing report.
	RemoveExportCancel RemoveExportType = "cancel"
)

// RemoveExportReport deletes or cancels an export report.
// Use "delete" for processed reports, "cancel" for queued/processing reports.
//
// API: POST /0/private/RemoveExport
// Docs: https://docs.kraken.com/api/docs/rest-api/remove-export
func (s *AccountService) RemoveExportReport(ctx context.Context, id string, removeType RemoveExportType) (*types.RemoveExportResult, error) {
	params := url.Values{}
	params.Set("id", id)
	params.Set("type", string(removeType))

	resp, err := s.client.restClient.DoPrivate(ctx, "/0/private/RemoveExport", params)
	if err != nil {
		return nil, err
	}

	var result types.RemoveExportResult
	if err := resp.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
