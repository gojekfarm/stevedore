package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Application represents spec to be deployed
type Application struct {
	Component Component         `json:"component" yaml:"component" validate:"required"`
	Configs   stevedore.Configs `json:"configs" yaml:"configs"`
	DependsOn []string          `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`
	Mounts    stevedore.Configs `json:"mounts,omitempty" yaml:"mounts,omitempty"`
}

// ReleaseSpecification convert applications to stevedore.ReleaseSpecification
func (application Application) ReleaseSpecification() stevedore.ReleaseSpecification {
	return stevedore.ReleaseSpecification{
		Release:   application.Component.Release(),
		Configs:   application.Configs,
		DependsOn: application.DependsOn,
		Mounts:    application.Mounts,
	}
}
