package netscaler

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config with all fields",
			input: map[string]any{
				"prefix":    "test-prefix-",
				"endpoint":  "https://netscaler.example.com",
				"username":  "admin",
				"password":  "secret",
				"sslVerify": true,
			},
			want: &Config{
				Prefix:    "test-prefix-",
				Endpoint:  "https://netscaler.example.com",
				Username:  "admin",
				Password:  "secret",
				SslVerify: true,
			},
			wantErr: false,
		},
		{
			name: "valid config with minimal fields",
			input: map[string]any{
				"endpoint": "https://netscaler.example.com",
				"username": "admin",
				"password": "secret",
			},
			want: &Config{
				Endpoint:  "https://netscaler.example.com",
				Username:  "admin",
				Password:  "secret",
				SslVerify: false,
			},
			wantErr: false,
		},
		{
			name: "valid config with empty prefix",
			input: map[string]any{
				"prefix":    "",
				"endpoint":  "https://netscaler.example.com",
				"username":  "admin",
				"password":  "secret",
				"sslVerify": false,
			},
			want: &Config{
				Prefix:    "",
				Endpoint:  "https://netscaler.example.com",
				Username:  "admin",
				Password:  "secret",
				SslVerify: false,
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			want:    &Config{},
			wantErr: false,
		},
		{
			name:    "empty map input",
			input:   map[string]any{},
			want:    &Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.Prefix != tt.want.Prefix {
				t.Errorf("NewConfig() Prefix = %v, want %v", got.Prefix, tt.want.Prefix)
			}
			if got.Endpoint != tt.want.Endpoint {
				t.Errorf("NewConfig() Endpoint = %v, want %v", got.Endpoint, tt.want.Endpoint)
			}
			if got.Username != tt.want.Username {
				t.Errorf("NewConfig() Username = %v, want %v", got.Username, tt.want.Username)
			}
			if got.Password != tt.want.Password {
				t.Errorf("NewConfig() Password = %v, want %v", got.Password, tt.want.Password)
			}
			if got.SslVerify != tt.want.SslVerify {
				t.Errorf("NewConfig() SslVerify = %v, want %v", got.SslVerify, tt.want.SslVerify)
			}
		})
	}
}

func TestNewConfig_InvalidJSON(t *testing.T) {
	// Test with a type that can't be marshaled to JSON
	invalidInput := make(chan int)

	_, err := NewConfig(invalidInput)
	if err == nil {
		t.Error("NewConfig() should return error for invalid input")
	}
}
