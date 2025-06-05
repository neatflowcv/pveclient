package proxmox

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

type VersionResponse struct {
	Data struct {
		Version string `json:"version"`
		Release string `json:"release"`
		RepoBit string `json:"repoid"`
	} `json:"data"`
}

func NewClient(baseURL string, apiToken string, opts ...ConfigOption) *Client {
	var config Config
	for _, opt := range opts {
		opt(&config)
	}

	var httpClient http.Client
	if config.insecureSkipTLS {
		httpClient.Transport = &http.Transport{ //nolint:exhaustruct
			TLSClientConfig: &tls.Config{ //nolint:exhaustruct
				InsecureSkipVerify: true, //nolint:gosec
			},
		}
	}

	return &Client{
		baseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: &httpClient,
	}
}

var ErrInvalidStatusCode = errors.New("invalid status code")

func (c *Client) Version(ctx context.Context) (string, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/version")
	if err != nil {
		return "", fmt.Errorf("failed to construct URL: %w", err)
	}

	headers := http.Header{}
	headers.Set("Authorization", c.apiToken)

	req, err := newGetRequest(ctx, endpoint, headers)
	if err != nil {
		return "", err
	}

	statusCode, content, err := c.call(req)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var versionResp VersionResponse
	if err := json.Unmarshal(content, &versionResp); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return versionResp.Data.Version, nil
}

func (c *Client) call(req *http.Request) (int, []byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp.StatusCode, content, nil
}

func newGetRequest(ctx context.Context, endpoint string, headers http.Header) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header = headers

	return req, nil
}
