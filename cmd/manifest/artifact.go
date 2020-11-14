package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gojek/stevedore/log"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Artifact interface to generate artifact
type Artifact interface {
	Save(info Info) error
}

// DefaultArtifact will not generate any artifact files
type DefaultArtifact struct{}

// Save will not generate any artifact files
func (artifact DefaultArtifact) Save(info Info) error {
	return nil
}

// PersistentArtifact will generate artifact files
type PersistentArtifact struct {
	path        string
	persistence Persistence
	fs          afero.Fs
}

func (artifact PersistentArtifact) write(filepath string, data []byte) error {
	return artifact.persistence.write(filepath, data)
}

// Save will generate artifact files and save it in the provided path
func (artifact PersistentArtifact) Save(info Info) error {
	artifactPath := artifact.path

	if err := artifact.fs.RemoveAll(artifactPath); err != nil {
		return fmt.Errorf("error cleaning up %v directory: %v", artifactPath, err)
	}
	if err := artifact.fs.Mkdir(artifactPath, os.ModePerm); err != nil {
		return fmt.Errorf("error while creating %v directory: %v", artifactPath, err)
	}

	for _, manifestFile := range info.ManifestFiles {
		for _, releaseSpecification := range manifestFile.Manifest.Spec {
			manifest := stevedore.Manifest{
				Kind:     stevedore.KindStevedoreManifest,
				Version:  stevedore.ManifestCurrentVersion,
				DeployTo: manifestFile.Manifest.DeployTo,
				Spec:     stevedore.ReleaseSpecifications{releaseSpecification},
			}
			artifactFilePath := filepath.Join(artifactPath, fmt.Sprintf("%s.yaml", releaseSpecification.Release.Name))

			data, err := yaml.Marshal(manifest)

			if err != nil {
				return fmt.Errorf("error marshalling manifest of aplication: %v %v", releaseSpecification.Release.Name, err)
			}

			log.Debug("creating manifest artifact: ", artifactFilePath)

			err = artifact.write(artifactFilePath, data)

			if err != nil {
				return fmt.Errorf("error writing manifest file to %v : %v", artifactFilePath, err)
			}

			log.Debug("successfully generated artifact at: ", artifactFilePath)

		}
	}
	return nil
}

// NewArtifact returns an artifact
func NewArtifact(fs afero.Fs, save bool, path string) Artifact {
	if save {
		var persistence Persistence = ConsolePersistence{}
		if path != "-" {
			persistence = DiskPersistence{fs: fs}
		}
		return PersistentArtifact{path: path, persistence: persistence, fs: fs}
	}

	return DefaultArtifact{}
}
