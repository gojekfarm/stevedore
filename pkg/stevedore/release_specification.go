package stevedore

import (
	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/merger"
)

// ReleaseSpecification represents spec to be deployed
type ReleaseSpecification struct {
	Release    Release  `json:"release" yaml:"release" validate:"required"`
	Configs    Configs  `json:"configs" yaml:"configs"`
	DependsOn  []string `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`
	Mounts     Configs  `json:"mounts,omitempty" yaml:"mounts,omitempty"`
	substitute Substitute
}

// NewReleaseSpecification returns an instance of spec
func NewReleaseSpecification(release Release, configs Configs, substitute Substitute) ReleaseSpecification {
	return ReleaseSpecification{Release: release, Configs: configs, substitute: substitute}
}

// EnrichWith will return enriched spec with final merged values
func (spec ReleaseSpecification) EnrichWith(context Context, overrides Overrides) ReleaseSpecification {
	predicate := NewPredicate(spec, context)
	matchedOverrides := overrides.CollateBy(predicate)
	enrichedComponent := spec.Release.EnrichValues(matchedOverrides)

	spec.Release = enrichedComponent

	return spec
}

// Replace will return spec with substituted values
func (spec ReleaseSpecification) Replace(stevedoreContext Context, envs Substitute, providers config.Providers) (ReleaseSpecification, error) {
	if variables, err := spec.Release.Values.Variables(); err != nil || len(variables) == 0 {
		return spec, err
	}

	appConfig, err := spec.Configs.Fetch(providers, stevedoreContext)
	if err != nil {
		return spec, err
	}

	substitutes, err := appConfig.Merge(envs)

	if err != nil {
		return spec, err
	}

	replacedComponent, err := spec.Release.Replace(substitutes)
	if err != nil {
		return spec, err
	}

	spec.Release = replacedComponent
	spec.substitute = substitutes

	return spec, nil
}

// Mount will return spec with mounted values
// Mount will use a JSONPath to mount the fetched configs
// Mounts at root if JSONPath is empty
func (spec ReleaseSpecification) Mount(stevedoreContext Context, providers config.Providers) (ReleaseSpecification, error) {
	if len(spec.Mounts) == 0 {
		return spec, nil
	}

	mountedConfigs, err := spec.Mounts.Fetch(providers, stevedoreContext)
	if err != nil {
		return spec, err
	}

	mergedConfigs, err := merger.Merge(spec.Release.Values, mountedConfigs)
	if err != nil {
		return spec, err
	}

	spec.Release.Values = mergedConfigs
	spec.Mounts = Configs{}

	return spec, nil
}

// HasBuildStep returns whether the chart has to be built for the spec
func (spec ReleaseSpecification) HasBuildStep() bool {
	return spec.Release.HasBuildStep()
}

// SubstitutedVariables return substituted values
func (spec ReleaseSpecification) SubstitutedVariables() Substitute {
	return spec.Release.usedSubstitute
}

// ContainsDependency returns whether the spec contains the given chart name as dependency
func (spec ReleaseSpecification) ContainsDependency(chartName string) (Dependencies, bool) {
	return spec.Release.ChartSpec.Dependencies.Contains(chartName)
}
