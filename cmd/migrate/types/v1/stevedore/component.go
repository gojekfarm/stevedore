package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Component represent metadata necessary for release specification
type Component struct {
	// Name of the helm release
	// Required: true
	Name string `json:"name" yaml:"name" validate:"required"`
	// Namespace in which the release needs to be deployed
	// Required: true
	Namespace string `json:"namespace" yaml:"namespace" validate:"required"`
	// Name of helm chart to be installed/upgraded.
	//
	// Use ChartSpec if you want to dynamically build the chart and install.
	//
	// Do not use Chart and ChartSpec together
	Chart string `json:"chart" yaml:"chart,omitempty"`
	// Chart Version to be deployed (By default, Stevedore will install latest version)
	ChartVersion string              `json:"chartVersion,omitempty" yaml:"chartVersion,omitempty"`
	ChartSpec    stevedore.ChartSpec `json:"chartSpec,omitempty" yaml:"chartSpec,omitempty"`
	// Current helm release version.
	//
	// It is used with plan. While planning we can get the current helm release version
	// and while applying we can assert that apply is done immediately after plan.
	CurrentReleaseVersion int32 `json:"currentReleaseVersion,omitempty" yaml:"currentReleaseVersion,omitempty"`
	// Required: true
	Values stevedore.Values `json:"values" yaml:"values"`
	// Set to true to use privileged kube-system tiller to upstall the release
	Privileged bool `json:"privileged,omitempty" yaml:"privileged,omitempty"`
}

// Release convert component as stevedore.Release
func (component Component) Release() stevedore.Release {
	return stevedore.Release{
		Name:                  component.Name,
		Namespace:             component.Namespace,
		Chart:                 component.Chart,
		ChartVersion:          component.ChartVersion,
		ChartSpec:             component.ChartSpec,
		CurrentReleaseVersion: component.CurrentReleaseVersion,
		Values:                component.Values,
		Privileged:            component.Privileged,
	}
}
