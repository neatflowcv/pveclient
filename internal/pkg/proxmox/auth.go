package proxmox

import "context"

type Auth interface {
	Authenticate(ctx context.Context, client *Client) error
	ModifyHeaders(headers map[string][]string) map[string][]string
}
