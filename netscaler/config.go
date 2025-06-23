package netscaler

import "encoding/json"

type Config struct {
	Prefix    string `json:"prefix,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	SslVerify bool   `json:"sslVerify,omitempty"`
}

func NewConfig(v any) (*Config, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var c = &Config{}
	err = json.Unmarshal(b, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
