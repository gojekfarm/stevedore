package helpers

import "github.com/gojek/stevedore/pkg/stevedore"

// CreateManifestFor is a helper function to create a stevedore.Manifest
// for the given chart name with release name.
func CreateManifestFor(chart string, releaseName string) stevedore.Manifest {
	return stevedore.Manifest{
		Kind:     stevedore.KindStevedoreManifest,
		Version:  stevedore.ManifestCurrentVersion,
		DeployTo: nil,
		Spec:     nil,
	}
}
