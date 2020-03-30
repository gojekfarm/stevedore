package migrate

import (
	yml "github.com/gojek/stevedore/client/yaml"
	v1Stevedore "github.com/gojek/stevedore/cmd/migrate/types/v1/stevedore"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// OverrideStrategy represents the migration strategy
// for migrating override
type OverrideStrategy struct {
	fs           afero.Fs
	overridePath string
	optimize     bool
}

func (strategy OverrideStrategy) read(path string) (v1Stevedore.Overrides, error) {
	v1Manifest := v1Stevedore.Overrides{}
	err := read(strategy.fs, path, &v1Manifest)
	return v1Manifest, err
}

func (strategy OverrideStrategy) save(path string, overrides stevedore.Overrides) error {
	return save(strategy.fs, path, overrides)
}

// Name return the override strategy name
func (strategy OverrideStrategy) Name() string {
	return "override"
}

// Do perform the override strategy
func (strategy OverrideStrategy) Do() error {
	return migrate(strategy)
}

// Files returns the matching override files
func (strategy OverrideStrategy) Files() ([]string, error) {
	files, err := yml.NewYamlFiles(strategy.fs, strategy.overridePath)
	if err != nil {
		return nil, err
	}
	return files.Names(), nil
}

// Convert converts the override to newer format and returns error if any
func (strategy OverrideStrategy) Convert(file string) error {
	overrides, err := strategy.read(file)
	if err != nil {
		return err
	}
	converted := overrides.Convert()
	return strategy.save(file, converted)
}

// NewOverrideStrategy returns Strategy for migrating
// manifest
func NewOverrideStrategy(fs afero.Fs, overridePath string, optimize bool) Strategy {
	return OverrideStrategy{
		fs:           fs,
		overridePath: overridePath,
		optimize:     optimize,
	}
}
