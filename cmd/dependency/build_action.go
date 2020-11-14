package dependency

import (
	"context"

	"github.com/gojek/stevedore/cmd/manifest"
	manifestProvider "github.com/gojek/stevedore/pkg/manifest"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// BuildAction holds necessary information for build action
type BuildAction struct {
	fs            afero.Fs
	repoName      string
	artifactsPath string
}

// NewBuildAction returns action
func NewBuildAction(fs afero.Fs, repoName string, artifactsPath string) Action {
	return BuildAction{fs: fs, repoName: repoName, artifactsPath: artifactsPath}
}

// Do performs build action
func (action BuildAction) Do(mProvider manifestProvider.ProviderImpl) error {
	manifestFiles, err := mProvider.Provider.Manifests(mProvider.Context)
	if err != nil {
		return err
	}

	dependencyBuilder, err := stevedore.CreateDependencyBuilder(manifestFiles, action.repoName)
	if err != nil {
		return err
	}

	generateArtifact := len(action.artifactsPath) != 0
	artifact := manifest.NewArtifact(action.fs, generateArtifact, action.artifactsPath)

	manifests, err := dependencyBuilder.Build(context.TODO(), manifestFiles)
	if err != nil {
		return err
	}

	info := manifest.Info{ManifestFiles: manifests}
	return artifact.Save(info)
}
