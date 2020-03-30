package provider

import (
	"github.com/gojek/stevedore/client/yaml"
	"github.com/gojek/stevedore/pkg/file"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// OverrideProvider is the OverrideSpecification Provider interface
type OverrideProvider interface {
	Overrides() (stevedore.Overrides, error)
}

// DefaultOverrideProvider represents the default override provider
// which reads the overrides from given dir
type DefaultOverrideProvider struct {
	fs   afero.Fs
	path string
}

// Overrides returns all the overrides
func (provider DefaultOverrideProvider) Overrides() (stevedore.Overrides, error) {
	if provider.path == "" {
		return stevedore.EmptyOverrides(), nil
	}

	yamlFiles, err := yaml.NewYamlFiles(provider.fs, provider.path)

	if _, ok := err.(yaml.EmptyFolderError); ok {
		return stevedore.EmptyOverrides(), nil
	} else if err != nil {
		return stevedore.EmptyOverrides(), err
	}

	overrideErrors := file.Errors{}
	specs := stevedore.OverrideSpecifications{}
	for _, yamlFile := range yamlFiles {
		ok, err := yamlFile.Check(stevedore.KindStevedoreOverride, stevedore.OverrideCurrentVersion)
		if err != nil {
			overrideErrors = append(overrideErrors, file.Error{Filename: yamlFile.Name, Reason: err})
			continue
		}
		if ok {
			result, err := stevedore.NewOverrides(yamlFile.Reader())
			if err != nil {
				overrideErrors = append(overrideErrors, file.Error{Filename: yamlFile.Name, Reason: err})
				continue
			}
			for index := range result.Spec {
				result.Spec[index].FileName = yamlFile.Name
			}
			specs = append(specs, result.Spec...)
		}
	}

	if len(overrideErrors) != 0 {
		return stevedore.EmptyOverrides(), overrideErrors
	}
	return stevedore.Overrides{Spec: specs}, nil
}

// NewOverrideProvider returns new default override provider
func NewOverrideProvider(fs afero.Fs, path string) OverrideProvider {
	return DefaultOverrideProvider{fs: fs, path: path}
}
