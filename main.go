package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	logger  hclog.Logger
	config  *proto.PluginConfig
	clients map[string]*netscaler.Client
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
		log.Printf("Invalid Config format for environment: %s", err.Error())
		return nil, fmt.Errorf("invalid Config format for environment: %s", err.Error())
	}

	for env, value := range environments {
		log.Printf("Processing environment %s", env)
		cfg, err := netscaler.NewConfig(value)
		if err != nil {
			log.Printf("Invalid Config format for environment %s: %s", env, err.Error())
			return nil, fmt.Errorf("invalid Config format for environment %s: %w", env, err)
		}

		log.Printf("Config for %s: endpoint=%s, username=%s, prefix=%s, sslVerify=%v",
			env, cfg.Endpoint, cfg.Username, cfg.Prefix, cfg.SslVerify)

		// Validate required fields
		if cfg.Endpoint == "" {
			log.Printf("Missing required field 'endpoint' for environment %s", env)
			return nil, fmt.Errorf("missing required field 'endpoint' for environment %s", env)
		}
		if cfg.Username == "" {
			log.Printf("Missing required field 'username' for environment %s", env)
			return nil, fmt.Errorf("missing required field 'username' for environment %s", env)
		}
		if cfg.Password == "" {
			log.Printf("Missing required field 'password' for environment %s", env)
			return nil, fmt.Errorf("missing required field 'password' for environment %s", env)
		}

		envConfigs[env] = *cfg
	}

	// Create Netscaler clients for each environment
	for env, cfg := range envConfigs {
		log.Printf("Creating Netscaler client for environment %s", env)
		clientConfig := &netscaler.ClientConfig{
			Endpoint:  cfg.Endpoint,
			Username:  cfg.Username,
			Password:  cfg.Password,
			SslVerify: cfg.SslVerify,
			Headers:   make(map[string]string),
		}

		client, err := netscaler.NewClient(cfg.Prefix, clientConfig)
		if err != nil {
			log.Printf("Failed to create Netscaler client for environment %s: %v", env, err)
			return nil, fmt.Errorf("failed to create Netscaler client for environment %s: %w", env, err)
		}
		log.Printf("Successfully created Netscaler client for environment %s", env)
		p.clients[env] = client
	}

	return &proto.InitializeResponse{}, nil
}

// GetMetadata implements the plugin.Plugin interface
func (p *NetscalerPlugin) GetMetadata(_ context.Context, req *proto.GetMetadataRequest) (*proto.GetMetadataResponse, error) {
	p.logger.Debug("GetMetadata called")

	// Create a new Metadata for the response
	metadata := proto.NewMetadata()

	var errs []string
	for env, client := range p.clients {
		m, err := client.GetCertificate(req.GetDomainEntry().GetDomain())
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to get metadata for environment %s: %s", env, err.Error()))
			continue
		}
		err = metadata.SetMap(env, m)
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to convert metadata for environment %s: %v", env, err.Error()))
		}
	}

	// Add errors to metadata if any occurred
	if len(errs) > 0 {
		metadata.SetError(strings.Join(errs, "; "))
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
		Name:   "netscaler-plugin",
		Level:  hclog.Trace,
		Output: os.Stdout,
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
