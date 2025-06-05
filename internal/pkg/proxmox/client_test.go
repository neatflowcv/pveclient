//go:build cluster

package proxmox_test

import (
	"os"
	"testing"

	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
	"github.com/stretchr/testify/require"
)

func TestClient_Version(t *testing.T) {
	// 환경변수에서 Proxmox 서버 URL을 가져옵니다
	baseURL := os.Getenv("PROXMOX_URL")
	if baseURL == "" {
		t.Skip("PROXMOX_URL environment variable not set, skipping integration test")
	}

	// 클라이언트 생성
	client := proxmox.NewClient(baseURL)

	// Version 메서드 테스트
	version, err := client.Version()
	require.NoError(t, err, "Version() should not return an error")

	// 버전이 비어있지 않은지 확인
	require.NotEmpty(t, version, "Version() should not return empty string")

	// 버전 형식이 올바른지 확인 (예: "8.1" 또는 "7.4-3")
	require.Contains(t, version, ".", "Version() should contain a dot in the format")

	t.Logf("Proxmox version: %s", version)
}
