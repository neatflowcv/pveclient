package proxmox

type Config struct {
	insecureSkipTLS bool
}

type ConfigOption func(*Config)

func WithInsecure() ConfigOption {
	return func(c *Config) {
		c.insecureSkipTLS = true
	}
}
