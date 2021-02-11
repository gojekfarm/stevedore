package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

var version string
var defaultVersion = "2.0-dev"
var configFilePath string

// PluginDirPath defines path of Stevedore Plugin Directory
var pluginDirPath string

// ManifestProvider defines name of manifest provider
var ManifestProvider string

// SetBuildVersion use this to override the version
// of the stevedore
func SetBuildVersion(v string, build string) {
	revision := strings.Replace(build, v, "", -1)
	version = fmt.Sprintf("%s%s", v, revision)
}

// BuildVersion returns build version specified at compile time
// returns defaultVersion if not present
func BuildVersion() string {
	if version != "" {
		return version
	}
	return defaultVersion
}

// SetConfigFilePath use this to override the major and minor version
// of the stevedore
func SetConfigFilePath(c string) {
	configFilePath = c
}

// FilePath returns path of stevedore.yaml
func FilePath() (string, error) {
	if path, ok := os.LookupEnv("STEVEDORE_CONFIG"); ok {
		return path, nil
	}
	if configFilePath != "" {
		return configFilePath, nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %v", err)
	}

	return filepath.Join(home, ".config", "stevedore", "config"), nil
}

// SetPluginDirPath sets the path of stevedore plugins directory
func SetPluginDirPath(path string) {
	pluginDirPath = path
}

// PluginDirPath returns path of stevedore plugin directory
func PluginDirPath() (string, error) {
	if pluginDirPath != "" {
		return pluginDirPath, nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %v", err)
	}

	return filepath.Join(home, ".config", "stevedore", "plugins"), nil
}
