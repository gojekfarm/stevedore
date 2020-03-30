package migrate

import (
	v1Stevedore "github.com/gojek/stevedore/cmd/migrate/types/v1/stevedore"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// IgnoreStrategy represents the migration strategy
// for migrating manifest
type IgnoreStrategy struct {
	fs    afero.Fs
	files []string
}

func (strategy IgnoreStrategy) read(path string) (v1Stevedore.Ignores, error) {
	v1Ignores := v1Stevedore.Ignores{}
	err := read(strategy.fs, path, &v1Ignores)
	return v1Ignores, err
}

func (strategy IgnoreStrategy) save(path string, ignores stevedore.Ignores) error {
	return save(strategy.fs, path, ignores)
}

// Name returns name for the strategy
func (strategy IgnoreStrategy) Name() string {
	return "ignore"
}

// Files returns the file names which will be converted
func (strategy IgnoreStrategy) Files() ([]string, error) {
	return strategy.files, nil
}

// Convert perform the conversion and report if any
func (strategy IgnoreStrategy) Convert(path string) error {
	ignores, err := strategy.read(path)
	if err != nil {
		return err
	}
	converted := ignores.Convert()
	return strategy.save(path, converted)
}

// Do performs the migration strategy and returns if any
func (strategy IgnoreStrategy) Do() error {
	return migrate(strategy)
}

// NewIgnoreStrategy returns Strategy for migrating
// ignores
func NewIgnoreStrategy(fs afero.Fs, files []string) Strategy {
	return IgnoreStrategy{fs: fs, files: files}
}
