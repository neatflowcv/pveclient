//go:build cluster

package proxmox_test

import (
	"context"
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
	ctx := context.Background()

	version, err := client.Version(ctx)

	require.NoError(t, err, "Version() should not return an error")
	require.NotEmpty(t, version, "Version() should not return empty string")
	require.Contains(t, version, ".", "Version() should contain a dot in the format")
}
