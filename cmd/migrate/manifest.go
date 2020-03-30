package migrate

import (
	yml "github.com/gojek/stevedore/client/yaml"
	v1Stevedore "github.com/gojek/stevedore/cmd/migrate/types/v1/stevedore"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// ManifestStrategy represents the migration strategy
// for migrating manifest
type ManifestStrategy struct {
	fs           afero.Fs
	manifestPath string
	contexts     stevedore.Contexts
	optimize     bool
}

// NewManifestStrategy returns Strategy for migrating
// manifest
func NewManifestStrategy(fs afero.Fs, manifestPath string, contexts stevedore.Contexts, optimize bool) Strategy {
	return ManifestStrategy{
		fs:           fs,
		manifestPath: manifestPath,
		contexts:     contexts,
		optimize:     optimize,
	}
}

func (strategy ManifestStrategy) read(path string) (v1Stevedore.Manifest, error) {
	v1Manifest := v1Stevedore.Manifest{}
	err := read(strategy.fs, path, &v1Manifest)
	return v1Manifest, err
}

func (strategy ManifestStrategy) save(path string, manifest stevedore.Manifest) error {
	return save(strategy.fs, path, manifest)
}

// Name returns name for the strategy
func (strategy ManifestStrategy) Name() string {
	return "manifest"
}

// Convert performs conversion and returns error if any
func (strategy ManifestStrategy) Convert(path string) error {
	manifest, err := strategy.read(path)
	if err != nil {
		return err
	}
	converted := manifest.Convert(strategy.contexts, strategy.optimize)
	return strategy.save(path, converted)
}

// Files returns the file which will be converted
func (strategy ManifestStrategy) Files() ([]string, error) {
	files, err := yml.NewYamlFiles(strategy.fs, strategy.manifestPath)
	if err != nil {
		return nil, err
	}
	return files.Names(), nil
}

// Do performs the migration strategy and returns if any
func (strategy ManifestStrategy) Do() error {
	return migrate(strategy)
}
