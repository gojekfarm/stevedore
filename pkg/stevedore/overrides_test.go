package stevedore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gojek/stevedore/pkg/stevedore"
)

func TestNewOverrides(t *testing.T) {
	t.Run("should create overrides", func(t *testing.T) {
		overrideString := `
kind: StevedoreOverride
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
    applicationName: x-service
  values:
    db:
      persistence:
        size: 5Gi`

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

		actual, err := stevedore.NewOverrides(strings.NewReader(overrideString))

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)

	})

	t.Run("should not create overrides if matches are not valid", func(t *testing.T) {
		overrideString := `
kind: StevedoreOverride
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

		_, err := stevedore.NewOverrides(strings.NewReader(overrideString))

		if assert.NotNil(t, err) {
			assert.Equal(t, "Key: 'OverrideSpecification.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag", err.Error())
		}
	})

	t.Run("should not create overrides if kind is not provided", func(t *testing.T) {
		overrideString := `
kind:
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
    applicationName: x-service
  values:
    db:
      persistence:
        size: 5Gi`

		_, err := stevedore.NewOverrides(strings.NewReader(overrideString))

		if assert.NotNil(t, err) {
			assert.Equal(t, "Key: 'Overrides.Kind' Error:Field validation for 'Kind' failed on the 'required' tag", err.Error())
		}
	})

	t.Run("should not create overrides if version is not provided", func(t *testing.T) {
		overrideString := `
kind: StevedoreOverride
version:
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

		_, err := stevedore.NewOverrides(strings.NewReader(overrideString))

		if assert.NotNil(t, err) {
			assert.Equal(t, "Key: 'Overrides.Version' Error:Field validation for 'Version' failed on the 'required' tag", err.Error())
		}
	})

	t.Run("should return error if override.yaml is invalid", func(t *testing.T) {
		overrideString := `
kind: StevedoreOverride
version: 2
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

		_, err := stevedore.NewOverrides(strings.NewReader(overrideString))

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "[NewOverrides] error when validating from file")
	})
}

func TestOverrideFormat(t *testing.T) {
	t.Run("it should be able to formatAs as yaml", func(t *testing.T) {
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

		expected := `kind: StevedoreOverride
version: "2"
spec:
- matches:
    applicationName: x-service
    contextType: components
    environmentType: staging
  values:
    db:
      persistence:
        size: 10Gi
    redis:
      metrics:
        enabled: true
- matches:
    applicationName: x-service
    contextName: env-components-staging
  values:
    db:
      persistence:
        size: 5Gi
`
		actual := fmt.Sprintf("%y", overrides)

		assert.Equal(t, expected, actual)
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

		actual, err := stevedore.NewOverrides(strings.NewReader(fmt.Sprintf("%j", overrides)))
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, overrides, actual)
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
