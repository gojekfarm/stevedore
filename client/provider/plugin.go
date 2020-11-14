package provider

import (
	"fmt"

	"github.com/gojek/stevedore/pkg/manifest"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"

	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/plugin"
)

// Plugins represents the map of plugin.Interface
type Plugins map[string]ClientPlugin

var defaultPlugins = Plugins{}

// DefaultPlugins is the Default Plugins
func DefaultPlugins() Plugins {
	return defaultPlugins
}

// ClientPlugin represents Plugin Implementation and its Client if present
type ClientPlugin struct {
	Client     *goplugin.Client
	PluginImpl plugin.Interface
}

// ConfigProviders returns the list of config providers
func (plugins Plugins) ConfigProviders() (config.Providers, error) {
	configProviders := config.Providers{}
	for k, v := range plugins {
		t, err := v.PluginImpl.Type()
		if err != nil {
			return nil, err
		}

		if t == plugin.TypeConfig {
			configProvider, ok := v.PluginImpl.(config.Provider)
			if !ok {
				return nil, fmt.Errorf("%s is not a config plugin", k)
			}

			p := config.ProviderImpl{Provider: configProvider, Name: k}
			configProviders = append(configProviders, p)
		}
	}

	return configProviders, nil
}

// ManifestProvider returns the manifest provider
func (plugins Plugins) ManifestProvider() (manifest.ProviderImpl, error) {
	for k, v := range plugins {
		t, err := v.PluginImpl.Type()
		if err != nil {
			return manifest.ProviderImpl{}, err
		}

		if t == plugin.TypeManifest {
			manifestProvider, ok := v.PluginImpl.(manifest.Provider)
			if !ok {
				return manifest.ProviderImpl{}, fmt.Errorf("%s is not a manifest plugin", k)
			}

			return manifest.ProviderImpl{Provider: manifestProvider, Name: k}, nil
		}
	}

	return manifest.ProviderImpl{}, fmt.Errorf("no manifest plugins found")
}

// PopulateFlags will populate the flags of plugin to the given command
func (plugins Plugins) PopulateFlags(cmd *cobra.Command) error {
	for pluginName, p := range plugins {
		flags, err := p.PluginImpl.Flags()
		if err != nil {
			return err
		}

		for _, f := range flags {
			flagName := fmt.Sprintf("%s-%s", pluginName, f.Name)
			cmd.PersistentFlags().StringP(flagName, f.Shorthand, f.Default, f.Usage)
			if f.Required {
				err = cmd.MarkPersistentFlagRequired(flagName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
