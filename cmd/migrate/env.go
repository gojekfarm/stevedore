package migrate

import (
	yml "github.com/gojek/stevedore/client/yaml"
	v1Stevedore "github.com/gojek/stevedore/cmd/migrate/types/v1/stevedore"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// EnvStrategy represents the migration strategy
// for migrating env
type EnvStrategy struct {
	fs       afero.Fs
	envPath  string
	optimize bool
}

func (strategy EnvStrategy) read(path string) (v1Stevedore.Envs, error) {
	v1Envs := v1Stevedore.Envs{}
	err := read(strategy.fs, path, &v1Envs)
	return v1Envs, err
}

func (strategy EnvStrategy) save(path string, envs stevedore.Env) error {
	return save(strategy.fs, path, envs)
}

// Name return the env strategy name
func (strategy EnvStrategy) Name() string {
	return "env"
}

// Do perform the env strategy
func (strategy EnvStrategy) Do() error {
	return migrate(strategy)
}

// Files returns the matching env files
func (strategy EnvStrategy) Files() ([]string, error) {
	files, err := yml.NewYamlFiles(strategy.fs, strategy.envPath)
	if err != nil {
		return nil, err
	}
	return files.Names(), nil
}

// Convert converts the env to newer format and returns error if any
func (strategy EnvStrategy) Convert(file string) error {
	envs, err := strategy.read(file)
	if err != nil {
		return err
	}
	converted := envs.Convert()
	return strategy.save(file, converted)
}

// NewEnvStrategy returns Strategy for migrating
// manifest
func NewEnvStrategy(fs afero.Fs, envPath string, optimize bool) Strategy {
	return EnvStrategy{
		fs:       fs,
		envPath:  envPath,
		optimize: optimize,
	}
}
