package provider

import (
	"fmt"

	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// ContextProvider is the Context ProviderImpl interface
type ContextProvider interface {
	Context() (stevedore.Context, error)
}

// DefaultContextProvider represents the default context provider
// which reads the context from file
type DefaultContextProvider struct {
	fs          afero.Fs
	file        string
	environment config.Environment
}

// Context returns the current context which is in use
func (provider DefaultContextProvider) Context() (stevedore.Context, error) {
	configurations, err := stevedore.NewConfigurationFromFile(provider.fs, provider.file, provider.environment)
	if err != nil {
		return stevedore.Context{}, fmt.Errorf("[currentContext] %v", err)
	}
	return configurations.CurrentContext()
}

// NewContextProvider returns new instance of context.ContextProvider
func NewContextProvider(fs afero.Fs, file string, environment config.Environment) ContextProvider {
	return DefaultContextProvider{fs: fs, file: file, environment: environment}
}
