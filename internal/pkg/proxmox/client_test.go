//go:build cluster

package proxmox_test

import (
	"testing"

	"github.com/neatflowcv/pveclient/internal/pkg/config"
	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
	"github.com/stretchr/testify/require"
)

func TestClient_Version(t *testing.T) {
	t.Parallel()
	client := newClient(t)
	ctx := t.Context()

	version, err := client.Version(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, version)
	require.Contains(t, version, ".")
}

func TestClient_ListNodes(t *testing.T) {
	t.Parallel()
	client := newClient(t)
	ctx := t.Context()

	nodes, err := client.ListNodes(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, nodes)

	for _, node := range nodes.Data {
		require.Equal(t, "node", node.Type) // 고정값
		require.NotEmpty(t, node.Node)
	}
}

func TestClient_ListDisks(t *testing.T) {
	t.Parallel()
	client := newClient(t)
	ctx := t.Context()
	nodes, _ := client.ListNodes(ctx)

	disks, err := client.ListDisks(ctx, nodes.Data[0].Node)

	require.NoError(t, err)
	require.NotEmpty(t, disks)
}

func newClient(t *testing.T) *proxmox.Client {
	t.Helper()

	config, err := config.LoadConfig()
	require.NoError(t, err)

	client := proxmox.NewClient(config.ProxmoxURL, config.APIToken, proxmox.WithInsecure())

	return client
}
