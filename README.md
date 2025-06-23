# Dehydrated API Metadata Plugin for Netscaler

A plugin for the Dehydrated API that provides metadata extraction and management capabilities for SSL/TLS certificates stored in Citrix Netscaler/ADC appliances.

## Overview

This plugin extends the Dehydrated API functionality by providing integration with Citrix Netscaler/ADC appliances for certificate management. It allows you to retrieve and manage SSL/TLS certificates stored in Netscaler instances, supporting multiple environments with different configurations.

## Features

- **Multi-environment support**: Configure and manage certificates across multiple Netscaler environments (dev, staging, prod, etc.)
- **Certificate retrieval**: Get all certificates or specific certificates by name
- **Prefix-based filtering**: Use environment-specific prefixes to organize and filter certificates
- **Secure authentication**: Supports username/password authentication with SSL verification options
- **Error handling**: Comprehensive error handling and reporting for connection and API issues
- **Integration ready**: Implements the Dehydrated API plugin interface for seamless integration

## Requirements

- Go 1.24 or later
- Citrix Netscaler/ADC appliance
- Network access to Netscaler instance
- Valid Netscaler credentials

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/schumann-it/dehydrated-api-metadata-plugin-netscaler.git
   cd dehydrated-api-metadata-plugin-netscaler
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the plugin:
   ```bash
   go build -o dehydrated-api-metadata-plugin-netscaler .
   ```

### From Releases

Download the latest release for your platform from the [Releases page](https://github.com/schumann-it/dehydrated-api-metadata-plugin-netscaler/releases).

## Configuration

The plugin supports configuration for multiple environments. Each environment requires the following configuration:

```json
{
  "environments": {
    "dev": {
      "endpoint": "https://netscaler-dev.example.com",
      "username": "admin",
      "password": "your-password",
      "prefix": "dev-",
      "sslVerify": false
    },
    "prod": {
      "endpoint": "https://netscaler-prod.example.com",
      "username": "admin",
      "password": "your-password",
      "prefix": "prod-",
      "sslVerify": true
    }
  }
}
```

### Configuration Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `endpoint` | Yes | The Netscaler API endpoint URL (e.g., `https://netscaler.example.com`) |
| `username` | Yes | Netscaler admin username |
| `password` | Yes | Netscaler admin password |
| `prefix` | No | Prefix for certificate names (e.g., `dev-`, `prod-`) |
| `sslVerify` | No | Whether to verify SSL certificates (default: `false`) |

## Usage

The plugin implements the Dehydrated API plugin interface and provides the following functionality:

### Plugin Methods

1. **Initialize**: Sets up the plugin with configuration for multiple environments
2. **GetMetadata**: Returns plugin metadata and capabilities
3. **Close**: Handles plugin cleanup and resource release

### Netscaler Client Methods

The plugin provides a Netscaler client with the following methods:

- `GetAllCertificates()`: Retrieves all certificates for the configured environment
- `GetCertificate(name)`: Retrieves a specific certificate by name

### Example Usage

```go
// Initialize the plugin with configuration
config := `{
  "environments": {
    "dev": {
      "endpoint": "https://netscaler-dev.example.com",
      "username": "admin",
      "password": "password",
      "prefix": "dev-"
    }
  }
}`

plugin := &NetscalerPlugin{}
err := plugin.Initialize(config)
if err != nil {
    log.Fatal(err)
}

// Get all certificates
certificates, err := plugin.GetAllCertificates("dev")
if err != nil {
    log.Fatal(err)
}

// Get a specific certificate
cert, err := plugin.GetCertificate("dev", "example.com")
if err != nil {
    log.Fatal(err)
}
```

## Testing

### Unit Tests

Run unit tests (excludes integration tests):
```bash
make test
# or
go test ./...
```

### Integration Tests

Integration tests verify that the Netscaler client methods work correctly with a real Netscaler instance. These tests are tagged with `integration` and are not run by default.

#### Prerequisites

- A running Netscaler instance (can be a test/development instance)
- Network access to the Netscaler instance
- Valid credentials

#### Environment Variables

Set the following environment variables to configure the integration tests:

```bash
export NETSCALER_ENDPOINT="https://your-netscaler-instance.com"
export NETSCALER_USERNAME="your-username"
export NETSCALER_PASSWORD="your-password"
export NETSCALER_PREFIX="test-"  # Optional, defaults to "test-"
export NETSCALER_SSL_VERIFY="false"  # Optional, defaults to "false"
```

#### Running Integration Tests

```bash
# Run only integration tests
make test-integration
# or
go test -v -tags=integration ./...

# Run both unit and integration tests
make test-all
```

#### Integration Test Coverage

The integration tests verify:
- Client creation and authentication
- `GetAllCertificates()` method functionality
- `GetCertificate()` method functionality
- Error handling for non-existent certificates
- Prefix handling in certificate names

#### Skipping Integration Tests

If you don't have access to a Netscaler instance, the integration tests will be skipped automatically. You can also explicitly exclude them:

```bash
# Run only unit tests (explicitly exclude integration)
go test -v ./... -tags="!integration"
```

## Development

The project structure is organized as follows:

```
.
├── main.go                    # Main plugin implementation
├── netscaler/                 # Netscaler client package
│   ├── client.go              # Netscaler client implementation
│   ├── client_test.go         # Unit tests for client
│   ├── config.go              # Configuration handling
│   ├── config_test.go         # Unit tests for config
│   └── integration_test.go    # Integration tests
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── .goreleaser.yml            # GoReleaser configuration
├── Makefile                   # Build and test automation
└── README.md                  # This file
```

### Building for Different Platforms

The project uses GoReleaser for building releases across multiple platforms:

```bash
# Build for current platform
go build

# Build for all platforms (requires GoReleaser)
goreleaser build --snapshot --clean
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new functionality
- Ensure all tests pass before submitting PRs
- Follow Go coding conventions
- Update documentation as needed

## Author

Jan Schumann

## Support

For issues and questions:

- Create an issue on [GitHub](https://github.com/schumann-it/dehydrated-api-metadata-plugin-netscaler/issues)
- Check the existing issues for similar problems
- Review the test examples for usage patterns

## Related Projects

- [Dehydrated API](https://github.com/schumann-it/dehydrated-api-go) - The core API
- [OpenSSL Plugin](https://github.com/Schumann-IT/dehydrated-api-metadata-plugin-openssl) - Similar plugin for OpenSSL integration 