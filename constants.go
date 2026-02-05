package kraken

// API base URLs
const (
	// DefaultBaseURL is the base URL for the Kraken REST API.
	DefaultBaseURL = "https://api.kraken.com"

	// APIVersion is the current API version.
	APIVersion = "0"

	// PublicWSURL is the WebSocket URL for public channels.
	PublicWSURL = "wss://ws.kraken.com/v2"

	// PrivateWSURL is the WebSocket URL for authenticated channels.
	PrivateWSURL = "wss://ws-auth.kraken.com/v2"
)

// API paths
const (
	// Public endpoints
	PathServerTime   = "/0/public/Time"
	PathSystemStatus = "/0/public/SystemStatus"
	PathAssets       = "/0/public/Assets"
	PathAssetPairs   = "/0/public/AssetPairs"
	PathTicker       = "/0/public/Ticker"
	PathOHLC         = "/0/public/OHLC"
	PathDepth        = "/0/public/Depth"
	PathTrades       = "/0/public/Trades"
	PathSpread       = "/0/public/Spread"

	// Private account endpoints
	PathBalance        = "/0/private/Balance"
	PathTradeBalance   = "/0/private/TradeBalance"
	PathOpenOrders     = "/0/private/OpenOrders"
	PathClosedOrders   = "/0/private/ClosedOrders"
	PathQueryOrders    = "/0/private/QueryOrders"
	PathTradesHistory  = "/0/private/TradesHistory"
	PathQueryTrades    = "/0/private/QueryTrades"
	PathOpenPositions  = "/0/private/OpenPositions"
	PathLedgers        = "/0/private/Ledgers"
	PathQueryLedgers   = "/0/private/QueryLedgers"
	PathTradeVolume    = "/0/private/TradeVolume"
	PathAddExport      = "/0/private/AddExport"
	PathExportStatus   = "/0/private/ExportStatus"
	PathRetrieveExport = "/0/private/RetrieveExport"
	PathRemoveExport   = "/0/private/RemoveExport"

	// Private trading endpoints
	PathAddOrder             = "/0/private/AddOrder"
	PathAddOrderBatch        = "/0/private/AddOrderBatch"
	PathEditOrder            = "/0/private/EditOrder"
	PathAmendOrder           = "/0/private/AmendOrder"
	PathCancelOrder          = "/0/private/CancelOrder"
	PathCancelOrderBatch     = "/0/private/CancelOrderBatch"
	PathCancelAll            = "/0/private/CancelAll"
	PathCancelAllOrdersAfter = "/0/private/CancelAllOrdersAfter"
	PathOrderAmends          = "/0/private/OrderAmends"

	// Private funding endpoints
	PathDepositMethods    = "/0/private/DepositMethods"
	PathDepositAddresses  = "/0/private/DepositAddresses"
	PathDepositStatus     = "/0/private/DepositStatus"
	PathWithdrawMethods   = "/0/private/WithdrawMethods"
	PathWithdrawAddresses = "/0/private/WithdrawAddresses"
	PathWithdrawInfo      = "/0/private/WithdrawInfo"
	PathWithdraw          = "/0/private/Withdraw"
	PathWithdrawStatus    = "/0/private/WithdrawStatus"
	PathWalletTransfer    = "/0/private/WalletTransfer"

	// Private earn endpoints
	PathEarnStrategies       = "/0/private/Earn/Strategies"
	PathEarnAllocations      = "/0/private/Earn/Allocations"
	PathEarnAllocate         = "/0/private/Earn/Allocate"
	PathEarnDeallocate       = "/0/private/Earn/Deallocate"
	PathEarnAllocateStatus   = "/0/private/Earn/AllocateStatus"
	PathEarnDeallocateStatus = "/0/private/Earn/DeallocateStatus"

	// Private subaccount endpoints
	PathCreateSubaccount = "/0/private/CreateSubaccount"

	// WebSocket token
	PathGetWebSocketsToken = "/0/private/GetWebSocketsToken"
)

// WebSocket channels
const (
	WSChannelTicker     = "ticker"
	WSChannelBook       = "book"
	WSChannelTrade      = "trade"
	WSChannelOHLC       = "ohlc"
	WSChannelInstrument = "instrument"
	WSChannelStatus     = "status"
	WSChannelHeartbeat  = "heartbeat"
	WSChannelExecutions = "executions"
	WSChannelBalances   = "balances"
)

// WebSocket methods
const (
	WSMethodPing        = "ping"
	WSMethodSubscribe   = "subscribe"
	WSMethodUnsubscribe = "unsubscribe"
	WSMethodAddOrder    = "add_order"
	WSMethodEditOrder   = "edit_order"
	WSMethodAmendOrder  = "amend_order"
	WSMethodCancelOrder = "cancel_order"
	WSMethodBatchAdd    = "batch_add"
	WSMethodBatchCancel = "batch_cancel"
	WSMethodCancelAll   = "cancel_all"
	WSMethodCancelAfter = "cancel_after"
)
