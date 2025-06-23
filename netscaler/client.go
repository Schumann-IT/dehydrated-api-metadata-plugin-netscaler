package netscaler

import (
	"fmt"

	"github.com/citrix/adc-nitro-go/service"
)

type Client struct {
	api    *service.NitroClient
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
		SslVerify: false,
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
	return c.api.FindAllResources(service.Sslcertkey.Type())
}

func (c *Client) GetCertificate(name string) (map[string]any, error) {
	return c.api.FindResource(service.Sslcertkey.Type(), fmt.Sprintf("%s%s", c.prefix, name))
}
