package stevedore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/pkg/config"

	"gopkg.in/yaml.v2"

	"github.com/spf13/afero"
)

// ContextEnv is holds the env key for stevedore context
const ContextEnv = "STEVEDORE_CONTEXT"

const defaultFileMode = 0660
const defaultDirMode = 0770

// AppConfigStore to fetch configurations specific to release specifications
type AppConfigStore struct {
	Host string `yaml:"host"`
}

// Configuration config to wrap all stevedore contexts and store config
type Configuration struct {
	Contexts   Contexts `yaml:"contexts"`
	UserLabels Labels   `yaml:"labels"`
	Current    string   `yaml:"current"`
	fs         afero.Fs
	filename   string
}

// NewConfigurationFromFile loads stevedore config from file
func NewConfigurationFromFile(fs afero.Fs, filename string, env config.Environment) (*Configuration, error) {
	stevedoreConfig := &Configuration{filename: filename, fs: fs}
	ok, err := afero.Exists(fs, filename)
	if err != nil {
		return nil, err
	}

	if ok {
		data, err := afero.ReadFile(fs, filename)
		if err != nil {
			return nil, fmt.Errorf("[NewConfigurationFromFile] error when reading from file: %v", err)
		}

		if err = yaml.Unmarshal(data, stevedoreConfig); err != nil {
			return nil, fmt.Errorf("[NewConfigurationFromFile] error when unmarshalling from file: %v", err)
		}

		if current, ok := OverriddenContext(env); ok {
			stevedoreConfig.Current = current
		}
	}

	return stevedoreConfig, nil
}

// OverriddenContext returns the overridden context name
// and true if its overridden else returns false
func OverriddenContext(env config.Environment) (string, bool) {
	envs := env.Fetch()
	if current, ok := envs[ContextEnv]; ok {
		return current.(string), true
	}
	return "", false
}

// Use `name` as current context if valid
func (s *Configuration) Use(name string) error {
	if _, exists := s.Contexts.Find(name); exists {
		s.Current = name
		return s.save()
	}

	return fmt.Errorf("invalid context name: %s", name)
}

// Delete the given context from all contexts
func (s *Configuration) Delete(name string) error {
	if index, exists := s.Contexts.Find(name); exists {
		s.Contexts = append(s.Contexts[:index], s.Contexts[index+1:len(s.Contexts)]...)
		if s.Current == name {
			s.Current = ""
		}
		return s.save()
	}
	return fmt.Errorf("invalid context name: %s", name)
}

// Add context to list of contexts if valid
func (s *Configuration) Add(ctx Context) error {
	if _, exists := s.Contexts.Find(ctx.Name); exists {
		return fmt.Errorf("context already exists: %s", ctx.Name)
	}

	s.Contexts = append(s.Contexts, ctx)
	return s.save()
}

// Rename the given context name to a different context
func (s *Configuration) Rename(fromCtx string, toCtx string) error {
	if _, exists := s.Contexts.Find(toCtx); exists {
		return fmt.Errorf("target context %s already exists", toCtx)
	}
	if index, exists := s.Contexts.Find(fromCtx); exists {
		s.Contexts[index].Name = toCtx

		if s.Current == fromCtx {
			s.Current = toCtx
		}

		return s.save()
	}
	return fmt.Errorf("invalid context name: %s", fromCtx)
}

// CurrentContext returns the context if set else returns appropriate error
func (s *Configuration) CurrentContext() (Context, error) {
	if s.Current == "" {
		return Context{}, fmt.Errorf("current context is not set")
	}
	if index, ok := s.Contexts.Find(s.Current); ok {
		return s.Contexts[index], nil
	}
	return Context{}, fmt.Errorf("unable to find current context %v", s.Current)
}

func (s *Configuration) save() error {
	dir, _ := filepath.Split(s.filename)
	err := s.fs.MkdirAll(dir, defaultDirMode)
	if err != nil {
		return err
	}

	f, err := s.fs.OpenFile(s.filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, defaultFileMode)
	if err != nil {
		return fmt.Errorf("[save] error when opening file %v", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			cli.Errorf("unable to close %s, reason :%v", s.filename, err)
		}
	}()

	return yaml.NewEncoder(f).Encode(s)
}

// Labels returns the user defined labels and appends the context and environment
func (s *Configuration) Labels() Labels {
	userLabels := Labels{}
	maxWeight := 0
	nextWeight := func(baseWeight int, incBy int) int {
		if baseWeight == 0 {
			return 0
		}
		return baseWeight + incBy
	}
	for _, label := range s.UserLabels {
		if maxWeight < label.Weight {
			maxWeight = label.Weight
		}
		userLabels = append(userLabels, label)
	}
	additionalLabels := []Label{
		{Name: ConditionEnvironment, Weight: nextWeight(maxWeight, 1)},
		{Name: ConditionContextName, Weight: nextWeight(maxWeight, 2)},
		{Name: ConditionApplicationName, Weight: nextWeight(maxWeight, 3)},
	}
	return append(userLabels, additionalLabels...)
}
