package main

import (
	"fmt"
	"log"
	"os"

	"github.com/neatflowcv/pveclient/internal/pkg/proxmox"
)

func main() {
	// Get Proxmox server URL from environment variable or use default
	proxmoxURL := os.Getenv("PROXMOX_URL")
	if proxmoxURL == "" {
		proxmoxURL = "https://your-proxmox-server:8006"
	}

	// Method 1: Using API Token (Recommended for scripts)
	apiToken := os.Getenv("PROXMOX_API_TOKEN")
	if apiToken != "" {
		fmt.Println("Using API Token authentication...")
		client := proxmox.NewInsecureClientWithAPIToken(proxmoxURL, apiToken)
		testConnection(client)
		return
	}

	// Method 2: Using Username/Password
	username := os.Getenv("PROXMOX_USERNAME")
	password := os.Getenv("PROXMOX_PASSWORD")
	if username != "" && password != "" {
		fmt.Println("Using username/password authentication...")
		client := proxmox.NewInsecureClientWithAuth(proxmoxURL, username, password)
		testConnection(client)
		return
	}

	// Method 3: Creating client and setting auth later
	fmt.Println("No authentication credentials found in environment variables.")
	fmt.Println("Please set one of the following:")
	fmt.Println("  - PROXMOX_API_TOKEN: API token in format 'user@realm!tokenid=secret'")
	fmt.Println("  - PROXMOX_USERNAME and PROXMOX_PASSWORD: Username and password")
	fmt.Println()
	fmt.Println("Example API token usage:")
	fmt.Println("  export PROXMOX_URL='https://your-proxmox-server:8006'")
	fmt.Println("  export PROXMOX_API_TOKEN='root@pam!mytoken=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'")
	fmt.Println()
	fmt.Println("Example username/password usage:")
	fmt.Println("  export PROXMOX_URL='https://your-proxmox-server:8006'")
	fmt.Println("  export PROXMOX_USERNAME='root@pam'")
	fmt.Println("  export PROXMOX_PASSWORD='your-password'")

	// Try without authentication (will likely fail with 401)
	fmt.Println("\nTrying without authentication (this will likely fail)...")
	client := proxmox.NewInsecureClient(proxmoxURL)
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
