package proxmox

import (
	"context"
	"errors"
)

var _ Auth = (*LoginAuth)(nil)

type LoginAuth struct {
	realm    string
	username string
	password string
	response *IssueTicketResponse
}

var ErrInvalidLoginAuth = errors.New("realm, username, and password are required")

func NewLoginAuth(realm, username, password string) (*LoginAuth, error) {
	if realm == "" || username == "" || password == "" {
		return nil, ErrInvalidLoginAuth
	}

	return &LoginAuth{
		realm:    realm,
		username: username,
		password: password,
		response: nil,
	}, nil
}

func (a *LoginAuth) Authenticate(ctx context.Context, client *Client) error {
	response, err := client.IssueTicket(ctx, a.realm, a.username, a.password)
	if err != nil {
		return err
	}

	a.response = response

	return nil
}

func (a *LoginAuth) ModifyHeaders(headers map[string][]string) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["CSRFPreventionToken"] = []string{a.response.Data.CSRFPreventionToken}
	headers["Cookie"] = []string{"PVEAuthCookie=" + a.response.Data.Ticket}

	return headers
}
