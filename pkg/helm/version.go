package helm

import "k8s.io/helm/pkg/version"

// GetHelmVersion returns version of helm
func GetHelmVersion() string {
	return version.Version
}

// GetBuildMetadata returns build metadata of helm
func GetBuildMetadata() string {
	return version.BuildMetadata
}

// SetHelmVersion sets version of helm
func SetHelmVersion(newVersion string) {
	version.Version = newVersion
}

// SetBuildMetadata will update the build metadata of helm
func SetBuildMetadata(buildMetadata string) {
	version.BuildMetadata = buildMetadata
}
