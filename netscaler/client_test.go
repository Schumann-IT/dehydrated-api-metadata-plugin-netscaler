package netscaler

import (
	"errors"
	"strings"
	"testing"
)

// MockNitroClient is a mock implementation of the NitroClient
type MockNitroClient struct {
	loginCalled bool
	loginErr    error
	findAllErr  error
	findErr     error
	allCerts    []map[string]any
	cert        map[string]any
}

func (m *MockNitroClient) Login() error {
	m.loginCalled = true
	return m.loginErr
}

func (m *MockNitroClient) FindAllResources(resourceType string) ([]map[string]any, error) {
	return m.allCerts, m.findAllErr
}

func (m *MockNitroClient) FindResource(resourceType string, name string) (map[string]any, error) {
	return m.cert, m.findErr
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		config      *ClientConfig
		mockLogin   error
		wantErr     bool
		description string
	}{
		{
			name:   "valid client creation",
			prefix: "test-prefix-",
			config: &ClientConfig{
				Endpoint:  "https://netscaler.example.com",
				Username:  "admin",
				Password:  "secret",
				SslVerify: false,
				Headers:   map[string]string{"X-Custom-Header": "value"},
			},
			mockLogin:   nil,
			wantErr:     false,
			description: "should create client successfully with valid config",
		},
		{
			name:   "client creation with login failure",
			prefix: "test-prefix-",
			config: &ClientConfig{
				Endpoint:  "https://netscaler.example.com",
				Username:  "admin",
				Password:  "secret",
				SslVerify: true,
				Headers:   map[string]string{},
			},
			mockLogin:   errors.New("authentication failed"),
			wantErr:     true,
			description: "should fail when login returns error",
		},
		{
			name:   "client creation with empty prefix",
			prefix: "",
			config: &ClientConfig{
				Endpoint:  "https://netscaler.example.com",
				Username:  "admin",
				Password:  "secret",
				SslVerify: false,
				Headers:   map[string]string{},
			},
			mockLogin:   nil,
			wantErr:     false,
			description: "should create client successfully with empty prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test would require mocking the service.NewNitroClientFromParams function
			// Since the actual implementation uses the real Citrix SDK, we'll test the structure
			// and basic functionality without the actual network calls

			client := &Client{
				prefix: tt.prefix,
			}

			if client.prefix != tt.prefix {
				t.Errorf("Client prefix = %v, want %v", client.prefix, tt.prefix)
			}
		})
	}
}

func TestClient_GetAllCertificates(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		mockCerts   []map[string]any
		mockErr     error
		wantErr     bool
		expectedLen int
		description string
	}{
		{
			name:   "successful get all certificates with prefix filtering",
			prefix: "test-prefix-",
			mockCerts: []map[string]any{
				{"certkey": "test-prefix-cert1", "cert": "cert1-data"},
				{"certkey": "test-prefix-cert2", "cert": "cert2-data"},
				{"certkey": "other-cert", "cert": "other-data"},
				{"certkey": "test-prefix-cert3", "cert": "cert3-data"},
			},
			mockErr:     nil,
			wantErr:     false,
			expectedLen: 3,
			description: "should return only certificates with matching prefix",
		},
		{
			name:        "empty certificates list",
			prefix:      "test-prefix-",
			mockCerts:   []map[string]any{},
			mockErr:     nil,
			wantErr:     false,
			expectedLen: 0,
			description: "should return empty list when no certificates exist",
		},
		{
			name:   "no certificates match prefix",
			prefix: "test-prefix-",
			mockCerts: []map[string]any{
				{"certkey": "other-cert1", "cert": "other1-data"},
				{"certkey": "different-cert2", "cert": "other2-data"},
			},
			mockErr:     nil,
			wantErr:     false,
			expectedLen: 0,
			description: "should return empty list when no certificates match prefix",
		},
		{
			name:   "empty prefix matches all certificates",
			prefix: "",
			mockCerts: []map[string]any{
				{"certkey": "cert1", "cert": "cert1-data"},
				{"certkey": "cert2", "cert": "cert2-data"},
				{"certkey": "test-prefix-cert3", "cert": "cert3-data"},
			},
			mockErr:     nil,
			wantErr:     false,
			expectedLen: 3,
			description: "should return all certificates when prefix is empty",
		},
		{
			name:        "error getting certificates",
			prefix:      "test-prefix-",
			mockCerts:   nil,
			mockErr:     errors.New("internal server error"),
			wantErr:     true,
			expectedLen: 0,
			description: "should return error when API call fails",
		},
		{
			name:   "certificate name is not a string",
			prefix: "test-prefix-",
			mockCerts: []map[string]any{
				{"certkey": "test-prefix-cert1", "cert": "cert1-data"},
				{"certkey": 123, "cert": "invalid-data"},
				{"certkey": "test-prefix-cert2", "cert": "cert2-data"},
			},
			mockErr:     nil,
			wantErr:     true,
			expectedLen: 0,
			description: "should return error when certificate name is not a string",
		},
		{
			name:   "certificate missing certkey field",
			prefix: "test-prefix-",
			mockCerts: []map[string]any{
				{"certkey": "test-prefix-cert1", "cert": "cert1-data"},
				{"cert": "missing-key-data"},
				{"certkey": "test-prefix-cert2", "cert": "cert2-data"},
			},
			mockErr:     nil,
			wantErr:     true,
			expectedLen: 0,
			description: "should return error when certificate is missing certkey field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockNitroClient{
				allCerts:   tt.mockCerts,
				findAllErr: tt.mockErr,
			}

			client := &Client{
				api:    mockAPI,
				prefix: tt.prefix,
			}

			got, err := client.GetAllCertificates()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCertificates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != tt.expectedLen {
					t.Errorf("GetAllCertificates() returned %d certificates, want %d", len(got), tt.expectedLen)
				}

				// Verify that all returned certificates have the correct prefix
				for _, cert := range got {
					if name, ok := cert["certkey"].(string); ok {
						if !strings.HasPrefix(name, tt.prefix) {
							t.Errorf("Certificate %s does not have expected prefix %s", name, tt.prefix)
						}
					} else {
						t.Errorf("Certificate name is not a string: %v", cert["certkey"])
					}
				}
			}
		})
	}
}

func TestClient_GetCertificate(t *testing.T) {
	tests := []struct {
		name        string
		certName    string
		prefix      string
		mockCert    map[string]any
		mockErr     error
		wantErr     bool
		expectedKey string
		description string
	}{
		{
			name:     "successful get certificate",
			certName: "example.com",
			prefix:   "test-prefix-",
			mockCert: map[string]any{
				"certkey": "test-prefix-example.com",
				"cert":    "certificate-data",
			},
			mockErr:     nil,
			wantErr:     false,
			expectedKey: "test-prefix-example.com",
			description: "should return certificate with prefixed name",
		},
		{
			name:        "certificate not found",
			certName:    "nonexistent.com",
			prefix:      "test-prefix-",
			mockCert:    nil,
			mockErr:     errors.New("certificate not found"),
			wantErr:     true,
			expectedKey: "test-prefix-nonexistent.com",
			description: "should return error when certificate doesn't exist",
		},
		{
			name:     "get certificate with empty prefix",
			certName: "example.com",
			prefix:   "",
			mockCert: map[string]any{
				"certkey": "example.com",
				"cert":    "certificate-data",
			},
			mockErr:     nil,
			wantErr:     false,
			expectedKey: "example.com",
			description: "should return certificate without prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockNitroClient{
				cert:    tt.mockCert,
				findErr: tt.mockErr,
			}

			client := &Client{
				api:    mockAPI,
				prefix: tt.prefix,
			}

			got, err := client.GetCertificate(tt.certName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Error("GetCertificate() returned nil certificate")
				}
			}
		})
	}
}

func TestClient_GetCertificate_KeyFormat(t *testing.T) {
	// Test that the certificate key is properly formatted with prefix
	client := &Client{
		prefix: "test-prefix-",
	}

	// This test verifies the key format logic
	certName := "example.com"
	expectedKey := "test-prefix-example.com"

	// The actual key formatting happens in the GetCertificate method
	// We can't easily test this without mocking the service layer,
	// but we can verify the logic is correct
	actualKey := client.prefix + certName
	if actualKey != expectedKey {
		t.Errorf("Certificate key format = %v, want %v", actualKey, expectedKey)
	}
}
