package proxmox

import (
	"context"
	"errors"
	"fmt"
)

var _ Auth = (*TokenAuth)(nil)

type TokenAuth struct {
	username    string
	realm       string
	tokenID     string
	tokenSecret string
}

var ErrInvalidTokenAuth = errors.New("realm, username, tokenID, and tokenSecret are required")

func NewTokenAuth(realm, username, tokenID, tokenSecret string) (*TokenAuth, error) {
	if realm == "" || username == "" || tokenID == "" || tokenSecret == "" {
		return nil, ErrInvalidTokenAuth
	}

	return &TokenAuth{
		realm:       realm,
		username:    username,
		tokenID:     tokenID,
		tokenSecret: tokenSecret,
	}, nil
}

func (a *TokenAuth) Authenticate(ctx context.Context, client *Client) error {
	return nil
}

func (a *TokenAuth) ModifyHeaders(headers map[string][]string) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["Authorization"] = []string{
		fmt.Sprintf("PVEAPIToken=%s@%s!%s=%s", a.username, a.realm, a.tokenID, a.tokenSecret),
	}

	return headers
}
