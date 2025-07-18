package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProxmoxURL string
	APIToken   string
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
		ProxmoxURL: proxmoxURL,
		APIToken:   apiToken,
	}, nil
}
