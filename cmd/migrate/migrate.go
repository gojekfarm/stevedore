package migrate

import (
	"fmt"
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/pkg/file"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Strategy represents migration strategy
type Strategy interface {
	Name() string
	Do() error
	Files() ([]string, error)
	Convert(file string) error
}

// Perform migrates manifest file with list of environments into list of matchers
func Perform(strategies ...Strategy) error {
	hasError := false
	cli.Infof("%d migration(s) will be performed", len(strategies))
	for index, strategy := range strategies {
		cli.Infof("\n%d. migrating %s file(s)", index+1, strategy.Name())
		if err := strategy.Do(); err != nil {
			if fileErrors, ok := err.(file.Errors); ok {
				cli.Errorf("\n%d error(s) occurred while migrating %s file(s):", len(fileErrors), strategy.Name())
				for index, err := range fileErrors {
					cli.Errorf("\n%d. %v", index+1, err)
				}
			}
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("there were some problem when migrating files")
	}
	return nil
}

func migrate(strategy Strategy) error {
	files, err := strategy.Files()
	if err != nil {
		return err
	}

	cli.Infof("found %d files", len(files))
	errors := file.Errors{}
	for index, yamlFile := range files {
		cli.Info(fmt.Sprintf("\t %d. %s...", index+1, yamlFile))
		err := strategy.Convert(yamlFile)
		if err != nil {
			errors = append(errors, file.Error{Filename: yamlFile, Reason: err})
			continue
		}
	}
	if len(errors) != 0 {
		return errors
	}
	return nil
}

func read(fs afero.Fs, path string, out interface{}) error {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}

func save(fs afero.Fs, path string, in interface{}) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, path, data, 066)
}
