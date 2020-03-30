package plugin

import (
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v2"
	"net/rpc"
)

var _ ConfigInterface = &RPCClient{}
var _ plugin.Plugin = &ConfigPlugin{}

// ConfigPlugin is the type that contains an implementation for stevedore plugin interface
type ConfigPlugin struct {
	Impl ConfigInterface
}

// Server is the implementation for plugin Server
func (p *ConfigPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ConfigRPCServer{Impl: p.Impl}, nil
}

// Client is the implementation for plugin Client
func (ConfigPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// ConfigRPCServer represents a type that contains an implementation for stevedore plugin interface
type ConfigRPCServer struct {
	Impl ConfigInterface
}

// Input represents input to plugin
type Input struct {
	Data    interface{}       `yaml:"data"`
	Context map[string]string `yaml:"context"`
}

// Fetch is the interface implementation for ConfigRPCServer Server
func (s *ConfigRPCServer) Fetch(args []byte, resp *map[string]interface{}) error {
	input := Input{}
	err := yaml.Unmarshal(args, &input)
	if err != nil {
		return err
	}

	fetched, err := s.Impl.Fetch(input.Context, input.Data)
	if err != nil {
		return err
	}

	*resp = fetched
	return nil
}

// Flags implementation for the ConfigRPCServer Server
func (s *ConfigRPCServer) Flags(args interface{}, resp *[]Flag) error {
	return Flags(s.Impl, resp)
}

// Type implementation for the ConfigRPCServer Server
func (s *ConfigRPCServer) Type(args interface{}, resp *Type) error {
	return GetType(s.Impl, resp)
}

// Version implementation for the ConfigRPCServer Server
func (s *ConfigRPCServer) Version(args interface{}, resp *string) error {
	return Version(s.Impl, resp)
}

// Help implementation for the ConfigRPCServer Server
func (s *ConfigRPCServer) Help(args interface{}, resp *string) error {
	return Help(s.Impl, resp)
}

// ConfigProviderKey represents plugin key
const ConfigProviderKey = "config_provider"

// ServeConfigPlugin can be called by config plugin implementations as part of main
func ServeConfigPlugin(p ConfigInterface) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: CreateHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			ConfigProviderKey: &ConfigPlugin{Impl: p}},
	})
}
