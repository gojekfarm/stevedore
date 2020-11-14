package provider_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDefaultOverrideProviderOverrides(t *testing.T) {
	t.Run("should return the overrides", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		overrideFileName := "/mock/overrides/override.yaml"
		overrideString := `
kind: StevedoreOverride
version: 2
spec:
- matches:
    environmentType: staging
    contextType: components
  values:
    db:
      persistence:
        size: 10Gi`
		anotherOverrideFileName := "/mock/overrides/anotherOverride.yaml"
		anotherOverrideString := `
kind: StevedoreOverride
version: 2
spec:
- matches:
    contextName: env-components-staging
    applicationName: abc-service
  values:
    db:
      persistence:
        size: 20Gi`

		err := afero.WriteFile(memFs, overrideFileName, []byte(overrideString), 0644)
		assert.NoError(t, err)

		err = afero.WriteFile(memFs, anotherOverrideFileName, []byte(anotherOverrideString), 0644)
		assert.NoError(t, err)

		overrideProvider := provider.NewOverrideProvider(memFs, "/mock/overrides")
		expected := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				stevedore.OverrideSpecification{
					FileName: "/mock/overrides/anotherOverride.yaml",
					Matches:  stevedore.Conditions{"contextName": "env-components-staging", "applicationName": "abc-service"},
					Values:   map[string]interface{}{"db": map[interface{}]interface{}{"persistence": map[interface{}]interface{}{"size": "20Gi"}}},
				},
				stevedore.OverrideSpecification{
					FileName: "/mock/overrides/override.yaml",
					Matches:  stevedore.Conditions{"environmentType": "staging", "contextType": "components"},
					Values:   map[string]interface{}{"db": map[interface{}]interface{}{"persistence": map[interface{}]interface{}{"size": "10Gi"}}},
				},
			},
		}

		overrides, err := overrideProvider.Overrides()

		assert.NoError(t, err)
		assert.Equal(t, expected.Kind, overrides.Kind)
		assert.Equal(t, expected.Version, overrides.Version)
		if !cmp.Equal(expected.Spec, overrides.Spec) {
			assert.Fail(t, cmp.Diff(expected.Spec, overrides.Spec))
		}
	})

	t.Run("it should not return error if the override files doesn't exists", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		_ = memFs.Mkdir("/mock/configs", 0644)

		overrideProvider := provider.NewOverrideProvider(memFs, "/mock/configs")

		overrides, err := overrideProvider.Overrides()

		assert.NoError(t, err)
		assert.Equal(t, stevedore.Overrides{Kind: stevedore.KindStevedoreOverride, Version: stevedore.OverrideCurrentVersion}, overrides)
	})
}
