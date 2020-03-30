package provider_test

import (
	"path/filepath"
	"testing"

	"github.com/gojek/stevedore/client/internal/mocks"
	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDefaultIgnoreProviderIgnores(t *testing.T) {
	currentWorkingDir := "mock"
	configsDir := filepath.Join(currentWorkingDir, "configs")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should return the ignores", func(t *testing.T) {
		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Cwd().Return(currentWorkingDir, nil)

		manifestIgnoreFilePath := filepath.Join(configsDir, ".stevedoreignore")
		manifestIgnoreFileContent := `
- matches:
    contextName: components-staging
  releases:
    - name: "x-stevedore"`

		cwdIgnoreFilePath := filepath.Join(currentWorkingDir, ".stevedoreignore")
		cwdIgnoreFileContent := `
- matches:
    contextName: components-staging
  releases:
    - name: "y-stevedore"`

		memFs := afero.NewMemMapFs()
		_ = afero.WriteFile(memFs, manifestIgnoreFilePath, []byte(manifestIgnoreFileContent), 0644)
		_ = afero.WriteFile(memFs, cwdIgnoreFilePath, []byte(cwdIgnoreFileContent), 0644)

		ignoreProvider, err := provider.NewIgnoreProvider(memFs, configsDir, mockEnvironment)
		assert.NoError(t, err)

		expected := stevedore.Ignores{
			stevedore.Ignore{
				Matches: stevedore.Conditions{"contextName": "components-staging"},
				Releases: stevedore.IgnoredReleases{
					stevedore.IgnoredRelease{Name: "x-stevedore"},
				},
			},
			stevedore.Ignore{
				Matches: stevedore.Conditions{"contextName": "components-staging"},
				Releases: stevedore.IgnoredReleases{
					stevedore.IgnoredRelease{Name: "y-stevedore"},
				},
			},
		}

		ignores, err := ignoreProvider.Ignores()

		assert.NoError(t, err)
		assert.Equal(t, expected, ignores)
	})

	t.Run("should return the ignores of the config directory", func(t *testing.T) {
		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Cwd().Return(currentWorkingDir, nil)

		manifestIgnoreFilePath := filepath.Join(configsDir, ".stevedoreignore")
		manifestIgnoreFileContent := `
- matches:
    contextName: components-staging
  releases:
    - name: "x-stevedore"`

		memFs := afero.NewMemMapFs()
		_ = afero.WriteFile(memFs, manifestIgnoreFilePath, []byte(manifestIgnoreFileContent), 0644)

		ignoreProvider, err := provider.NewIgnoreProvider(memFs, filepath.Join(configsDir, "x-stevedore.yaml"), mockEnvironment)
		assert.NoError(t, err)

		expected := stevedore.Ignores{
			stevedore.Ignore{
				Matches: stevedore.Conditions{"contextName": "components-staging"},
				Releases: stevedore.IgnoredReleases{
					stevedore.IgnoredRelease{Name: "x-stevedore"},
				},
			},
		}

		ignores, err := ignoreProvider.Ignores()

		assert.NoError(t, err)
		assert.Equal(t, expected, ignores)
	})

	t.Run("it should not return error if the ignore file doesn't exists", func(t *testing.T) {
		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Cwd().Return(currentWorkingDir, nil)

		memFs := afero.NewMemMapFs()
		_ = memFs.Mkdir(configsDir, 0644)

		ignoreProvider, err := provider.NewIgnoreProvider(memFs, configsDir, mockEnvironment)
		assert.NoError(t, err)

		ignores, err := ignoreProvider.Ignores()

		assert.Nil(t, err)
		assert.Equal(t, stevedore.Ignores{}, ignores)
	})
}
