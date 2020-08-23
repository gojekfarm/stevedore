package helm

//import "helm.sh/helm/v3/internal/version"

// GetHelmVersion returns version of helm
func GetHelmVersion() string {
	//return version.GetVersion()
	return "3"
}

// GetBuildMetadata returns build metadata of helm
func GetBuildMetadata() string {
	//return version.Get().Version //TODO
	return "3"
}

//// SetHelmVersion sets version of helm
//func SetHelmVersion(newVersion string) {
//	version.Version = newVersion
//}
//
//// SetBuildMetadata will update the build metadata of helm
//func SetBuildMetadata(buildMetadata string) {
//	version.BuildMetadata = buildMetadata
//}
