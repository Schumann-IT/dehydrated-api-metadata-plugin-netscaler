package netscaler

import (
	"fmt"
	"strings"

	"github.com/citrix/adc-nitro-go/service"
)

// NitroClientInterface defines the interface for NitroClient methods we use
type NitroClientInterface interface {
	Login() error
	FindAllResources(resourceType string) ([]map[string]any, error)
	FindResource(resourceType string, name string) (map[string]any, error)
}

type Client struct {
	api    NitroClientInterface
	prefix string
}

type ClientConfig struct {
	Endpoint  string
	Username  string
	Password  string
	SslVerify bool
	Headers   map[string]string
}

func NewClient(prefix string, config *ClientConfig) (*Client, error) {
	c := &Client{
		prefix: prefix,
	}

	api, err := service.NewNitroClientFromParams(service.NitroParams{
		Url:       config.Endpoint,
		Username:  config.Username,
		Password:  config.Password,
		SslVerify: config.SslVerify,
		Headers:   config.Headers,
	})
	if err != nil {
		return nil, err
	}

	c.api = api

	err = c.api.Login()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) GetAllCertificates() ([]map[string]any, error) {
	all, err := c.api.FindAllResources(service.Sslcertkey.Type())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve certificates: %w", err)
	}

	var certs = []map[string]any{}
	for _, cert := range all {
		if name, ok := cert["certkey"].(string); ok {
			if strings.HasPrefix(name, c.prefix) {
				certs = append(certs, cert)
			}
		} else {
			return nil, fmt.Errorf("certificate name is not a string: %v", cert)
		}
	}

	return certs, nil
}

func (c *Client) GetCertificate(name string) (map[string]any, error) {
	return c.api.FindResource(service.Sslcertkey.Type(), fmt.Sprintf("%s%s", c.prefix, name))
}
