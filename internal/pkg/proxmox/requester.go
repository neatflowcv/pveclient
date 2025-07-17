package proxmox

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
)

type StatusCode int

const (
	StatusCodeUnknown StatusCode = 0
	StatusCodeOK      StatusCode = 200
)

type Requester struct {
	client *http.Client
}

func NewRequester(isInsecure bool) *Requester {
	var httpClient http.Client
	if isInsecure {
		httpClient.Transport = &http.Transport{ //nolint:exhaustruct
			TLSClientConfig: &tls.Config{ //nolint:exhaustruct
				InsecureSkipVerify: true, //nolint:gosec
			},
		}
	}

	return &Requester{
		client: &httpClient,
	}
}

// Call return status code, response body, error.
func (r *Requester) Call(req *Request) (StatusCode, []byte, error) {
	resp, err := r.client.Do(req.req)
	if err != nil {
		return StatusCodeUnknown, nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatusCodeUnknown, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return StatusCode(resp.StatusCode), content, nil
}
