// Package domain holds core business types.
package domain

// Report is a profit/loss summary for a date range.
type Report struct {
	From    string
	To      string
	Income  float64
	Expense float64
}
