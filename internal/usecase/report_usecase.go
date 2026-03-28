// Package usecase contains application use cases orchestrating domain and ports.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/domain"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/ports"
)

// ReportUsecase loads reports through a ReportFetcher port.
type ReportUsecase struct {
	fetcher ports.ReportFetcher
}

// NewReportUsecase returns a use case backed by fetcher.
func NewReportUsecase(fetcher ports.ReportFetcher) *ReportUsecase {
	return &ReportUsecase{fetcher: fetcher}
}

// GetReport validates date strings and returns a domain report for the inclusive range.
func (r *ReportUsecase) GetReport(ctx context.Context, from, to string) (*domain.Report, error) {
	if _, err := time.Parse("2006-01-02", from); err != nil {
		return nil, fmt.Errorf("invalid from date: %w", err)
	}
	if _, err := time.Parse("2006-01-02", to); err != nil {
		return nil, fmt.Errorf("invalid to date: %w", err)
	}
	resp, err := r.fetcher.FetchReport(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("fetch report: %w", err)
	}
	return &domain.Report{From: from, To: to, Income: resp.Income, Expense: resp.Expense}, nil
}
