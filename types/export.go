package types

// ExportReportType represents the type of export report.
type ExportReportType string

const (
	// ExportReportTrades exports trade data.
	ExportReportTrades ExportReportType = "trades"

	// ExportReportLedgers exports ledger data.
	ExportReportLedgers ExportReportType = "ledgers"
)

// ExportFormat represents the file format for exports.
type ExportFormat string

const (
	// ExportFormatCSV exports as CSV.
	ExportFormatCSV ExportFormat = "CSV"

	// ExportFormatTSV exports as TSV.
	ExportFormatTSV ExportFormat = "TSV"
)

// ExportReportStatus represents the status of an export report.
type ExportReportStatus string

const (
	// ExportStatusQueued means the report is queued.
	ExportStatusQueued ExportReportStatus = "Queued"

	// ExportStatusProcessing means the report is being processed.
	ExportStatusProcessing ExportReportStatus = "Processing"

	// ExportStatusProcessed means the report is ready.
	ExportStatusProcessed ExportReportStatus = "Processed"
)

// ExportRequest represents the result of requesting an export.
type ExportRequest struct {
	// ID is the report ID.
	ID string `json:"id"`
}

// ExportReport represents an export report status.
type ExportReport struct {
	// ID is the report ID.
	ID string `json:"id"`

	// Descr is the report description.
	Descr string `json:"descr"`

	// Format is the file format (CSV/TSV).
	Format string `json:"format"`

	// Report is the report type (trades/ledgers).
	Report string `json:"report"`

	// Subtype is the report subtype.
	Subtype string `json:"subtype,omitempty"`

	// Status is the report status.
	Status ExportReportStatus `json:"status"`

	// Fields is the list of fields included.
	Fields string `json:"fields"`

	// CreatedTM is the UNIX timestamp of report request.
	CreatedTM string `json:"createdtm"`

	// StartTM is the UNIX timestamp when processing began.
	StartTM string `json:"starttm,omitempty"`

	// CompletedTM is the UNIX timestamp when processing finished.
	CompletedTM string `json:"completedtm,omitempty"`

	// DataStartTM is the UNIX timestamp of report data start.
	DataStartTM string `json:"datastarttm"`

	// DataEndTM is the UNIX timestamp of report data end.
	DataEndTM string `json:"dataendtm"`

	// Asset is the asset filter (if any).
	Asset string `json:"asset,omitempty"`
}

// RemoveExportResult represents the result of removing an export.
type RemoveExportResult struct {
	// Delete indicates if deletion was successful.
	Delete bool `json:"delete,omitempty"`

	// Cancel indicates if cancellation was successful.
	Cancel bool `json:"cancel,omitempty"`
}
