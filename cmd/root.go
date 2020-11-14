package cmd

import (
	"os"
	"strings"

	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/cmd/config"
	"github.com/gojek/stevedore/cmd/plugin"
	"github.com/gojek/stevedore/log"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var cfgFile string
var logLevel string
var fs = afero.NewOsFs()

var pluginLoadingCommands = []string{"apply", "plan", "render", "server", "init", "plugin", "dependency"}

var rootCmd = &cobra.Command{
	Use:   "stevedore",
	Short: `Stevedore loads the cluster with containers for kubernetes to orchestrate`,
}

// Execute run the root command, entry point for the cli
func Execute() {
	rootCmd.Version = config.BuildVersion()
	filePath, err := config.FilePath()
	if err != nil {
		cli.DieIf(err, closePlugins)
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", filePath, "config file")
	defer closePlugins()

	if err := rootCmd.Execute(); err != nil {
		cli.DieIf(err, closePlugins)
	}
}

const logLevelEnv = "LOG_LEVEL"
const defaultManifestProvider = "manifests"

func initConfig() {
	envLogLevel := os.Getenv(logLevelEnv)
	if envLogLevel != "" {
		log.SetLogLevel(envLogLevel)
	} else {
		log.SetLogLevel(logLevel)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "error", "set the logger level")
}

// Load plugins before init()
var _ = func() error {
	var pluginDirPath string
	defaultPluginDirPath, err := config.PluginDirPath()
	if err != nil {
		cli.DieIf(err, func() {})
	}

	rootCmd.PersistentFlags().StringVarP(&pluginDirPath, "plugin-dir-path", "p", defaultPluginDirPath,
		"plugin directory path")
	rootCmd.PersistentFlags().StringVarP(&config.ManifestProvider, "manifest-provider", "", defaultManifestProvider,
		"manifest provider to use")

	pluginDirPath = flagHackLookup("--plugin-dir-path")
	if len(pluginDirPath) == 0 {
		pluginDirPath = defaultPluginDirPath
	} else {
		config.SetPluginDirPath(pluginDirPath)
		exists, err := afero.DirExists(fs, pluginDirPath)
		if err != nil {
			cli.DieIf(err, func() {})
		}

		if !exists {
			cli.Warn("WARN: Directory ", pluginDirPath, " does not exist. Will use default plugins...")
		}
	}

	manifestProvider := flagHackLookup("--manifest-provider")
	if manifestProvider == "" {
		manifestProvider = defaultManifestProvider
	}
	config.ManifestProvider = manifestProvider

	loader := plugin.NewPluginLoader(pluginDirPath, config.ManifestProvider, fs)

	if shouldLoadPlugin(pluginLoadingCommands) {
		err = loader.Load()
		cli.DieIf(err, closePlugins)
	}

	return nil
}()

func shouldLoadPlugin(nonPluginRelatedCommands []string) bool {
	callingCommand := strings.Join(os.Args, " ")
	for _, command := range nonPluginRelatedCommands {
		if strings.Contains(callingCommand, command) {
			return true
		}
	}
	return false
}

func closePlugins() {
	loader, err := plugin.GetPluginLoader()
	if err != nil {
		cli.DieIf(err, func() {})
	}
	loader.Close()
}
