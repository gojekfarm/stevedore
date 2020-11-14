package provider_test

import (
	"testing"

	"github.com/gojek/stevedore/client/provider"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDefaultEnvProviderIgnores(t *testing.T) {
	t.Run("should return the envs", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		envsFileName := "/mock/envs/env.yaml"
		envsString := `
kind: StevedoreEnv
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
    env:
      size: 10Gi
      value: 'value1'
  - matches:
      contextName: env-components-staging
    env:
      persistence: true`

		_ = afero.WriteFile(memFs, envsFileName, []byte(envsString), 0644)

		envProvider := provider.NewEnvProvider(memFs, "/mock/envs")
		expected := provider.EnvsFiles{
			provider.EnvsFile{
				Name: "/mock/envs/env.yaml",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{"environmentType": "staging", "contextType": "components"},
						Values:  stevedore.Substitute{"value": "value1", "size": "10Gi"},
					},
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{"contextName": "env-components-staging"},
						Values:  stevedore.Substitute{"persistence": true},
					},
				},
			},
		}

		envs, err := envProvider.Envs()

		assert.NoError(t, err)
		assert.Equal(t, expected, envs)
	})

	t.Run("should skip reading the file if kind is not provided", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		envsFileName := "/mock/envs/env.yaml"
		envsString := `
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
    env:
      size: 10Gi
      value: 'value1'
  - matches:
      contextName: env-components-staging
    env:
      persistence: true`

		_ = afero.WriteFile(memFs, envsFileName, []byte(envsString), 0644)

		envProvider := provider.NewEnvProvider(memFs, "/mock/envs")
		expected := provider.EnvsFiles{}

		envs, err := envProvider.Envs()

		assert.NoError(t, err)
		assert.Equal(t, expected, envs)
	})

	t.Run("should skip reading the file if kind is not StevedoreEnv", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		envsFileName := "/mock/envs/env.yaml"
		envsString := `
kind: StevedoreManifest
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
    env:
      size: 10Gi
      value: 'value1'
  - matches:
      contextName: env-components-staging
    env:
      persistence: true`

		_ = afero.WriteFile(memFs, envsFileName, []byte(envsString), 0644)

		envProvider := provider.NewEnvProvider(memFs, "/mock/envs")
		expected := provider.EnvsFiles{}

		envs, err := envProvider.Envs()

		assert.NoError(t, err)
		assert.Equal(t, expected, envs)
	})

	t.Run("should skip reading the file if kind is valid and version doesn't match", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		envsFileName := "/mock/envs/env.yaml"
		envsString := `
kind: StevedoreEnv
version: 3
spec:
  - matches:
      environmentType: staging
      contextType: components
    env:
      size: 10Gi
      value: 'value1'
  - matches:
      contextName: env-components-staging
    env:
      persistence: true`

		_ = afero.WriteFile(memFs, envsFileName, []byte(envsString), 0644)

		envProvider := provider.NewEnvProvider(memFs, "/mock/envs")
		expected := provider.EnvsFiles{}

		envs, err := envProvider.Envs()

		assert.NoError(t, err)
		assert.Equal(t, expected, envs)
	})

	t.Run("it should not return error if the env files doesn't exists", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		_ = memFs.Mkdir("/mock/configs", 0644)

		envProvider := provider.NewEnvProvider(memFs, "/mock/configs")

		envs, err := envProvider.Envs()

		assert.NoError(t, err)
		assert.Equal(t, provider.EnvsFiles{}, envs)
	})
}
