package stevedore

import (
	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/utils/string"
)

// ReleaseSpecifications is a collection of helm release specification
type ReleaseSpecifications []ReleaseSpecification

// EnrichWith will return enriched applications with final merged values
func (specs ReleaseSpecifications) EnrichWith(context Context, overrides Overrides, labels Labels) ReleaseSpecifications {
	applications := ReleaseSpecifications{}
	for _, app := range specs {
		applications = append(applications, app.EnrichWith(context, overrides, labels))
	}
	return applications
}

// Replace will return applications with substituted values
func (specs ReleaseSpecifications) Replace(stevedoreContext Context, envs Substitute, providers config.Providers) (ReleaseSpecifications, error) {
	replacedApps := ReleaseSpecifications{}
	substituteErrors := SubstituteError{}
	for _, app := range specs {
		replacedApp, err := app.Replace(stevedoreContext, envs, providers)
		if err != nil {
			if substituteErr, ok := err.(SubstituteError); ok {
				substituteErrors = append(substituteErrors, substituteErr...)
			} else {
				return specs, err
			}
		} else {
			replacedApps = append(replacedApps, replacedApp)
		}
	}

	if len(substituteErrors) != 0 {
		return specs, substituteErrors
	}
	return replacedApps, nil
}

// Mount will return applications with mounted values
func (specs ReleaseSpecifications) Mount(stevedoreContext Context, providers config.Providers) (ReleaseSpecifications, error) {
	mountedApps := ReleaseSpecifications{}
	for _, app := range specs {
		replacedApp, err := app.Mount(stevedoreContext, providers)
		if err != nil {
			return specs, err
		}
		mountedApps = append(mountedApps, replacedApp)
	}
	return mountedApps, nil
}

// Namespaces will return the list of namespaces
func (specs ReleaseSpecifications) Namespaces() []string {
	var namespaces []string
	for _, app := range specs {
		namespaces = append(namespaces, app.Release.Namespace)
	}
	return stringutils.Unique(namespaces)
}

// HasBuildStep returns whether the chart has to be built for the release specification
func (specs ReleaseSpecifications) HasBuildStep() bool {
	for _, app := range specs {
		if app.HasBuildStep() {
			return true
		}
	}
	return false
}
