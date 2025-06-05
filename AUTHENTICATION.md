# Proxmox VE Authentication Methods

This client supports multiple authentication methods for accessing the Proxmox VE API.

## 1. API Token Authentication (Recommended)

API tokens are the recommended method for programmatic access as they don't require storing passwords and can have limited permissions.

### Creating an API Token

1. Log into your Proxmox VE web interface
2. Go to **Datacenter** > **Permissions** > **API Tokens**
3. Click **Add** to create a new token
4. Fill in the details:
   - **User**: Select the user (e.g., `root@pam`)
   - **Token ID**: Give your token a name (e.g., `mytoken`)
   - **Privilege Separation**: Uncheck if you want the token to have the same permissions as the user
5. Click **Add** and save the generated secret

### Using API Token

```go
// Using environment variables
client := proxmox.NewInsecureClientWithAPIToken(
    "https://your-proxmox-server:8006",
    os.Getenv("PROXMOX_API_TOKEN"),
)

// Or directly
client := proxmox.NewInsecureClientWithAPIToken(
    "https://your-proxmox-server:8006",
    "root@pam!mytoken=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
)
```

### Environment Variables
```bash
export PROXMOX_URL="https://your-proxmox-server:8006"
export PROXMOX_API_TOKEN="root@pam!mytoken=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

## 2. Username/Password Authentication

This method uses traditional username and password authentication.

### Using Username/Password

```go
// Using environment variables
client := proxmox.NewInsecureClientWithAuth(
    "https://your-proxmox-server:8006",
    os.Getenv("PROXMOX_USERNAME"),
    os.Getenv("PROXMOX_PASSWORD"),
)

// Or directly
client := proxmox.NewInsecureClientWithAuth(
    "https://your-proxmox-server:8006",
    "root@pam",
    "your-password",
)
```

### Environment Variables
```bash
export PROXMOX_URL="https://your-proxmox-server:8006"
export PROXMOX_USERNAME="root@pam"
export PROXMOX_PASSWORD="your-password"
```

## 3. Setting Authentication After Client Creation

You can also create a client first and then set authentication credentials:

```go
client := proxmox.NewInsecureClient("https://your-proxmox-server:8006")

// Set API token
client.SetAPIToken("root@pam!mytoken=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")

// Or set username/password
client.SetAuth("root@pam", "your-password")
```

## 4. Manual Login (Ticket-based Authentication)

When using username/password, you can manually trigger the login process:

```go
client := proxmox.NewInsecureClientWithAuth(
    "https://your-proxmox-server:8006",
    "root@pam",
    "your-password",
)

// Manually login to get ticket
err := client.Login()
if err != nil {
    log.Fatal(err)
}

// Now the client has a valid ticket for subsequent requests
```

## TLS Configuration

### Insecure Client (Skip TLS Verification)
Use this for self-signed certificates (common in lab environments):

```go
client := proxmox.NewInsecureClientWithAPIToken(url, token)
```

### Secure Client (Verify TLS)
Use this for production environments with valid certificates:

```go
client := proxmox.NewClientWithAPIToken(url, token)
```

### Configuring TLS After Creation
```go
client := proxmox.NewClient(url)
client.SetInsecureSkipTLS(true) // Skip TLS verification
```

## Common Authentication Realms

- `@pam`: System users (Linux PAM)
- `@pve`: Proxmox VE users
- `@ldap`: LDAP users (if configured)
- `@ad`: Active Directory users (if configured)

## Error Handling

Common authentication errors:

- **401 Unauthorized**: Invalid credentials or expired token
- **403 Forbidden**: Valid credentials but insufficient permissions
- **Connection errors**: Check URL and network connectivity

## Security Best Practices

1. **Use API tokens** instead of passwords for automated scripts
2. **Enable privilege separation** on API tokens when possible
3. **Use secure HTTPS connections** in production
4. **Rotate API tokens regularly**
5. **Store credentials securely** (environment variables, secret managers)
6. **Use least privilege principle** - grant only necessary permissions

## Example Usage

```bash
# Set environment variables
export PROXMOX_URL="https://192.168.1.100:8006"
export PROXMOX_API_TOKEN="root@pam!automation=12345678-1234-1234-1234-123456789abc"

# Run the client
./pveclient
```

The client will automatically detect and use the appropriate authentication method based on the environment variables provided. 