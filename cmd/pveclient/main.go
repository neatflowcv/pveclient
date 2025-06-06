package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
)

type Config struct {
	proxmoxURL string
	apiToken   string
}

var ErrEnvNotSet = errors.New("environment variable is not set")

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Println("Proceeding with system environment variables...")
	} else {
		log.Println("Loaded environment variables from .env file")
	}

	proxmoxURL := os.Getenv("PROXMOX_URL")
	if proxmoxURL == "" {
		return nil, fmt.Errorf("PROXMOX_URL: %w", ErrEnvNotSet)
	}

	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("PROXMOX_API_TOKEN: %w", ErrEnvNotSet)
	}

	return &Config{
		proxmoxURL: proxmoxURL,
		apiToken:   apiToken,
	}, nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := proxmox.NewClient(config.proxmoxURL, config.apiToken, proxmox.WithInsecure())
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
