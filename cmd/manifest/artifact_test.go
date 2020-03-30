package manifest

import (
	"path/filepath"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestNewArtifact(t *testing.T) {
	t.Run("should return default artifact if save is disabled", func(t *testing.T) {
		mapFs := afero.NewMemMapFs()
		artifact := NewArtifact(mapFs, false, "")

		assert.IsType(t, DefaultArtifact{}, artifact)
	})

	t.Run("should return persistent artifact with writeTo as Console if save is enabled and path is -", func(t *testing.T) {
		mapFs := afero.NewMemMapFs()
		artifact := NewArtifact(mapFs, true, "-")

		if assert.IsType(t, PersistentArtifact{}, artifact) {
			pArtifact, _ := artifact.(PersistentArtifact)
			assert.IsType(t, ConsolePersistence{}, pArtifact.persistence)
		}
	})

	t.Run("should return persistent artifact with writeTo as Disk if save is enabled and path is not -", func(t *testing.T) {
		mapFs := afero.NewMemMapFs()
		artifact := NewArtifact(mapFs, true, "temp")

		if assert.IsType(t, PersistentArtifact{}, artifact) {
			pArtifact, _ := artifact.(PersistentArtifact)
			assert.IsType(t, DiskPersistence{}, pArtifact.persistence)
		}
	})
}

func TestArtifactSave(t *testing.T) {
	info := Info{
		ManifestFiles: stevedore.ManifestFiles{
			{File: "a.yaml",
				Manifest: stevedore.Manifest{
					DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "staging"}},
					Spec: stevedore.ReleaseSpecifications{
						{
							Release: stevedore.Release{Name: "component-a", Values: stevedore.Values{}},
							Configs: stevedore.Configs{"store": []interface{}{map[interface{}]interface{}{}}},
							Mounts:  stevedore.Configs(nil),
						},
						{
							Release: stevedore.Release{Name: "component-b", Values: stevedore.Values{}},
							Configs: stevedore.Configs{"store": []interface{}{map[interface{}]interface{}{}}},
							Mounts:  stevedore.Configs(nil),
						},
					},
				},
			},
			{File: "b.yaml",
				Manifest: stevedore.Manifest{
					DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "staging"}},
					Spec: stevedore.ReleaseSpecifications{
						{
							Release: stevedore.Release{Name: "component-c", Values: stevedore.Values{"value": "value1"}},
							Configs: stevedore.Configs{"store": []interface{}{map[interface{}]interface{}{}}},
							Mounts:  stevedore.Configs(nil),
						},
						{
							Release: stevedore.Release{Name: "component-d", Values: stevedore.Values{}},
							Configs: stevedore.Configs{"store": []interface{}{map[interface{}]interface{}{}}},
							Mounts:  stevedore.Configs(nil),
						},
					},
				},
			},
		},
	}

	expectedComponentAManifest := stevedore.Manifest{
		Kind:     stevedore.KindStevedoreManifest,
		Version:  stevedore.ManifestCurrentVersion,
		DeployTo: info.ManifestFiles[0].Manifest.DeployTo,
		Spec:     stevedore.ReleaseSpecifications{info.ManifestFiles[0].Manifest.Spec[0]},
	}
	expectedComponentBManifest := stevedore.Manifest{
		Kind:     stevedore.KindStevedoreManifest,
		Version:  stevedore.ManifestCurrentVersion,
		DeployTo: info.ManifestFiles[0].Manifest.DeployTo,
		Spec:     stevedore.ReleaseSpecifications{info.ManifestFiles[0].Manifest.Spec[1]},
	}
	expectedComponentCManifest := stevedore.Manifest{
		Kind:     stevedore.KindStevedoreManifest,
		Version:  stevedore.ManifestCurrentVersion,
		DeployTo: info.ManifestFiles[1].Manifest.DeployTo,
		Spec:     stevedore.ReleaseSpecifications{info.ManifestFiles[1].Manifest.Spec[0]},
	}
	expectedComponentDManifest := stevedore.Manifest{
		Kind:     stevedore.KindStevedoreManifest,
		Version:  stevedore.ManifestCurrentVersion,
		DeployTo: info.ManifestFiles[1].Manifest.DeployTo,
		Spec:     stevedore.ReleaseSpecifications{info.ManifestFiles[1].Manifest.Spec[1]},
	}

	t.Run("should write manifests to files", func(t *testing.T) {

		fs := afero.NewMemMapFs()
		artifactDir := "temp"
		artifact := NewArtifact(fs, true, artifactDir)

		err := artifact.Save(info)

		assert.Nil(t, err)

		componentAYamlData, err := afero.ReadFile(fs, filepath.Join(artifactDir, "component-a.yaml"))
		assert.NoError(t, err)
		assertComponent(t, expectedComponentAManifest, componentAYamlData)

		componentBYamlData, err := afero.ReadFile(fs, filepath.Join(artifactDir, "component-b.yaml"))
		assert.NoError(t, err)
		assertComponent(t, expectedComponentBManifest, componentBYamlData)

		componentCYamlData, err := afero.ReadFile(fs, filepath.Join(artifactDir, "component-c.yaml"))
		assert.NoError(t, err)
		assertComponent(t, expectedComponentCManifest, componentCYamlData)

		componentDYamlData, err := afero.ReadFile(fs, filepath.Join(artifactDir, "component-d.yaml"))
		assert.NoError(t, err)
		assertComponent(t, expectedComponentDManifest, componentDYamlData)
	})

}

func assertComponent(t *testing.T, expectedManifest stevedore.Manifest, componentYamlData []byte) {
	componentManifest := stevedore.Manifest{}
	err := yaml.Unmarshal(componentYamlData, &componentManifest)
	assert.NoError(t, err)
	assert.Equal(t, expectedManifest, componentManifest)
}
