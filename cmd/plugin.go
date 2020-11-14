package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/cmd/config"
	"github.com/gojek/stevedore/cmd/plugin"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage stevedore plugins",
}

var pluginListCmd = &cobra.Command{
	Use:           "list",
	Aliases:       []string{"ls"},
	Short:         "List stevedore plugins",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginLoader, err := plugin.GetPluginLoader()
		if err != nil {
			return err
		}
		stevedorePlugins := pluginLoader.GetAllPlugins()
		pluginInfos, err := getPluginsInfo(stevedorePlugins)
		if err != nil {
			return err
		}

		if len(stevedorePlugins) == 0 {
			fmt.Println("No plugins installed")
			return nil
		}

		table := cli.NewTableRenderer(os.Stdout)
		table.SetHeader([]string{"NAME", "TYPE", "VERSION"})

		for _, pluginInfo := range pluginInfos {
			table.Append([]string{pluginInfo.name, pluginInfo.pluginType.String(), pluginInfo.version})
		}

		table.Render()
		return nil
	},
}

var pluginHelpCmd = &cobra.Command{
	Use:     "help <plugin-name>",
	Aliases: []string{"desc", "describe"},
	Short:   "Show help for stevedore plugin",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("pass exactly one <plugin-name> as argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName := args[0]
		pluginLoader, err := plugin.GetPluginLoader()
		if err != nil {
			return err
		}

		plugins := pluginLoader.GetAllPlugins()
		for name, p := range plugins {
			if name == pluginName {
				help, err := p.PluginImpl.Help()
				if err != nil {
					return fmt.Errorf("error fetching help for the plugin %s: %v", pluginName, err)
				}

				fmt.Println(help)
				return nil
			}
		}

		return fmt.Errorf("plugin %s not found", pluginName)
	},
}

// PluginInfo represents details of the plugin
type PluginInfo struct {
	name       string
	pluginType pkgPlugin.Type
	version    string
}

func getPluginsInfo(plugins provider.Plugins) ([]PluginInfo, error) {
	result := make([]PluginInfo, 0, len(plugins))
	for name, p := range plugins {
		pluginType, err := p.PluginImpl.Type()
		if err != nil {
			return nil, fmt.Errorf("error while getting type of plugin %s: %v", name, err)
		}
		version, err := p.PluginImpl.Version()
		if err != nil {
			return nil, fmt.Errorf("error while getting version for plugin %s: %v", name, err)
		}
		result = append(result, PluginInfo{name: name, pluginType: pluginType, version: version})
	}

	return result, nil
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <name> <path|uri>",
	Short: "Install a plugin",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("pass 2 arguments: <name> <path|uri>")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName := args[0]
		sourcePath := args[1]
		pluginDir, err := config.PluginDirPath()
		if err != nil {
			return err
		}

		pluginFileContent, err := readPlugin(fs, sourcePath)
		if err != nil {
			return err
		}

		pluginPath := filepath.Join(pluginDir, fmt.Sprintf(plugin.NameFormat, pluginName))
		err = afero.WriteFile(fs, pluginPath, pluginFileContent, 0744)
		if err != nil {
			return err
		}

		fmt.Println("plugin successfully installed to", pluginPath)
		return nil
	},
}

type sourceType int

const (
	local sourceType = iota
	remote
)

func readPlugin(fs afero.Fs, sourcePath string) ([]byte, error) {
	source, err := findSource(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("error reading plugin: %v", err)
	}

	switch source {
	case local:
		return afero.ReadFile(fs, sourcePath)
	case remote:
		return downloadFile(sourcePath)
	}

	return nil, fmt.Errorf("error reading plugin")
}

func findSource(src string) (sourceType, error) {
	if !(strings.HasPrefix(src, "http") || strings.HasPrefix(src, "https")) {
		return local, nil
	}
	_, err := url.ParseRequestURI(src)
	if err == nil {
		return remote, nil
	}
	return remote, fmt.Errorf("error occurred while trying to find the source")
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading %s: %v", url, err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			cli.Error(err.Error())
		}
	}()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error getting response: %v", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("error closing response: %v", err)
	}

	return buf.Bytes(), nil
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginHelpCmd)

	rootCmd.AddCommand(pluginCmd)
}
