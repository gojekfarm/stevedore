package stevedore_test

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvs(t *testing.T) {
	t.Run("should create envs", func(t *testing.T) {
		envString := `
kind: StevedoreOverride
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
      applicationName: x-service
    env:
      size: 10Gi
  - matches:
      contextName: env-components-staging
      applicationName: x-service
    env:
      persistence: true`

		expected := stevedore.Env{
			Kind:    "StevedoreOverride",
			Version: "2",
			Spec: stevedore.EnvSpecifications{
				stevedore.EnvSpecification{
					Matches: stevedore.Conditions{
						"environmentType": "staging",
						"contextType":     "components",
						"applicationName": "x-service",
					},
					Values: stevedore.Substitute{"size": "10Gi"},
				},
				stevedore.EnvSpecification{
					Matches: stevedore.Conditions{
						"contextName":     "env-components-staging",
						"applicationName": "x-service",
					},
					Values: stevedore.Substitute{"persistence": true},
				},
			},
		}

		actual, err := stevedore.NewEnv(strings.NewReader(envString))

		assert.Nil(t, err)
		if !cmp.Equal(expected, actual) {
			assert.Fail(t, cmp.Diff(expected, actual))
		}
	})

	t.Run("should not create envs if matches are not valid", func(t *testing.T) {
		envString := `
kind: StevedoreEnv
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
      unknown: x-service
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

		actual, err := stevedore.NewEnv(strings.NewReader(envString))

		if assert.NotNil(t, err) {
			expected := "[NewEnv] error when validating from file:\nKey: 'Env.Spec[0].Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"
			if !cmp.Equal(expected, err.Error()) {
				assert.Fail(t, cmp.Diff(expected, err.Error()))
			}
		}
		assert.Equal(t, stevedore.Env{}, actual)
	})

	t.Run("should return error if env.yaml is invalid", func(t *testing.T) {
		envString := `
kind: StevedoreEnv
version: 2
spec:
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
    name: x-service
  values:
    db:
      persistence:
        size: 5Gi`

		actual, err := stevedore.NewEnv(strings.NewReader(envString))

		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "Key: 'Env.Spec[1].Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag")
		}
		assert.Equal(t, stevedore.Env{}, actual)

	})
}

func TestEnvIsValid(t *testing.T) {
	type scenario struct {
		name       string
		valid      bool
		conditions stevedore.Conditions
		error      string
	}

	scenarios := []scenario{
		{name: "should be valid if conditions are empty", valid: true, conditions: stevedore.Conditions{}},
		{name: "should be valid for environmentType", valid: true, conditions: stevedore.Conditions{"environmentType": "staging"}},
		{name: "should be valid for environment", valid: true, conditions: stevedore.Conditions{"environment": "staging"}},
		{name: "should be valid for contextType", valid: true, conditions: stevedore.Conditions{"contextType": "components"}},
		{name: "should be valid for contextName", valid: true, conditions: stevedore.Conditions{"contextName": "env-components"}},
		{name: "should be valid for applicationName", valid: true, conditions: stevedore.Conditions{"applicationName": "x-service"}},
		{name: "should be invalid for unknown", valid: false, conditions: stevedore.Conditions{"unknown": "invalid"}, error: "Key: 'Env.Spec[0].Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
		{name: "should be invalid even if one field is not valid", valid: false, conditions: stevedore.Conditions{"applicationName": "x-service", "unknown": "invalid"}, error: "Key: 'Env.Spec[0].Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("%v", scenario.name), func(t *testing.T) {
			env := stevedore.Env{
				Kind:    stevedore.KindStevedoreEnv,
				Version: stevedore.EnvCurrentVersion,
				Spec: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{Matches: scenario.conditions},
				},
			}

			err := env.IsValid()

			if scenario.valid {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Equal(t, scenario.error, err.Error())
				}
			}
		})
	}
}

func TestEnvSpecificationIsApplicableFor(t *testing.T) {

	t.Run("should return true", func(t *testing.T) {

		context := stevedore.Context{Name: "staging", Environment: "env", EnvironmentType: "staging"}

		env := stevedore.EnvSpecification{
			Matches: stevedore.Conditions{"environmentType": "staging"},
			Values:  stevedore.Substitute{"env": "staging"},
		}

		isApplicable := env.IsApplicableFor(context)

		assert.True(t, isApplicable)
	})

	t.Run("should return false", func(t *testing.T) {

		context := stevedore.Context{Name: "staging", Environment: "env", EnvironmentType: "production"}

		env := stevedore.EnvSpecification{
			Matches: stevedore.Conditions{"environmentType": "staging"},
			Values:  stevedore.Substitute{"env": "staging"},
		}

		isApplicable := env.IsApplicableFor(context)

		assert.False(t, isApplicable)
	})
}

func TestEnvSpecificationsSort(t *testing.T) {

	t.Run("should sort envs based on the pre-determined order", func(t *testing.T) {

		envs := stevedore.EnvSpecifications{
			{
				Matches: stevedore.Conditions{
					"contextName": "cluster",
				},
				Values: stevedore.Substitute{"env": "staging"},
			},
			{
				Matches: stevedore.Conditions{
					"contextType": "env",
				},
				Values: stevedore.Substitute{"env": "staging"},
			},
			{
				Matches: stevedore.Conditions{
					"contextType": "env",
					"environment": "staging",
				},
				Values: stevedore.Substitute{"env": "staging"},
			},
		}

		expectedEnvs := stevedore.EnvSpecifications{
			{
				Matches: stevedore.Conditions{
					"contextType": "env",
				},
				Values: stevedore.Substitute{"env": "staging"},
			},
			{
				Matches: stevedore.Conditions{
					"contextType": "env",
					"environment": "staging",
				},
				Values: stevedore.Substitute{"env": "staging"},
			},
			{
				Matches: stevedore.Conditions{
					"contextName": "cluster",
				},
				Values: stevedore.Substitute{"env": "staging"},
			},
		}

		envs.Sort()

		assert.Equal(t, expectedEnvs, envs)
	})
}

func TestEnvFormat(t *testing.T) {
	t.Run("should format as yaml", func(t *testing.T) {
		envString := `
kind: StevedoreOverride
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
      applicationName: x-service
    env:
      size: 10Gi
      persistence: true
  - matches:
      contextName: env-components-staging
      applicationName: x-service
    env:
      z-index: 0
      persistence: true`

		expected := `kind: StevedoreOverride
version: "2"
spec:
- matches:
    applicationName: x-service
    contextType: components
    environmentType: staging
  env:
    persistence: true
    size: 10Gi
- matches:
    applicationName: x-service
    contextName: env-components-staging
  env:
    persistence: true
    z-index: 0
`
		envs, err := stevedore.NewEnv(strings.NewReader(envString))
		assert.Nil(t, err)

		actual := fmt.Sprintf("%y", envs)
		assert.Equal(t, expected, actual)
	})

	t.Run("it should be able to formatAs as json", func(t *testing.T) {
		envString := `
kind: StevedoreOverride
version: 2
spec:
  - matches:
      environmentType: staging
      contextType: components
      applicationName: x-service
    env:
      size: 10Gi
      persistence: true
  - matches:
      contextName: env-components-staging
      applicationName: x-service
    env:
      z-index: 0
      persistence: true`

		env, err := stevedore.NewEnv(strings.NewReader(envString))
		if !assert.NoError(t, err) {
			return
		}

		actual, err := stevedore.NewEnv(strings.NewReader(fmt.Sprintf("%j", env)))
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, env, actual)
	})

	t.Run("it should be able to formatAs as json", func(t *testing.T) {
		overrideString := `
kind: StevedoreOverride
version: 2
spec:
- matches:
    environmentType: staging
    contextType: components
    applicationName: x-service
  values:
    redis:
      metrics:
        enabled: true
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

		overrides, err := stevedore.NewOverrides(strings.NewReader(overrideString))
		if !assert.NoError(t, err) {
			return
		}

		actual, err := stevedore.NewOverrides(strings.NewReader(fmt.Sprintf("%#j", overrides)))
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, overrides, actual)
	})
}
