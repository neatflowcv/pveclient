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
	Username  string
	Password  string
	APIToken  string
	Ticket    string
	CSRFToken string
}

type VersionResponse struct {
	Data struct {
		Version string `json:"version"`
		Release string `json:"release"`
		RepoBit string `json:"repoid"`
	} `json:"data"`
}

type AuthResponse struct {
	Data struct {
		Ticket              string `json:"ticket"`
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Username            string `json:"username"`
		Cap                 struct {
			Access map[string]int `json:"access"`
		} `json:"cap"`
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

// NewClientWithAuth creates a new client with username/password authentication
func NewClientWithAuth(baseURL, username, password string) *Client {
	client := NewClient(baseURL)
	client.Username = username
	client.Password = password
	return client
}

// NewInsecureClientWithAuth creates a new insecure client with username/password authentication
func NewInsecureClientWithAuth(baseURL, username, password string) *Client {
	client := NewInsecureClient(baseURL)
	client.Username = username
	client.Password = password
	return client
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

// SetAuth sets username and password for authentication
func (c *Client) SetAuth(username, password string) {
	c.Username = username
	c.Password = password
}

// SetAPIToken sets the API token for authentication
func (c *Client) SetAPIToken(apiToken string) {
	c.APIToken = apiToken
}

// Login authenticates with username/password and obtains a ticket
func (c *Client) Login() error {
	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("username and password are required for login")
	}

	endpoint, err := url.JoinPath(c.BaseURL, "/api2/json/access/ticket")
	if err != nil {
		return fmt.Errorf("failed to construct URL: %w", err)
	}

	// Prepare form data
	data := url.Values{}
	data.Set("username", c.Username)
	data.Set("password", c.Password)

	resp, err := c.HTTPClient.PostForm(endpoint, data)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	c.Ticket = authResp.Data.Ticket
	c.CSRFToken = authResp.Data.CSRFPreventionToken

	return nil
}

// addAuthHeaders adds appropriate authentication headers to the request
func (c *Client) addAuthHeaders(req *http.Request) {
	if c.APIToken != "" {
		// API Token authentication
		req.Header.Set("Authorization", "PVEAPIToken="+c.APIToken)
	} else if c.Ticket != "" {
		// Ticket-based authentication
		req.Header.Set("Cookie", "PVEAuthCookie="+c.Ticket)
		if c.CSRFToken != "" {
			req.Header.Set("CSRFPreventionToken", c.CSRFToken)
		}
	}
}

// makeAuthenticatedRequest creates and executes an authenticated request
func (c *Client) makeAuthenticatedRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	c.addAuthHeaders(req)

	// Set content type for POST requests
	if method == "POST" && body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

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

	// Try to authenticate first if we have credentials but no ticket/token
	if c.Username != "" && c.Password != "" && c.Ticket == "" && c.APIToken == "" {
		if err := c.Login(); err != nil {
			return "", fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	var resp *http.Response
	if c.APIToken != "" || c.Ticket != "" {
		// Use authenticated request
		resp, err = c.makeAuthenticatedRequest("GET", endpoint, nil)
	} else {
		// Try without authentication first (for public endpoints)
		resp, err = c.HTTPClient.Get(endpoint)
	}

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
