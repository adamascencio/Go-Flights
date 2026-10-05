package googleflights

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://www.google.com/travel/flights"

type Client struct {
	HTTPClient *http.Client
	Options    URLOptions
}

func NewClient(options URLOptions) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
		Options:    options,
	}
}

func (c *Client) Fetch(ctx context.Context, query TFS) ([]byte, error) {
	u, err := BuildURL(query, c.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to build search URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/html")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flights: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch flights: HTTP %s", res.Status)
	}

	const maxBodySize = 10 << 20 // 10MB
	body, err := io.ReadAll(io.LimitReader(res.Body, maxBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if len(body) > maxBodySize {
		return nil, fmt.Errorf("response exceeds %d bytes", maxBodySize)
	}

	return body, nil
}
