package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nuonco/nuon-ext-api/internal/config"
	"github.com/nuonco/nuon-ext-api/internal/debug"
)

// Client makes authenticated HTTP requests to the Nuon API.
type Client struct {
	http    *http.Client
	baseURL string
	token   string
	orgID   string
}

// New creates a Client from the loaded config.
func New(cfg *config.Config) *Client {
	return &Client{
		http:    &http.Client{},
		baseURL: strings.TrimRight(cfg.APIURL, "/"),
		token:   cfg.APIToken,
		orgID:   cfg.OrgID,
	}
}

// Response holds the raw result of an API call.
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Do executes an HTTP request against the API.
func (c *Client) Do(method, path, payload string) (*Response, error) {
	url := c.baseURL + path

	var body io.Reader
	if payload != "" {
		body = strings.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	// Auth headers
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.orgID != "" {
		req.Header.Set("X-Nuon-Org-ID", c.orgID)
	}

	// Content headers
	if payload != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	debug.Log("http: %s %s", method, url)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	debug.Log("http: %d (%d bytes)", resp.StatusCode, len(respBody))

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}, nil
}
