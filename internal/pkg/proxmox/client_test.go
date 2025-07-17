//go:build cluster

package proxmox_test

import (
	"os"
	"testing"

	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
	"github.com/stretchr/testify/require"
)

func TestClient_Version(t *testing.T) {
	baseURL := os.Getenv("PROXMOX_URL")
	if baseURL == "" {
		t.Skip("PROXMOX_URL environment variable not set, skipping integration test")
	}
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken == "" {
		t.Skip("PROXMOX_API_TOKEN environment variable not set, skipping integration test")
	}
	client := proxmox.NewClient(baseURL, apiToken, proxmox.WithInsecure())
	ctx := t.Context()

	version, err := client.Version(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, version)
	require.Contains(t, version, ".")
}

func TestClient_ListNodes(t *testing.T) {
	baseURL := os.Getenv("PROXMOX_URL")
	if baseURL == "" {
		t.Skip("PROXMOX_URL environment variable not set, skipping integration test")
	}
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken == "" {
		t.Skip("PROXMOX_API_TOKEN environment variable not set, skipping integration test")
	}
	client := proxmox.NewClient(baseURL, apiToken, proxmox.WithInsecure())
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
	baseURL := os.Getenv("PROXMOX_URL")
	if baseURL == "" {
		t.Skip("PROXMOX_URL environment variable not set, skipping integration test")
	}
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken == "" {
		t.Skip("PROXMOX_API_TOKEN environment variable not set, skipping integration test")
	}
	client := proxmox.NewClient(baseURL, apiToken, proxmox.WithInsecure())
	ctx := t.Context()
	nodes, _ := client.ListNodes(ctx)

	disks, err := client.ListDisks(ctx, nodes.Data[0].Node)

	require.NoError(t, err)
	require.NotEmpty(t, disks)
}
