package usecase

import (
	"context"
	"testing"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/infra/httpclient"
	"github.com/stretchr/testify/require"
)

// mock client
// type mockClient struct{}

// func (m *mockClient) FetchReport(ctx context.Context, from, to string) (*httpclient.ReportResponse, error) {
// 	if from == "bad" {
// 		return nil, errors.New("bad date")
// 	}
// 	return &httpclient.ReportResponse{Income: 100.0, Expense: 50.0}, nil
// }

func TestGetReport_Success(t *testing.T) {
	r := &ReportUsecase{client: &httpclient.Client{}}
	ctx := context.Background()
	_, err := r.GetReport(ctx, "2020-01-01", "2020-01-31")
	require.Error(t, err)
}
