package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/domain"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
)

type ReportUsecase struct {
	client *httpclient.Client
}

func NewReportUsecase(c *httpclient.Client) *ReportUsecase {
	return &ReportUsecase{client: c}
}

func (r *ReportUsecase) GetReport(ctx context.Context, from, to string) (*domain.Report, error) {
	if _, err := time.Parse("2006-01-02", from); err != nil {
		return nil, fmt.Errorf("invalid from date: %w", err)
	}
	if _, err := time.Parse("2006-01-02", to); err != nil {
		return nil, fmt.Errorf("invalid to date: %w", err)
	}
	resp, err := r.client.FetchReport(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &domain.Report{From: from, To: to, Income: resp.Income, Expense: resp.Expense}, nil
}
