package stevedore

import (
	"fmt"
	"github.com/gojek/stevedore/pkg/config"

	"github.com/gojek/stevedore/pkg/file"
)

// ManifestFile represents a stevedore release request
type ManifestFile struct {
	File     string
	Manifest `json:"manifest" validate:"required"`
}

// Namespaces returns namespaces from a manifest file uniquely
func (manifestFile ManifestFile) Namespaces() []string {
	return manifestFile.Spec.Namespaces()
}

// HasBuildStep returns whether the chart has to be built for the releaseSpecification
func (manifestFile ManifestFile) HasBuildStep() bool {
	return manifestFile.Spec.HasBuildStep()
}

// ManifestFiles is a collection of stevedore release requests
type ManifestFiles []ManifestFile

func (manifestFiles ManifestFiles) applicableFor(context Context) (ManifestFiles, IgnoredReleases) {
	applicableManifestFiles := ManifestFiles{}
	var ignoreComponents IgnoredReleases

	for _, manifestFile := range manifestFiles {
		if manifestFile.IsApplicableFor(context) {
			applicableManifestFiles = append(applicableManifestFiles, manifestFile)
		} else {
			for _, releaseSpecification := range manifestFile.Spec {
				reason := fmt.Sprintf("Not applicable for the context '%s'", context.Name)
				ignoreComponent := IgnoredRelease{Name: releaseSpecification.Release.Name, Reason: reason}
				ignoreComponents = append(ignoreComponents, ignoreComponent)
			}
		}
	}

	return applicableManifestFiles, ignoreComponents
}

// Filter filters the manifests by removing ignored releaseSpecification
func (manifestFiles ManifestFiles) Filter(ignores Ignores, context Context) (ManifestFiles, IgnoredReleases) {
	result := ManifestFiles{}
	applicableManifestFiles, ignoreComponents := manifestFiles.applicableFor(context)

	for _, manifest := range applicableManifestFiles {
		newManifest := ManifestFile{
			File: manifest.File,
			Manifest: Manifest{
				DeployTo: manifest.DeployTo,
				Spec:     ReleaseSpecifications{},
			},
		}

		for _, releaseSpecification := range manifest.Spec {
			predicate := NewPredicate(releaseSpecification, context)
			matchedIgnoredReleases := ignores.Filter(predicate)

			if ignoredComponent, found := matchedIgnoredReleases.Find(releaseSpecification.Release.Name); found {
				ignoreComponents = append(ignoreComponents, ignoredComponent)
			} else {
				newManifest.Spec = append(newManifest.Spec, releaseSpecification)
			}
		}

		if len(newManifest.Spec) != 0 {
			result = append(result, newManifest)
		}
	}
	return result, ignoreComponents
}

// Enrich filters the applicable manifest, enriches it with overrides and substitutes
func (manifestFiles ManifestFiles) Enrich(
	overrides Overrides,
	stevedoreContext Context,
	envs Substitute,
	ignores Ignores,
	providers config.Providers,
) (ManifestFiles, IgnoredReleases, error) {
	result := ManifestFiles{}
	manifestErrors := file.Errors{}
	filteredManifests, ignoredComponents := manifestFiles.Filter(ignores, stevedoreContext)

	for _, manifest := range filteredManifests {
		enrichedManifest := manifest.EnrichWith(stevedoreContext, overrides)
		populatedManifest, err := enrichedManifest.Replace(stevedoreContext, envs, providers)
		if err != nil {
			manifestErrors = append(manifestErrors, file.Error{Filename: manifest.File, Reason: err})
		}

		mountedManifest, err := populatedManifest.Mount(stevedoreContext, providers)

		if err != nil {
			manifestErrors = append(manifestErrors, file.Error{Filename: manifest.File, Reason: err})
		}
		result = append(result, ManifestFile{File: manifest.File, Manifest: mountedManifest})
	}
	if len(manifestErrors) != 0 {
		return filteredManifests, ignoredComponents, manifestErrors
	}
	return result, ignoredComponents, nil
}

// HasBuildStep returns whether the chart has to be built for the releaseSpecification
func (manifestFiles ManifestFiles) HasBuildStep() bool {
	for _, manifestFile := range manifestFiles {
		if manifestFile.HasBuildStep() {
			return true
		}
	}
	return false
}
