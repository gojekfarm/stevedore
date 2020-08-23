package stevedore

// ChartSpec is to represent chart name and their dependencies
type ChartSpec struct {
	// Name of helm chart to be build/published/installed
	Name         string       `json:"name" yaml:"name"`
	Dependencies Dependencies `json:"dependencies" yaml:"dependencies"`
}

// Release represent metadata necessary for release specification
type Release struct {
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
	ChartVersion string    `json:"chartVersion,omitempty" yaml:"chartVersion,omitempty"`
	ChartSpec    ChartSpec `json:"chartSpec,omitempty" yaml:"chartSpec,omitempty"`
	// Current helm release version.
	//
	// It is used with plan. While planning we can get the current helm release version
	// and while applying we can assert that apply is done immediately after plan.
	CurrentReleaseVersion int32 `json:"currentReleaseVersion,omitempty" yaml:"currentReleaseVersion,omitempty"`
	// Required: true
	Values         Values `json:"values" yaml:"values"`
	usedSubstitute Substitute
	overrides      Overrides
}

// NewRelease returns a new Release
// NOTE: It is only used in test. Do not use it in source.
// Use pass by value to clone Release object to avoid leaving new fields
func NewRelease(name, namespace, chart, chartVersion string, chartSpec ChartSpec, currentReleaseVersion int32, values Values, usedSubstitute Substitute, overrides Overrides) Release {
	return Release{Name: name, Namespace: namespace, Chart: chart, ChartVersion: chartVersion, ChartSpec: chartSpec, CurrentReleaseVersion: currentReleaseVersion, Values: values, usedSubstitute: usedSubstitute, overrides: overrides}
}

// EnrichValues will return enriched release with final merged values
func (release Release) EnrichValues(overrides Overrides) Release {
	result := release.Values.MergeWith(overrides)
	usedSubstitute := release.usedSubstitute
	if usedSubstitute == nil {
		usedSubstitute = Substitute{}
	}

	release.Values = result
	release.overrides = overrides
	release.usedSubstitute = usedSubstitute

	// This is pass by value and it is a copy. Its not a mutation.
	return release
}

// Replace will return release with substituted values
func (release Release) Replace(with Substitute) (Release, error) {
	result, usedSubstitute, err := release.Values.Replace(with)
	if err != nil {
		return release, err
	}

	release.Values = result
	release.usedSubstitute = usedSubstitute

	// This is pass by value and it is a copy. Its not a mutation.
	return release, nil
}

// Overrides returns values enriched by overrides
func (release Release) Overrides() Overrides {
	return release.overrides
}

// HasBuildStep returns whether the chart has to be built
func (release Release) HasBuildStep() bool {
	return len(release.ChartSpec.Dependencies) != 0
}
