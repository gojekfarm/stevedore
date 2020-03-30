package manifest

import (
	"github.com/gojek/stevedore/pkg/config"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/gojek/stevedore/pkg/utils/string"
)

// Info holds info manifestFiles and context
type Info struct {
	stevedore.ManifestFiles
	Ignored stevedore.IgnoredReleases
	stevedore.Context
}

// FilterBy returns the info which matches the release names and populate the release data
func (info *Info) FilterBy(responses stevedore.Responses) Info {
	newInfo := Info{Ignored: info.Ignored, Context: info.Context}

	for _, manifestFile := range info.ManifestFiles {
		releaseSpecifications := stevedore.ReleaseSpecifications{}
		manifest := manifestFile.Manifest
		releaseNames := responses.GetReleaseNames()
		for _, application := range manifest.Spec {
			releaseName := application.Release.Name
			if stringutils.Contains(releaseNames, releaseName) {
				matchedResponse := responses.Find(releaseName)
				application.Release.ChartVersion = matchedResponse.ChartVersion
				application.Release.Chart = matchedResponse.ChartName
				application.Release.ChartSpec = stevedore.ChartSpec{}
				application.Release.CurrentReleaseVersion = matchedResponse.CurrentReleaseVersion
				releaseSpecifications = append(releaseSpecifications, application)
			}
		}
		if len(releaseSpecifications) > 0 {
			newManifestFile := stevedore.ManifestFile{
				File:     manifestFile.File,
				Manifest: stevedore.Manifest{Spec: releaseSpecifications, DeployTo: manifest.DeployTo},
			}
			newInfo.ManifestFiles = append(newInfo.ManifestFiles, newManifestFile)
		}
	}
	return newInfo
}

func info(manifests stevedore.ManifestFiles,
	overrides stevedore.Overrides,
	stevedoreContext stevedore.Context,
	envs stevedore.Substitute,
	ignores stevedore.Ignores,
	providers config.Providers) (*Info, error) {

	enrichedManifestFiles, ignoredComponents, err := manifests.Enrich(overrides, stevedoreContext, envs, ignores, providers)
	if err != nil {
		return nil, err
	}

	info := &Info{
		Context:       stevedoreContext,
		ManifestFiles: enrichedManifestFiles,
		Ignored:       ignoredComponents,
	}
	return info, nil
}
