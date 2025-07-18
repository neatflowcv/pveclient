package main

import (
	"context"
	"log"

	"github.com/neatflowcv/pveclient/internal/pkg/config"
	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := proxmox.NewClient(config.ProxmoxURL, config.APIToken, proxmox.WithInsecure())
	testConnection(client)
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
