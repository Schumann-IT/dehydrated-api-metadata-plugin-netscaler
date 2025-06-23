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

## Development

The project structure is organized as follows:

```
.
├── main.go           # Main plugin implementation
├── internal/         # Internal package
│   ├── certificate.go # Certificate analysis
│   └── key.go        # Private key analysis
├── go.mod           # Go module definition
└── go.sum           # Go module checksums
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

[Add contribution guidelines here]

## Author

Jan Schumann 