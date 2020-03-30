package migrate

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvStrategyDo(t *testing.T) {
	t.Run("should convert older stevedore env file into newer one", func(t *testing.T) {
		env := `
- matches:
    environmentType: staging
    contextType: components
    applicationName: x-service
  env:
    HOST: host-name
    NAME: x-service
- matches:
    contextName: env-components-staging
    applicationName: x-service
  env:
    HOST: another-host
    NAME: x-service`

		memFs := afero.NewMemMapFs()
		envFilePath := "old/sample.yaml"
		_ = memFs.Mkdir("old", 0666)
		_ = afero.WriteFile(memFs, envFilePath, []byte(env), 0666)

		expected := stevedore.Env{
			Kind:    stevedore.KindStevedoreEnv,
			Version: stevedore.EnvCurrentVersion,
			Spec: stevedore.EnvSpecifications{
				{
					Matches: stevedore.Conditions{
						"environmentType": "staging",
						"contextType":     "components",
						"applicationName": "x-service",
					},
					Values: stevedore.Substitute{
						"HOST": "host-name",
						"NAME": "x-service",
					},
				},
				{
					Matches: stevedore.Conditions{
						"contextName":     "env-components-staging",
						"applicationName": "x-service",
					},
					Values: stevedore.Substitute{
						"HOST": "another-host",
						"NAME": "x-service",
					},
				},
			},
		}

		envStrategy := NewEnvStrategy(memFs, envFilePath, false)

		err := envStrategy.Do()
		assert.NoError(t, err)

		actual := stevedore.Env{}
		err = read(memFs, envFilePath, &actual)
		assert.NoError(t, err)

		ignoreTypes := cmpopts.IgnoreUnexported()
		if !cmp.Equal(expected, actual, ignoreTypes) {
			assert.Fail(t, cmp.Diff(expected, actual, ignoreTypes))
		}
	})
}
