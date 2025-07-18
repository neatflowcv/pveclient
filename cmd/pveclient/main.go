package main

import (
	"context"
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
		log.Println("Using token authentication")

		auth, err := proxmox.NewTokenAuth(config.Realm, config.Username, config.TokenID, config.TokenSecret)
		if err != nil {
			panic(err)
		}

		client, err := proxmox.NewClient(context.Background(), config.ProxmoxURL, auth, opts...)
		if err != nil {
			panic(err)
		}

		return client

	case "password":
		log.Println("Using password authentication")

		auth, err := proxmox.NewLoginAuth(config.Realm, config.Username, config.Password)
		if err != nil {
			panic(err)
		}

		client, err := proxmox.NewClient(context.Background(), config.ProxmoxURL, auth, opts...)
		if err != nil {
			panic(err)
		}

		return client

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
