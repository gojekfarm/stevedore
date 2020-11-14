package provider

import (
	"fmt"
	"strconv"

	pkgConfig "github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/plugin"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/consul"
	"github.com/mitchellh/mapstructure"
)

func init() {
	defaultPlugins["consul"] = ClientPlugin{PluginImpl: ConsulConfigProvider{conf: config.NewConfig()}}
}

// ConsulConfigProvider is the struct represents the ConfigProvider for Consul
type ConsulConfigProvider struct {
	conf config.Config
}

var _ pkgConfig.Provider = ConsulConfigProvider{}
var _ plugin.ConfigInterface = ConsulConfigProvider{}

const (
	consulHostFlag        = "host"
	consulPortFlag        = "port"
	consulPrefixFlag      = "prefix"
	consulStripPrefixFlag = "strip-prefix"
)

// ConsulConfig represents struct for ConsulConfig
type ConsulConfig struct {
	Path []string
}

// Fetch configuration from ConsulProvider
func (p ConsulConfigProvider) Fetch(
	context map[string]string,
	data interface{},
) (map[string]interface{}, error) {
	host, ok := context[consulHostFlag]
	if !ok {
		return nil, fmt.Errorf("%s not set", consulHostFlag)
	}

	port, ok := context[consulPortFlag]
	if !ok {
		return nil, fmt.Errorf("%s not set", consulPortFlag)
	}

	prefix, ok := context[consulPrefixFlag]
	if !ok {
		return nil, fmt.Errorf("%s not set", consulPrefixFlag)
	}

	stripPrefix, ok := context[consulStripPrefixFlag]
	if !ok {
		return nil, fmt.Errorf("%s not set", consulStripPrefixFlag)
	}

	shouldStripPrefix, err := strconv.ParseBool(stripPrefix)
	if err != nil {
		return nil, fmt.Errorf("invalid value for %s: %v", consulStripPrefixFlag, err)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	var consulConfig ConsulConfig
	err = mapstructure.Decode(data, &consulConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid consul configs: %v", err)
	}

	consulSource := consul.NewSource(
		consul.WithAddress(address),
		consul.WithPrefix(prefix),
		consul.StripPrefix(shouldStripPrefix),
	)

	err = p.conf.Load(consulSource)
	if err != nil {
		return nil, fmt.Errorf("error loading from consul: %v", err)
	}
	finalConfigs := map[string]interface{}{}

	for _, path := range consulConfig.Path {

		err = p.conf.Get(path).Scan(&finalConfigs)
		if err != nil {
			return nil, fmt.Errorf("unable to decode configs under path %s: %v", path, err)
		}
	}
	return finalConfigs, nil
}

// Init the ConsulConfigProvider
func (ConsulConfigProvider) Init() error {
	return nil
}

// Version of the ConsulConfigProvider
func (ConsulConfigProvider) Version() (string, error) {
	return "v0.0.1", nil
}

// Flags for the Consul ConfigProvider
func (ConsulConfigProvider) Flags() ([]plugin.Flag, error) {
	return []plugin.Flag{
		{Name: consulHostFlag, Default: "http://127.0.0.1", Usage: "host for consul"},
		{Name: consulPortFlag, Default: "8500", Usage: "port for consul"},
		{Name: consulPrefixFlag, Default: "/", Usage: "prefix for contacting consul"},
		{Name: consulStripPrefixFlag, Default: "true", Usage: "strip-prefix indicates whether to remove the prefix from config entries, or leave it in place."},
	}, nil
}

// Type of the Provider. ConsulConfigProvider is of type Config
func (ConsulConfigProvider) Type() (plugin.Type, error) {
	return plugin.TypeConfig, nil
}

// Help returns the help message for the plugin
func (ConsulConfigProvider) Help() (string, error) {
	return "Help Consul!", nil
}

// Close the plugin
func (ConsulConfigProvider) Close() error {
	return nil
}
