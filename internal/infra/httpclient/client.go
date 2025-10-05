package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

type ReportResponse struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

func (c *Client) FetchReport(ctx context.Context, from, to string) (*ReportResponse, error) {
	url := fmt.Sprintf("%s/reports?from=%s&to=%s", c.base, from, to)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	var rr ReportResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&rr); err != nil {
		return nil, err
	}
	return &rr, nil
}
