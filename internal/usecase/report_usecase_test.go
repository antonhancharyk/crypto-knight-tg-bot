package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/domain"
	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/ports"
	"github.com/stretchr/testify/require"
)

type mockReportFetcher struct {
	result *ports.ReportResult
	err    error
}

func (m *mockReportFetcher) FetchReport(_ context.Context, _, _ string) (*ports.ReportResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func TestGetReport(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		fetcher  *mockReportFetcher
		want     *domain.Report
		wantErr  bool
		errContains string
	}{
		{
			name: "success",
			from: "2020-01-01",
			to:   "2020-01-31",
			fetcher: &mockReportFetcher{
				result: &ports.ReportResult{Income: 100.5, Expense: 50.25},
			},
			want: &domain.Report{
				From: "2020-01-01", To: "2020-01-31",
				Income: 100.5, Expense: 50.25,
			},
		},
		{
			name:    "invalid from date",
			from:    "bad",
			to:      "2020-01-31",
			fetcher: &mockReportFetcher{},
			wantErr: true,
			errContains: "invalid from date",
		},
		{
			name:    "invalid to date",
			from:    "2020-01-01",
			to:      "not-a-date",
			fetcher: &mockReportFetcher{},
			wantErr: true,
			errContains: "invalid to date",
		},
		{
			name: "fetcher error",
			from: "2020-01-01",
			to:   "2020-01-31",
			fetcher: &mockReportFetcher{
				err: errors.New("api unavailable"),
			},
			wantErr: true,
			errContains: "api unavailable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewReportUsecase(tt.fetcher)
			got, err := uc.GetReport(context.Background(), tt.from, tt.to)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
