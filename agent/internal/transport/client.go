package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	http *http.Client
	serverURL string
}

// WARNING: TLS certificate verification is disabled for development.
// This must be replaced with proper TLS configuration before production use.
func New(serverURL string) *Client {
	transport := &http.Transport {
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // dev-only, tracked in #TODO
		},
	}

	return &Client{
		http: &http.Client{
			Timeout: defaultTimeout,
			Transport: transport,
		},
		serverURL: serverURL,
	}
}

func (c *Client) Post(ctx context.Context, path string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("transport: marshal request body: %w", err)
	}

	url := c.serverURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("transport: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport: execute request: %w", err)
	}

	return resp, nil
}