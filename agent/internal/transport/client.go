package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/viniizn/Argus/agent/internal/credential"
)

const defaultTimeout = 10 * time.Second

// Client is the HTTP client used by the agent to communicate
// with the Argus server.
type Client struct {
	http      *http.Client
	serverURL string
	cred      credential.Secret
}

// New creates a new Client targeting serverURL.
//
// WARNING: TLS certificate verification is disabled for development.
// This must be replaced with proper TLS configuration before production use.
func New(serverURL string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // dev-only, tracked in #TODO
		},
	}

	return &Client{
		http: &http.Client{
			Timeout:   defaultTimeout,
			Transport: transport,
		},
		serverURL: serverURL,
	}
}

// SetCredential configures the bearer credential sent on subsequent requests.
func (c *Client) SetCredential(s credential.Secret) {
	c.cred = s
}

// Post sends a JSON-encoded body to the given path on the server.
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
	if c.cred != "" {
		req.Header.Set("Authorization", "Bearer "+c.cred.Reveal())
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport: execute request: %w", err)
	}

	return resp, nil
}