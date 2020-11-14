package plugin

import (
	"net/rpc"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/hashicorp/go-plugin"
)

var _ ManifestInterface = &RPCClient{}
var _ plugin.Plugin = &ManifestPlugin{}

// ManifestPlugin is the plugin for creating manifest
type ManifestPlugin struct {
	Impl ManifestInterface
}

// Server is the implementation for plugin Server
func (p *ManifestPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ManifestRPCServer{Impl: p.Impl}, nil
}

// Client is the implementation for plugin Client
func (ManifestPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// ManifestRPCServer represents a type that contains an implementation for manifest interface
type ManifestRPCServer struct {
	Impl ManifestInterface
}

// Manifests implementation for the ManifestRPCServer Server
func (s *ManifestRPCServer) Manifests(args map[string]string, resp *stevedore.ManifestFiles) error {
	manifestFiles, err := s.Impl.Manifests(args)
	if err != nil {
		return err
	}
	*resp = manifestFiles
	return nil
}

// Flags implementation for the ManifestRPCServer Server
func (s *ManifestRPCServer) Flags(args interface{}, resp *[]Flag) error {
	return Flags(s.Impl, resp)
}

// Type implementation for the ManifestRPCServer Server
func (s *ManifestRPCServer) Type(args interface{}, resp *Type) error {
	return GetType(s.Impl, resp)
}

// Version implementation for the ManifestRPCServer Server
func (s *ManifestRPCServer) Version(args interface{}, resp *string) error {
	return Version(s.Impl, resp)
}

// Help implementation for the ManifestRPCServer Server
func (s *ManifestRPCServer) Help(args interface{}, resp *string) error {
	return Help(s.Impl, resp)
}

// ManifestProviderKey represents plugin key
const ManifestProviderKey = "manifest_provider"

// ServeManifestPlugin can be called by manifest plugin implementations as part of main
func ServeManifestPlugin(m ManifestInterface) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: CreateHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			ManifestProviderKey: &ManifestPlugin{Impl: m}},
	})
}
