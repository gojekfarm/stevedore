package provider

import (
	"github.com/gojek/stevedore/client/yaml"
	"github.com/gojek/stevedore/pkg/file"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// EnvProvider is the Env Provider interface
type EnvProvider interface {
	Envs() (EnvsFiles, error)
}

// DefaultEnvProvider represents the default env provider
// which reads the envs from given dir
type DefaultEnvProvider struct {
	fs   afero.Fs
	path string
}

// Envs returns all the envs
func (provider DefaultEnvProvider) Envs() (EnvsFiles, error) {
	envs := EnvsFiles{}
	if provider.path == "" {
		return envs, nil
	}

	yamlFiles, err := yaml.NewYamlFiles(provider.fs, provider.path)

	if _, ok := err.(yaml.EmptyFolderError); ok {
		return envs, nil
	} else if err != nil {
		return nil, err
	}

	envErrors := file.Errors{}
	for _, yamlFile := range yamlFiles {
		ok, err := yamlFile.Check(stevedore.KindStevedoreEnv, stevedore.EnvCurrentVersion)
		if err != nil {
			envErrors = append(envErrors, file.Error{Filename: yamlFile.Name, Reason: err})
			continue
		}
		if ok {
			result, err := stevedore.NewEnv(yamlFile.Reader())
			if err != nil {
				envErrors = append(envErrors, file.Error{Filename: yamlFile.Name, Reason: err})
				continue
			}
			envs = append(envs, EnvsFile{Name: yamlFile.Name, EnvSpecifications: result.Spec})
		}
	}

	if len(envErrors) != 0 {
		return nil, envErrors
	}
	return envs, nil
}

// NewEnvProvider returns new default env provider
func NewEnvProvider(fs afero.Fs, path string) EnvProvider {
	return DefaultEnvProvider{fs: fs, path: path}
}
