package main

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/schumann-it/dehydrated-api-go/plugin/proto"
	"github.com/schumann-it/dehydrated-api-metadata-plugin-netscaler/netscaler"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the netscaler client
type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetCertificate(domain string) (map[string]any, error) {
	args := m.Called(domain)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockClient) GetAllCertificates() ([]map[string]any, error) {
	args := m.Called()
	return args.Get(0).([]map[string]any), args.Error(1)
}

func mockClientFactory(prefix string, config *netscaler.ClientConfig) (*netscaler.Client, error) {
	return &netscaler.Client{}, nil
}

func TestNetscalerPlugin_Initialize(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]any
		wantErr     bool
		description string
	}{
		{
			name: "valid configuration",
			config: map[string]any{
				"environments": map[string]any{
					"prod": map[string]any{
						"endpoint":  "https://netscaler-prod.example.com",
						"username":  "admin",
						"password":  "secret",
						"prefix":    "prod-",
						"sslVerify": true,
					},
					"dev": map[string]any{
						"endpoint":  "https://netscaler-dev.example.com",
						"username":  "admin",
						"password":  "secret",
						"prefix":    "dev-",
						"sslVerify": false,
					},
				},
			},
			wantErr:     false,
			description: "should initialize successfully with valid config",
		},
		{
			name: "missing endpoint",
			config: map[string]any{
				"environments": map[string]any{
					"prod": map[string]any{
						"username": "admin",
						"password": "secret",
					},
				},
			},
			wantErr:     true,
			description: "should fail when endpoint is missing",
		},
		{
			name: "missing username",
			config: map[string]any{
				"environments": map[string]any{
					"prod": map[string]any{
						"endpoint": "https://netscaler.example.com",
						"password": "secret",
					},
				},
			},
			wantErr:     true,
			description: "should fail when username is missing",
		},
		{
			name: "missing password",
			config: map[string]any{
				"environments": map[string]any{
					"prod": map[string]any{
						"endpoint": "https://netscaler.example.com",
						"username": "admin",
					},
				},
			},
			wantErr:     true,
			description: "should fail when password is missing",
		},
		{
			name: "invalid environments format",
			config: map[string]any{
				"environments": "not a map",
			},
			wantErr:     true,
			description: "should fail when environments is not a map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "test",
				Level:  hclog.Trace,
				Output: hclog.DefaultOutput,
			})

			plugin := &NetscalerPlugin{
				logger:        logger,
				config:        proto.NewPluginConfig(),
				clientFactory: mockClientFactory,
			}

			// Set the config
			for k, v := range tt.config {
				plugin.config.Set(k, v)
			}

			protoConfig, err := plugin.config.ToProto()
			if err != nil {
				t.Fatalf("Failed to convert config to proto: %v", err)
			}

			req := &proto.InitializeRequest{
				Config: protoConfig,
			}

			_, err = plugin.Initialize(context.Background(), req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify that clients were created
				if len(plugin.clients) == 0 {
					t.Error("Initialize() should create clients when successful")
				}
			}
		})
	}
}

func TestNetscalerPlugin_GetMetadata(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Trace,
		Output: hclog.DefaultOutput,
	})

	plugin := &NetscalerPlugin{
		logger:  logger,
		config:  proto.NewPluginConfig(),
		clients: make(map[string]*netscaler.Client),
	}

	req := &proto.GetMetadataRequest{
		DomainEntry: &proto.DomainEntry{
			Domain: "example.com",
		},
	}

	resp, err := plugin.GetMetadata(context.Background(), req)
	if err != nil {
		t.Errorf("GetMetadata() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("GetMetadata() should return a response")
	}
}

func TestNetscalerPlugin_Close(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Trace,
		Output: hclog.DefaultOutput,
	})

	plugin := &NetscalerPlugin{
		logger: logger,
		config: proto.NewPluginConfig(),
	}

	req := &proto.CloseRequest{}
	resp, err := plugin.Close(context.Background(), req)

	if err != nil {
		t.Errorf("Close() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("Close() should return a response")
	}
}
