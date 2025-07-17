package proxmox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Client struct {
	baseURL   string
	apiToken  string
	requester *Requester
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

	return &Client{
		baseURL:   baseURL,
		apiToken:  apiToken,
		requester: NewRequester(config.insecureSkipTLS),
	}
}

var ErrInvalidStatusCode = errors.New("invalid status code")

func (c *Client) Version(ctx context.Context) (string, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/version")
	if err != nil {
		return "", fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewGetRequest(ctx, endpoint, map[string][]string{
		"Authorization": {c.apiToken},
	})

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return "", err
	}

	if statusCode != StatusCodeOK {
		return "", fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var versionResp VersionResponse

	err = json.Unmarshal(content, &versionResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return versionResp.Data.Version, nil
}

type ListNodesResponse struct {
	Data []struct {
		Status         string  `json:"status"`
		Maxcpu         int     `json:"maxcpu"`
		Mem            int64   `json:"mem"`
		CPU            float64 `json:"cpu"`
		Level          string  `json:"level"`
		Maxdisk        int64   `json:"maxdisk"`
		ID             string  `json:"id"`
		Maxmem         int64   `json:"maxmem"`
		Disk           int64   `json:"disk"`
		Type           string  `json:"type"`
		SslFingerprint string  `json:"ssl_fingerprint"`
		Node           string  `json:"node"`
		Uptime         int     `json:"uptime"`
	} `json:"data"`
}

func (c *Client) ListNodes(ctx context.Context) (*ListNodesResponse, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/nodes")
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewGetRequest(ctx, endpoint, map[string][]string{
		"Authorization": {c.apiToken},
	})

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return nil, err
	}

	if statusCode != StatusCodeOK {
		return nil, fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var ret ListNodesResponse

	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &ret, nil
}

type DiskType string

const (
	DiskTypeHDD DiskType = "hdd"
	DiskTypeSSD DiskType = "ssd"
)

type ListDisksResponse struct {
	Data []struct {
		Size         int64       `json:"size"`
		OsdidList    interface{} `json:"osdid-list"`
		Osdid        IntValue    `json:"osdid"`
		Used         string      `json:"used,omitempty"`
		Wwn          string      `json:"wwn"`
		Health       string      `json:"health"`
		Rpm          IntValue    `json:"rpm"`
		Gpt          int         `json:"gpt"`
		Type         DiskType    `json:"type"`
		ByIDLink     string      `json:"by_id_link"`
		Serial       string      `json:"serial"`
		Devpath      string      `json:"devpath"`
		Wearout      Wearout     `json:"wearout"`
		Model        string      `json:"model"`
		Vendor       string      `json:"vendor"`
		Db           int         `json:"db,omitempty"`
		Bluestore    int         `json:"bluestore,omitempty"`
		Osdencrypted int         `json:"osdencrypted,omitempty"`
	} `json:"data"`
}

func (c *Client) ListDisks(ctx context.Context, node string) (*ListDisksResponse, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/nodes", node, "disks/list")
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewGetRequest(ctx, endpoint, map[string][]string{
		"Authorization": {c.apiToken},
	})

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return nil, err
	}

	if statusCode != StatusCodeOK {
		return nil, fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var ret ListDisksResponse

	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &ret, nil
}

type Wearout struct {
	IsAvailable bool
	Value       int
}

func (w *Wearout) UnmarshalJSON(data []byte) error {
	if string(data) == "\"N/A\"" {
		w.IsAvailable = false
		w.Value = 0

		return nil
	}

	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse wearout: %w", err)
	}

	w.IsAvailable = true
	w.Value = int(value)

	return nil
}

type IntValue struct {
	Value int
}

func (i *IntValue) UnmarshalJSON(data []byte) error {
	// -1인 경우 처리
	if string(data) == "-1" {
		i.Value = -1

		return nil
	}

	// 따옴표로 둘러싸인 문자열인 경우 처리
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		// 따옴표 제거
		str := string(data[1 : len(data)-1])

		value, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse osdid from quoted string: %w", err)
		}

		i.Value = int(value)

		return nil
	}

	// 일반 숫자인 경우 처리
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse osdid: %w", err)
	}

	i.Value = int(value)

	return nil
}
