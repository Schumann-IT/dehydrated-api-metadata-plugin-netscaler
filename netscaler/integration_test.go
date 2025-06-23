//go:build integration
// +build integration

package netscaler

import (
	"os"
	"testing"
)

// Integration test configuration
type IntegrationTestConfig struct {
	Endpoint  string
	Username  string
	Password  string
	Prefix    string
	SslVerify bool
}

func getIntegrationTestConfig() *IntegrationTestConfig {
	return &IntegrationTestConfig{
		Endpoint:  os.Getenv("NETSCALER_ENDPOINT"),
		Username:  os.Getenv("NETSCALER_USERNAME"),
		Password:  os.Getenv("NETSCALER_PASSWORD"),
		Prefix:    os.Getenv("NETSCALER_PREFIX"),
		SslVerify: os.Getenv("NETSCALER_SSL_VERIFY") == "true",
	}
}

func TestIntegration_NewClient(t *testing.T) {
	config := getIntegrationTestConfig()

	clientConfig := &ClientConfig{
		Endpoint:  config.Endpoint,
		Username:  config.Username,
		Password:  config.Password,
		SslVerify: config.SslVerify,
		Headers:   make(map[string]string),
	}

	client, err := NewClient(config.Prefix, clientConfig)
	if err != nil {
		t.Skipf("Skipping integration test - failed to connect to Netscaler: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.prefix != config.Prefix {
		t.Errorf("Expected prefix %s, got %s", config.Prefix, client.prefix)
	}
}

func TestIntegration_GetAllCertificates(t *testing.T) {
	config := getIntegrationTestConfig()

	clientConfig := &ClientConfig{
		Endpoint:  config.Endpoint,
		Username:  config.Username,
		Password:  config.Password,
		SslVerify: config.SslVerify,
		Headers:   make(map[string]string),
	}

	client, err := NewClient(config.Prefix, clientConfig)
	if err != nil {
		t.Skipf("Skipping integration test - failed to connect to Netscaler: %v", err)
	}

	certificates, err := client.GetAllCertificates()
	if err != nil {
		t.Fatalf("Failed to get all certificates: %v", err)
	}

	// Verify that we got a valid response (even if empty)
	if certificates == nil {
		t.Fatal("Expected certificates slice, got nil")
	}

	t.Logf("Retrieved %d certificates from Netscaler", len(certificates))

	// If there are certificates, verify their structure
	for i, cert := range certificates {
		if cert == nil {
			t.Errorf("Certificate at index %d is nil", i)
			continue
		}

		// Check for required fields (adjust based on actual Netscaler response structure)
		if certkey, exists := cert["certkey"]; !exists {
			t.Errorf("Certificate at index %d missing 'certkey' field", i)
		} else {
			t.Logf("Certificate %d: certkey = %v", i, certkey)
		}
	}
}

func TestIntegration_GetCertificate(t *testing.T) {
	config := getIntegrationTestConfig()

	clientConfig := &ClientConfig{
		Endpoint:  config.Endpoint,
		Username:  config.Username,
		Password:  config.Password,
		SslVerify: config.SslVerify,
		Headers:   make(map[string]string),
	}

	client, err := NewClient(config.Prefix, clientConfig)
	if err != nil {
		t.Skipf("Skipping integration test - failed to connect to Netscaler: %v", err)
	}

	// First, get all certificates to find a valid certificate name
	certificates, err := client.GetAllCertificates()
	if err != nil {
		t.Skipf("Skipping certificate test - failed to get certificates: %v", err)
	}

	if len(certificates) == 0 {
		t.Skip("Skipping certificate test - no certificates found")
	}

	// Use the first certificate for testing
	firstCert := certificates[0]
	certkey, exists := firstCert["certkey"]
	if !exists {
		t.Skip("Skipping certificate test - no certkey found in first certificate")
	}

	certName, ok := certkey.(string)
	if !ok {
		t.Skip("Skipping certificate test - certkey is not a string")
	}

	// Remove the prefix if it exists to get the base name
	baseName := certName
	if len(config.Prefix) > 0 && len(certName) > len(config.Prefix) {
		baseName = certName[len(config.Prefix):]
	}

	t.Logf("Testing GetCertificate with name: %s (base: %s)", certName, baseName)

	// Test getting the specific certificate
	cert, err := client.GetCertificate(baseName)
	if err != nil {
		t.Fatalf("Failed to get certificate '%s': %v", baseName, err)
	}

	if cert == nil {
		t.Fatal("Expected certificate, got nil")
	}

	// Verify the certificate structure
	if retrievedCertkey, exists := cert["certkey"]; !exists {
		t.Error("Retrieved certificate missing 'certkey' field")
	} else {
		t.Logf("Retrieved certificate certkey: %v", retrievedCertkey)
	}

	// Test getting a non-existent certificate
	_, err = client.GetCertificate("non-existent-certificate-12345")
	if err == nil {
		t.Error("Expected error when getting non-existent certificate, got nil")
	} else {
		t.Logf("Correctly got error for non-existent certificate: %v", err)
	}
}

func TestIntegration_GetCertificate_WithPrefix(t *testing.T) {
	config := getIntegrationTestConfig()

	clientConfig := &ClientConfig{
		Endpoint:  config.Endpoint,
		Username:  config.Username,
		Password:  config.Password,
		SslVerify: config.SslVerify,
		Headers:   make(map[string]string),
	}

	client, err := NewClient(config.Prefix, clientConfig)
	if err != nil {
		t.Skipf("Skipping integration test - failed to connect to Netscaler: %v", err)
	}

	// Test with a sample certificate name
	testCertName := "test-certificate"

	// This should attempt to find a certificate with the prefix
	_, err = client.GetCertificate(testCertName)

	// We don't fail the test if the certificate doesn't exist
	// We just verify that the method works correctly
	if err != nil {
		t.Logf("GetCertificate with prefix returned error (expected if certificate doesn't exist): %v", err)
	} else {
		t.Logf("Successfully retrieved certificate: %s", testCertName)
	}
}

// Helper function to run integration tests only when explicitly requested
func TestIntegration_Setup(t *testing.T) {
	config := getIntegrationTestConfig()
	t.Logf("Integration test configuration:")
	t.Logf("  Endpoint: %s", config.Endpoint)
	t.Logf("  Username: %s", config.Username)
	t.Logf("  Prefix: %s", config.Prefix)
	t.Logf("  SSL Verify: %v", config.SslVerify)

	// Test basic connectivity
	clientConfig := &ClientConfig{
		Endpoint:  config.Endpoint,
		Username:  config.Username,
		Password:  config.Password,
		SslVerify: config.SslVerify,
		Headers:   make(map[string]string),
	}

	_, err := NewClient(config.Prefix, clientConfig)
	if err != nil {
		t.Logf("Warning: Cannot connect to Netscaler for integration tests: %v", err)
		t.Logf("Set NETSCALER_ENDPOINT, NETSCALER_USERNAME, and NETSCALER_PASSWORD environment variables to run integration tests")
	} else {
		t.Logf("Successfully connected to Netscaler")
	}
}
