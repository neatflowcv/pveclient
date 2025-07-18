package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProxmoxURL  string
	Insecure    bool
	Username    string
	Realm       string
	AuthMethod  string
	TokenID     string // if authMethod is token
	TokenSecret string // if authMethod is token
	Password    string // if authMethod is password
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
		log.Println("Proceeding with system environment variables...")
	} else {
		log.Println("Loaded environment variables from .env file")
	}

	proxmoxURL := getEnv("PROXMOX_URL", "https://localhost:8006")
	insecure := getEnv("PROXMOX_INSECURE", "false") == "true"
	username := getEnv("PROXMOX_USERNAME", "")
	realm := getEnv("PROXMOX_REALM", "")
	authMethod := getEnv("PROXMOX_AUTH_METHOD", "token")
	tokenID := getEnv("PROXMOX_TOKEN_ID", "")
	tokenSecret := getEnv("PROXMOX_TOKEN_SECRET", "")
	password := getEnv("PROXMOX_PASSWORD", "")

	return &Config{
		ProxmoxURL:  proxmoxURL,
		Realm:       realm,
		AuthMethod:  authMethod,
		TokenID:     tokenID,
		TokenSecret: tokenSecret,
		Username:    username,
		Password:    password,
		Insecure:    insecure,
	}
}

func getEnv(key string, defaultValue string) string {
	ret := os.Getenv(key)
	if ret == "" {
		return defaultValue
	}

	return ret
}
