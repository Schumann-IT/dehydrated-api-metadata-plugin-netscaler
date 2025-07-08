package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/schumann-it/dehydrated-api-go/plugin/proto"
	"github.com/schumann-it/dehydrated-api-go/plugin/server"

	"github.com/schumann-it/dehydrated-api-metadata-plugin-netscaler/netscaler"
)

type envConfig map[string]netscaler.Config

var (
	// These variables are set by GoReleaser during build
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// NetscalerPlugin is a simple plugin implementation
type NetscalerPlugin struct {
	proto.UnimplementedPluginServer
	logger        hclog.Logger
	config        *proto.PluginConfig
	clients       map[string]*netscaler.Client
	clientFactory func(prefix string, config *netscaler.ClientConfig) (*netscaler.Client, error)
}

// Initialize implements the plugin.Plugin interface
func (p *NetscalerPlugin) Initialize(_ context.Context, req *proto.InitializeRequest) (*proto.InitializeResponse, error) {
	p.logger.Debug("Initialize called")
	p.config.FromProto(req.Config)

	p.clients = make(map[string]*netscaler.Client)

	// Convert the protobuf Config to our envConfig type
	envConfigs := make(envConfig)
	environments, err := p.config.GetMap("environments")
	if err != nil {
		return nil, fmt.Errorf("invalid config format: %s", err.Error())
	}

	for env, value := range environments {
		p.logger.Debug("Creating client", "environment", env)
		cfg, err := netscaler.NewConfig(value)
		if err != nil {
			return nil, fmt.Errorf("invalid Config format for environment %s: %w", env, err)
		}

		p.logger.Debug("Config",
			"environment", env,
			"endpoint", cfg.Endpoint,
			"username", cfg.Username,
			"prefix", cfg.Prefix,
			"sslverify", cfg.SslVerify)

		// Validate required fields
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("missing required field 'endpoint' for environment %s", env)
		}
		if cfg.Username == "" {
			return nil, fmt.Errorf("missing required field 'username' for environment %s", env)
		}
		if cfg.Password == "" {
			return nil, fmt.Errorf("missing required field 'password' for environment %s", env)
		}

		envConfigs[env] = *cfg
	}

	// Create Netscaler clients for each environment
	for env, cfg := range envConfigs {
		p.logger.Debug("Creating Netscaler client", "environment", env)
		clientConfig := &netscaler.ClientConfig{
			Endpoint:  cfg.Endpoint,
			Username:  cfg.Username,
			Password:  cfg.Password,
			SslVerify: cfg.SslVerify,
			Headers:   make(map[string]string),
		}

		factory := p.clientFactory
		if factory == nil {
			factory = netscaler.NewClient
		}
		client, err := factory(cfg.Prefix, clientConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Netscaler client for environment %s: %w", env, err)
		}
		p.clients[env] = client
	}

	return &proto.InitializeResponse{}, nil
}

// GetMetadata implements the plugin.Plugin interface
func (p *NetscalerPlugin) GetMetadata(_ context.Context, req *proto.GetMetadataRequest) (*proto.GetMetadataResponse, error) {
	p.logger.Debug("GetMetadata called")

	// Create a new Metadata for the response
	metadata := proto.NewMetadata()

	for env, client := range p.clients {
		name := req.GetDomainEntry().GetDomain()
		if req.DomainEntry.GetAlias() != "" {
			name = req.DomainEntry.GetAlias()
		}
		cert, err := client.GetCertificate(name)
		if err != nil {
			cert = map[string]any{
				"error": fmt.Sprintf("failed to retrieve certificate for domain %s: %v", req.GetDomainEntry().GetDomain(), err),
			}
		}

		_ = metadata.SetMap(env, cert)
	}

	return metadata.ToGetMetadataResponse()
}

// Close implements the plugin.Plugin interface
func (p *NetscalerPlugin) Close(_ context.Context, _ *proto.CloseRequest) (*proto.CloseResponse, error) {
	p.logger.Debug("Close called")
	return &proto.CloseResponse{}, nil
}

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		printVersionInfoAndExit()
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "netscaler-plugin",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	plugin := &NetscalerPlugin{
		logger: logger,
		config: proto.NewPluginConfig(),
	}

	server.NewPluginServer(plugin).Serve()
}

// printVersionInfoAndExit prints the version information as a formatted string and exists the program
func printVersionInfoAndExit() {
	fmt.Printf("Version: %s\nCommit: %s\nBuild Time: %s\n", Version, Commit, BuildTime)
	os.Exit(0)
}
