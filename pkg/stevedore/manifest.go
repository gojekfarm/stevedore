package stevedore

import (
	"fmt"
	"io"

	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/utils/string"
)

// Manifest is a Stevedore Manifest to wrap all configurations necessary for Stevedore to deploy
// Manifest is the struct representation of Stevedore yaml
type Manifest struct {
	Kind     string                `json:"kind" yaml:"kind" validate:"required"`
	Version  string                `json:"version" yaml:"version" validate:"required"`
	DeployTo Matchers              `json:"deployTo" yaml:"deployTo" validate:"required"`
	Spec     ReleaseSpecifications `json:"spec" yaml:"spec" validate:"required,dive"`
}

// EnrichWith will return enriched manifest with final merged values
func (manifest Manifest) EnrichWith(context Context, overrides Overrides) Manifest {
	enrichedApplications := manifest.Spec.EnrichWith(context, overrides)
	return Manifest{DeployTo: manifest.DeployTo, Spec: enrichedApplications}
}

// Replace will return manifest with substituted values
func (manifest Manifest) Replace(stevedoreContext Context, envs Substitute, providers config.Providers) (Manifest, error) {
	replacedApplications, err := manifest.Spec.Replace(stevedoreContext, envs, providers)
	return Manifest{DeployTo: manifest.DeployTo, Spec: replacedApplications}, err
}

// Mount will return manifest with mounted values
func (manifest Manifest) Mount(stevedoreContext Context, providers config.Providers) (Manifest, error) {
	replacedApplications, err := manifest.Spec.Mount(stevedoreContext, providers)
	return Manifest{DeployTo: manifest.DeployTo, Spec: replacedApplications}, err
}

// IsApplicableFor returns true if Manifest is applicable for the given environment
func (manifest Manifest) IsApplicableFor(context Context) bool {
	return manifest.DeployTo.Contains(context)
}

// HasBuildStep returns whether the chart has to be built for the release specification
func (manifest Manifest) HasBuildStep() bool {
	return manifest.Spec.HasBuildStep()
}

// Format implements fmt.Formatter. It accepts the formats
// 'y' (yaml)
// 'j' (json)
// '#j' (prettier json).
func (manifest Manifest) Format(f fmt.State, c rune) {
	formatAsJSONOrYaml(f, c, manifest)
}

// NewManifest to Validate the Stevedore manifest configuration
func NewManifest(reader io.Reader) (*Manifest, error) {
	manifest := &Manifest{}
	if err := ValidateAndGenerate(reader, manifest); err != nil {
		return nil, fmt.Errorf("[Schema Validation Failed] error when validating from file:\n%v", err)
	}
	return manifest, nil
}

// Manifests is a collection of Stevedore Manifest
type Manifests []Manifest

// Filter returns the filtered Manifests which are applicable for the given context,
// all the inapplicable will be returned
func (manifests Manifests) Filter(context Context) (Manifests, IgnoredReleases, bool) {
	ignoredReleases := IgnoredReleases{}
	filteredManifests := Manifests{}
	ignored := false

	for _, manifest := range manifests {
		if !manifest.IsApplicableFor(context) {
			for _, releaseSpecification := range manifest.Spec {
				reason := fmt.Sprintf("Not applicable for the context '%s'", context.Name)
				ignoredRelease := IgnoredRelease{Name: releaseSpecification.Release.Name, Reason: reason}
				ignoredReleases = append(ignoredReleases, ignoredRelease)
			}

			ignored = true
		} else {
			filteredManifests = append(filteredManifests, manifest)
		}
	}

	return filteredManifests, ignoredReleases, ignored
}

// Namespaces returns unique namespaces of manifests
func (manifests Manifests) Namespaces() []string {
	var namespaces []string

	for _, manifest := range manifests {
		namespaces = append(namespaces, manifest.Spec.Namespaces()...)
	}
	return stringutils.Unique(namespaces)
}

// HasBuildStep returns whether the chart has to be built for the releaseSpecification
func (manifests Manifests) HasBuildStep() bool {
	for _, manifest := range manifests {
		if manifest.HasBuildStep() {
			return true
		}
	}
	return false
}
