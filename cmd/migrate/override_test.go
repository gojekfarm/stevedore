package migrate

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOverrideStrategyDo(t *testing.T) {
	t.Run("should convert older stevedore override file into newer one", func(t *testing.T) {
		override := `
- matches:
    environmentType: staging
    contextType: components
    applicationName: x-service
  values:
    db:
      persistence:
        size: 10Gi
- matches:
    contextName: env-components-staging
    applicationName: x-service
  values:
    db:
      persistence:
        size: 5Gi`

		memFs := afero.NewMemMapFs()
		overrideFilePath := "old/sample.yaml"
		_ = memFs.Mkdir("old", 0666)
		_ = afero.WriteFile(memFs, overrideFilePath, []byte(override), 0666)

		expected := stevedore.Overrides{
			Kind:    stevedore.KindStevedoreOverride,
			Version: stevedore.OverrideCurrentVersion,
			Spec: stevedore.OverrideSpecifications{
				{
					Matches: stevedore.Conditions{
						"environmentType": "staging",
						"contextType":     "components",
						"applicationName": "x-service",
					},
					Values: stevedore.Values{
						"db": map[interface{}]interface{}{
							"persistence": map[interface{}]interface{}{
								"size": "10Gi",
							},
						},
					},
				},
				{
					Matches: stevedore.Conditions{
						"contextName":     "env-components-staging",
						"applicationName": "x-service",
					},
					Values: stevedore.Values{
						"db": map[interface{}]interface{}{
							"persistence": map[interface{}]interface{}{
								"size": "5Gi",
							},
						},
					},
				},
			},
		}

		overrideStrategy := NewOverrideStrategy(memFs, overrideFilePath, false)

		err := overrideStrategy.Do()
		assert.NoError(t, err)

		actual := stevedore.Overrides{}
		err = read(memFs, overrideFilePath, &actual)
		assert.NoError(t, err)

		ignoreTypes := cmpopts.IgnoreUnexported()
		if !cmp.Equal(expected, actual, ignoreTypes) {
			assert.Fail(t, cmp.Diff(expected, actual, ignoreTypes))
		}
	})
}
