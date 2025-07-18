package proxmox

import (
	"bytes"
	"context"
	"net/http"
)

type Request struct {
	req *http.Request
}

func NewGetRequest(ctx context.Context, endpoint string, headers map[string][]string) *Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		panic(err)
	}

	req.Header = headers

	return &Request{
		req: req,
	}
}

func NewPostRequest(ctx context.Context, endpoint string, headers map[string][]string, data []byte) *Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	req.Header = headers

	return &Request{
		req: req,
	}
}
