package ports

import "context"

// ReportResult is the DTO returned by report providers (e.g. HTTP API).
type ReportResult struct {
	Income  float64
	Expense float64
}

// ReportFetcher fetches report data for a date range from an external source.
type ReportFetcher interface {
	FetchReport(ctx context.Context, from, to string) (*ReportResult, error)
}
