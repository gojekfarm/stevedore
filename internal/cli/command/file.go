package command

import (
	"fmt"
	"github.com/gojek/stevedore/pkg/stevedore"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// SaveEnvs save envs into temporary dir
func SaveEnvs(envs stevedore.Env) (string, error) {
	baseDir, err := ioutil.TempDir("", "envs")
	if err != nil {
		return "", err
	}
	for index, env := range envs.Spec {
		envFilePath := path.Join(baseDir, fmt.Sprintf("env-%d.yaml", index))
		data, err := yaml.Marshal(env)
		if err != nil {
			return envFilePath, err
		}

		err = ioutil.WriteFile(envFilePath, data, os.ModePerm)
		if err != nil {
			return envFilePath, err
		}
	}
	return baseDir, nil
}

// SaveManifest save manifest into temporary dir
func SaveManifest(manifest stevedore.Manifest) (string, error) {
	baseDir, err := ioutil.TempDir("", "manifest")
	if err != nil {
		return "", err
	}

	manifestFilePath := path.Join(baseDir, "manifest.yaml")
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return manifestFilePath, err
	}
	return manifestFilePath, ioutil.WriteFile(manifestFilePath, data, os.ModePerm)
}

// SaveOverrides save overrides into temporary dir
func SaveOverrides(overrides stevedore.Overrides) (string, error) {
	baseDir, err := ioutil.TempDir("", "overrides")
	if err != nil {
		return "", err
	}
	for index, override := range overrides.Spec {
		overrideFilePath := path.Join(baseDir, fmt.Sprintf("override-%d.yaml", index))
		data, err := yaml.Marshal(override)
		if err != nil {
			return overrideFilePath, err
		}

		err = ioutil.WriteFile(overrideFilePath, data, os.ModePerm)
		if err != nil {
			return overrideFilePath, err
		}
	}
	return baseDir, nil
}

// Namespaces represents kubernetes namespace
type Namespaces struct {
	Names []string `yaml:"namespaces"`
}

// SaveNamespace save namespaces into temporary dir
func SaveNamespace(ns Namespaces) (string, error) {
	baseDir, err := ioutil.TempDir("", "namespace")
	if err != nil {
		return "", err
	}

	namespacesFilePath := path.Join(baseDir, "Namespaces.yaml")
	data, err := yaml.Marshal(ns)
	if err != nil {
		return namespacesFilePath, err
	}
	return namespacesFilePath, ioutil.WriteFile(namespacesFilePath, data, os.ModePerm)
}
