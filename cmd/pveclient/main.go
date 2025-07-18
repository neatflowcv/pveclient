package main

import (
	"context"
	"fmt"
	"log"

	"github.com/neatflowcv/pveclient/internal/pkg/config"
	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
)

func main() {
	config := config.LoadConfig()
	client := newClient(config)
	testConnection(client)
}

func newClient(config *config.Config) *proxmox.Client {
	var opts []proxmox.ConfigOption

	if config.Insecure {
		opts = append(opts, proxmox.WithInsecure())
	}

	switch config.AuthMethod {
	case "token":
		secret := fmt.Sprintf("PVEAPIToken=%s@%s!%s=%s", config.Username, config.Realm, config.TokenID, config.TokenSecret)

		return proxmox.NewClient(config.ProxmoxURL, secret, opts...)

	case "password":
		panic("unimplemented")

	default:
		panic("invalid auth method: " + config.AuthMethod)
	}
}

func testConnection(client *proxmox.Client) {
	ctx := context.Background()

	version, err := client.Version(ctx)
	if err != nil {
		log.Printf("Failed to get Proxmox version: %v", err)

		return
	}

	log.Printf("Successfully connected to Proxmox VE version: %s\n", version)
}
