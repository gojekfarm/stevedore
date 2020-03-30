package plugin

import (
	"fmt"
	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/log"
	pluginPkg "github.com/gojek/stevedore/pkg/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/afero"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// NameFormat is the format in which plugins should be named
const NameFormat = "stevedore-%s-plugin"

// Loader loads all the plugins in the given plugin directory
type Loader struct {
	pluginDir        string
	fs               afero.Fs
	plugins          provider.Plugins
	manifestProvider string
}

var loader *Loader

// NewPluginLoader creates plugin loader if not present
func NewPluginLoader(pluginDir string, manifestProvider string, fs afero.Fs) *Loader {
	if loader != nil {
		return loader
	}
	loader = &Loader{pluginDir: pluginDir, manifestProvider: manifestProvider, fs: fs}
	return loader
}

// GetPluginLoader returns plugin loader if present, returns error otherwise
func GetPluginLoader() (Loader, error) {
	if loader == nil {
		return Loader{}, fmt.Errorf("loader not created")
	}
	return *loader, nil
}

// Load returns all the plugins in given plugin directory
func (l *Loader) Load() error {
	pattern := filepath.Join(l.pluginDir, fmt.Sprintf(NameFormat, "*"))
	pluginPaths, err := afero.Glob(l.fs, pattern)
	if err != nil {
		return fmt.Errorf("error globbing for plugins: %v", err)
	}

	defaultPlugins := provider.DefaultPlugins()
	result := defaultPlugins
	for _, pluginPath := range pluginPaths {
		rpcPlugin, err := createRPCPlugin(pluginPath)
		if err != nil {
			return fmt.Errorf("error creating config provider: %v", err)
		}
		result[getPluginName(pluginPath)] = *rpcPlugin
	}

	l.plugins = result
	return nil
}

// Close closes all open plugin connections
func (l *Loader) Close() {
	for _, p := range l.plugins {
		rpcClient := p.Client
		if rpcClient != nil {
			client, err := rpcClient.Client()
			if err != nil {
				log.Error(err)
			}
			err = client.Close()
			if err != nil {
				log.Error(err)
			}
		}
		err := p.PluginImpl.Close()
		if err != nil {
			log.Error(err)
		}
	}
}

// GetAllEligiblePlugins returns all the plugins in given plugin directory
func (l Loader) GetAllEligiblePlugins() (provider.Plugins, error) {
	configPlugins, err := l.GetPluginsByType(pluginPkg.TypeConfig)
	if err != nil {
		return nil, err
	}
	manifestPlugin, err := l.GetManifestPlugin()
	if err != nil {
		return nil, err
	}
	for k, v := range manifestPlugin {
		configPlugins[k] = v
	}
	return configPlugins, nil
}

// GetAllPlugins returns all the plugins in given plugin directory
func (l Loader) GetAllPlugins() provider.Plugins {
	return l.plugins
}

// GetPluginsByType returns all the plugins of given type in given plugin directory
func (l Loader) GetPluginsByType(pluginType pluginPkg.Type) (provider.Plugins, error) {
	allStevedorePlugins := l.plugins
	filteredStevedorePlugins := make(provider.Plugins, len(allStevedorePlugins))
	for pluginName, stevedorePlugin := range allStevedorePlugins {
		t, err := stevedorePlugin.PluginImpl.Type()
		if err != nil {
			return nil, err
		}

		if t == pluginType {
			filteredStevedorePlugins[pluginName] = stevedorePlugin
		}
	}

	return filteredStevedorePlugins, nil
}

// GetManifestPlugin df
func (l *Loader) GetManifestPlugin() (provider.Plugins, error) {
	acc := provider.Plugins{}
	for key, value := range l.plugins {
		pluginType, err := value.PluginImpl.Type()
		if err != nil {
			return nil, err
		}
		if (pluginType == pluginPkg.TypeManifest) && (key == l.manifestProvider) {
			acc[key] = value
		}
	}
	return acc, nil
}

func getPluginName(pluginPath string) string {
	filename := filepath.Base(pluginPath)
	regexpPattern := "^" + fmt.Sprintf(NameFormat, "(.+)") + "$"
	re := regexp.MustCompile(regexpPattern)
	groups := re.FindStringSubmatch(filename)
	if len(groups) < 2 {
		panic(fmt.Sprintf("error in stevedore plugin filename %s", filename))
	}
	return strings.Split(groups[1], "-")[0]
}

func createRPCPlugin(pluginPath string) (*provider.ClientPlugin, error) {
	var pluginName string
	var plugins map[string]plugin.Plugin
	if strings.Contains(filepath.Base(pluginPath), pluginPkg.TypeManifest.String()) {
		pluginName = pluginPkg.ManifestProviderKey
		plugins = map[string]plugin.Plugin{
			pluginName: &pluginPkg.ManifestPlugin{},
		}
	} else if strings.Contains(filepath.Base(pluginPath), pluginPkg.TypeConfig.String()) {
		pluginName = pluginPkg.ConfigProviderKey
		plugins = map[string]plugin.Plugin{
			pluginName: &pluginPkg.ConfigPlugin{},
		}
	} else {
		return nil, fmt.Errorf("invalid plugin name: plugin type should present in plugin name")
	}
	client := createClient(pluginPath, plugins)

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("error starting rpc client: %v", err)
	}
	rawPluginConfigProvider, err := rpcClient.Dispense(pluginName)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("error dispensing provider: %v", err)
	}
	configProvider, ok := rawPluginConfigProvider.(pluginPkg.Interface)
	if !ok {
		client.Kill()
		return nil, fmt.Errorf("dispensed provider is not a ConfigProvider")
	}

	return &provider.ClientPlugin{Client: client, PluginImpl: configProvider}, nil
}

func createClient(pluginPath string, plugins map[string]plugin.Plugin) *plugin.Client {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr,
		Level:  hclog.Info,
	})
	handshakeConfig := pluginPkg.CreateHandshakeConfig()
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         plugins,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
	})
	return client
}
