package dependency

import (
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/pkg/manifest"
)

// ShowAction holds necessary information for build action
type ShowAction struct{}

// Do perform show action
func (action ShowAction) Do(impl manifest.ProviderImpl) error {
	cli.Info("Not yet implemented")
	return nil
}
