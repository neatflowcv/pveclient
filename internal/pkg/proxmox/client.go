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
	BaseURL         string
	HTTPClient      *http.Client
	InsecureSkipTLS bool

	// Authentication fields
	APIToken string
}

type VersionResponse struct {
	Data struct {
		Version string `json:"version"`
		Release string `json:"release"`
		RepoBit string `json:"repoid"`
	} `json:"data"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:         baseURL,
		HTTPClient:      &http.Client{},
		InsecureSkipTLS: false,
	}
}

func NewInsecureClient(baseURL string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		BaseURL:         baseURL,
		HTTPClient:      &http.Client{Transport: tr},
		InsecureSkipTLS: true,
	}
}

// NewClientWithAPIToken creates a new client with API token authentication
func NewClientWithAPIToken(baseURL, apiToken string) *Client {
	client := NewClient(baseURL)
	client.APIToken = apiToken
	return client
}

// NewInsecureClientWithAPIToken creates a new insecure client with API token authentication
func NewInsecureClientWithAPIToken(baseURL, apiToken string) *Client {
	client := NewInsecureClient(baseURL)
	client.APIToken = apiToken
	return client
}

func (c *Client) SetInsecureSkipTLS(skip bool) {
	c.InsecureSkipTLS = skip
	if skip {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.HTTPClient.Transport = tr
	} else {
		c.HTTPClient.Transport = nil // Use default transport
	}
}

// SetAPIToken sets the API token for authentication
func (c *Client) SetAPIToken(apiToken string) {
	c.APIToken = apiToken
}

// addAuthHeaders adds appropriate authentication headers to the request
func (c *Client) addAuthHeaders(req *http.Request) {
	// API Token authentication
	req.Header.Set("Authorization", c.APIToken)
}

// makeAuthenticatedRequest creates and executes an authenticated request
func (c *Client) makeAuthenticatedRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeaders(req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make %s request: %w", method, err)
	}
	return resp, nil
}

func (c *Client) Version() (string, error) {
	// /api2/json/version
	endpoint, err := url.JoinPath(c.BaseURL, "/api2/json/version")
	if err != nil {
		return "", fmt.Errorf("failed to construct URL: %w", err)
	}

	var resp *http.Response
	resp, err = c.makeAuthenticatedRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
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
