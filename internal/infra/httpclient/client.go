package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/antonhancharyk/crypto-knight-tg-bot/internal/ports"
)

type Client struct {
	base string
	http *http.Client
}

func New(base string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		base: base,
		http: &http.Client{Timeout: timeout},
	}
}

type reportResponse struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

func (c *Client) FetchReport(ctx context.Context, from, to string) (*ports.ReportResult, error) {
	url := fmt.Sprintf("%s/reports?from=%s&to=%s", c.base, from, to)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	var rr reportResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&rr); err != nil {
		return nil, fmt.Errorf("decode report json: %w", err)
	}
	return &ports.ReportResult{Income: rr.Income, Expense: rr.Expense}, nil
}
