package plugin

import (
	"fmt"
	"github.com/gojek/stevedore/pkg/stevedore"
	"gopkg.in/yaml.v2"
	"net/rpc"
)

// RPCClient represents a type that contains a RPCClient Client
type RPCClient struct{ client *rpc.Client }

// Version is the interface implementation for RPCClient Client
func (g *RPCClient) Version() (string, error) {
	var resp string
	rpcFuncCall := func() error { return g.client.Call("Plugin.Version", new(interface{}), &resp) }
	err := rpcCallWithTimeout(g.client, rpcFuncCall)
	return resp, err
}

// Help is the interface implementation for RPCClient Client
func (g *RPCClient) Help() (string, error) {
	var resp string
	rpcFuncCall := func() error { return g.client.Call("Plugin.Help", new(interface{}), &resp) }
	err := rpcCallWithTimeout(g.client, rpcFuncCall)
	return resp, err
}

// Type is the interface implementation for RPCClient Client
func (g *RPCClient) Type() (Type, error) {
	var resp Type
	rpcFuncCall := func() error { return g.client.Call("Plugin.Type", new(interface{}), &resp) }
	err := rpcCallWithTimeout(g.client, rpcFuncCall)
	return resp, err
}

// Flags fetch flags from the plugin
func (g *RPCClient) Flags() ([]Flag, error) {
	var resp []Flag
	rpcFuncCall := func() error { return g.client.Call("Plugin.Flags", new(interface{}), &resp) }
	err := rpcCallWithTimeout(g.client, rpcFuncCall)
	return resp, err
}

// Close the plugin
func (g *RPCClient) Close() error {
	return g.client.Close()
}

// Fetch retrieves data from the plugin
func (g *RPCClient) Fetch(
	context map[string]string,
	data interface{},
) (map[string]interface{}, error) {
	input := Input{data, context}
	pluginData, err := yaml.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("could not marshal configuration, err: %v", err)
	}
	var resp = make(map[string]interface{})
	rpcFuncCall := func() error { return g.client.Call("Plugin.Fetch", pluginData, &resp) }
	err = rpcCallWithTimeout(g.client, rpcFuncCall)
	return resp, err
}

// Manifests is the interface implementation for RPCClient Client
func (g *RPCClient) Manifests(data map[string]string) (stevedore.ManifestFiles, error) {
	var resp stevedore.ManifestFiles
	rpcFuncCall := func() error { return g.client.Call("Plugin.Manifests", data, &resp) }
	err := rpcCallWithTimeout(g.client, rpcFuncCall)
	return resp, err
}
