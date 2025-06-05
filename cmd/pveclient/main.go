package main

import (
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
		return nil, fmt.Errorf("PROXMOX_URL is not set")
	}
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("PROXMOX_API_TOKEN is not set")
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
	client := proxmox.NewInsecureClientWithAPIToken(config.proxmoxURL, config.apiToken)
	testConnection(client)
}

func testConnection(client *proxmox.Client) {
	// Test the connection by getting the version
	version, err := client.Version()
	if err != nil {
		log.Printf("Failed to get Proxmox version: %v", err)
		return
	}

	fmt.Printf("Successfully connected to Proxmox VE version: %s\n", version)
}
