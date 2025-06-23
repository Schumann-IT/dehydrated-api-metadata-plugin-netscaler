# Dehydrated API Metadata Plugin for Netscaler

A plugin for the Dehydrated API that provides metadata extraction and analysis capabilities for SSL/TLS certificates and private keys using OpenSSL.

## Overview

This plugin extends the Dehydrated API functionality by providing detailed metadata about SSL/TLS certificates and private keys. It analyzes certificate files and private keys to extract important information such as:

- Certificate metadata (subject, issuer, validity periods)
- Private key information (type, size)
- Certificate chain analysis

## Features

- Extracts metadata from various certificate files:
  - Private keys (`privkey.pem`)
  - Certificates (`cert.pem`)
  - Certificate chains (`chain.pem`)
  - Full certificate chains (`fullchain.pem`)
- Supports multiple key types:
  - RSA
  - ECDSA
  - Ed25519
- Provides detailed certificate information:
  - Subject DN
  - Issuer DN
  - Validity periods
  - Key type and size
- Error handling and reporting for invalid or corrupted files

## Requirements

- Go 1.x
- OpenSSL
- Dehydrated API Go client

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/schumann-it/dehydrated-api-metadata-plugin-openssl.git
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the plugin:
   ```bash
   go build
   ```

## Usage

The plugin implements the Dehydrated API plugin interface and provides the following functionality:

1. **Initialize**: Sets up the plugin with configuration
2. **GetMetadata**: Analyzes certificate files and returns metadata
3. **Close**: Handles plugin cleanup

The plugin processes the following files in the domain directory:
- `privkey.pem`: Private key file
- `cert.pem`: Certificate file
- `chain.pem`: Certificate chain file
- `fullchain.pem`: Full certificate chain file

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
├── main.go           # Main plugin implementation
├── netscaler/        # Netscaler client package
│   ├── client.go     # Netscaler client implementation
│   ├── client_test.go # Unit tests for client
│   ├── config.go     # Configuration handling
│   ├── config_test.go # Unit tests for config
│   └── integration_test.go # Integration tests
├── go.mod           # Go module definition
└── go.sum           # Go module checksums
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

[Add contribution guidelines here]

## Author

Jan Schumann 