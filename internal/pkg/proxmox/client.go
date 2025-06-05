package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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
				InsecureSkipVerify: true,
			},
		}
	}
	return &Client{
		baseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: &httpClient,
	}
}

func (c *Client) Version() (string, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/version")
	if err != nil {
		return "", fmt.Errorf("failed to construct URL: %w", err)
	}
	headers := http.Header{}
	headers.Set("Authorization", c.apiToken)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var versionResp VersionResponse
	if err := json.Unmarshal(body, &versionResp); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return versionResp.Data.Version, nil
}
