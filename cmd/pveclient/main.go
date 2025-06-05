package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Println("Proceeding with system environment variables...")
	} else {
		log.Println("Loaded environment variables from .env file")
	}

	// Get Proxmox server URL from environment variable or use default
	proxmoxURL := os.Getenv("PROXMOX_URL")
	if proxmoxURL == "" {
		log.Fatal("PROXMOX_URL is not set")
	}
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken == "" {
		log.Fatal("PROXMOX_API_TOKEN is not set")
	}

	client := proxmox.NewInsecureClientWithAPIToken(proxmoxURL, apiToken)
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
