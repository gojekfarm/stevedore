package stevedore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gojek/stevedore/pkg/stevedore"
)

func TestNewIgnores(t *testing.T) {
	t.Run("should create ignores", func(t *testing.T) {
		ignoreString := `
- matches:
    environmentType: staging
    contextType: components
  releases:
    - name: app
- matches:
    contextName: env-components-staging
  releases:
    - name: app
      reason: temporarily ignoring
`

		expected := stevedore.Ignores{
			stevedore.Ignore{
				Matches: stevedore.Conditions{
					"environmentType": "staging",
					"contextType":     "components",
				},
				Releases: stevedore.IgnoredReleases{{Name: "app"}},
			},
			stevedore.Ignore{
				Matches: stevedore.Conditions{
					"contextName": "env-components-staging",
				},
				Releases: stevedore.IgnoredReleases{{Name: "app", Reason: "temporarily ignoring"}},
			},
		}

		actual, err := stevedore.NewIgnores(strings.NewReader(ignoreString))

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)

	})

	t.Run("should not create ignores if matches are not valid", func(t *testing.T) {
		ignoreString := `
- matches:
    environmentType: staging
    contextType: components
    unknown: x-service
  releases:
    - name: app
- matches:
    contextName: env-components-staging
  releases:
    - name: app`

		actual, err := stevedore.NewIgnores(strings.NewReader(ignoreString))

		if assert.NotNil(t, err) {
			assert.Equal(t, "Key: 'Ignore.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag", err.Error())
		}
		assert.Nil(t, actual)
	})

	t.Run("should return error if override.yaml is invalid", func(t *testing.T) {
		ignoreString := `
- matches:
    environmentType: staging
    contextType: components
    applicationName: x-service
  releases: ["app"]
- matches:
    contextName: env-components-staging
    applicationName: x-service
    releases: ["app"]`

		actual, err := stevedore.NewIgnores(strings.NewReader(ignoreString))

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "[NewIgnores] error when validating from file")
		assert.Nil(t, actual)

	})
}

func TestMatchedIgnores(t *testing.T) {

	rideServiceFullMatchIgnore := stevedore.Ignore{
		Matches: stevedore.Conditions{
			"applicationName": "app",
			"contextType":     "components",
			"contextName":     "env-components-production",
			"environmentType": "production",
			"environment":     "production",
		},
		Releases: stevedore.IgnoredReleases{{Name: "app"}},
	}

	contextNameIgnore := stevedore.Ignore{
		Matches: stevedore.Conditions{
			"contextName": "env-components-production",
		},
		Releases: stevedore.IgnoredReleases{{Name: "app"}},
	}

	contextNameAndAppNameIgnore := stevedore.Ignore{
		Matches: stevedore.Conditions{
			"contextName":     "env-components-production",
			"applicationName": "service-one",
		},
		Releases: stevedore.IgnoredReleases{{Name: "service-one"}},
	}

	appServiceFullMatchIgnore := stevedore.Ignore{
		Matches: stevedore.Conditions{
			"applicationName": "app-service",
			"contextType":     "components",
			"contextName":     "env-components-production",
			"environmentType": "production",
			"environment":     "production",
		},
		Releases: stevedore.IgnoredReleases{{Name: "app-service"}},
	}

	productionIgnore := stevedore.Ignore{
		Matches: stevedore.Conditions{
			"environmentType": "production",
		},
		Releases: stevedore.IgnoredReleases{{Name: "service-x"}},
	}

	stagingIgnore := stevedore.Ignore{
		Matches: stevedore.Conditions{
			"environmentType": "staging",
		},
		Releases: stevedore.IgnoredReleases{{Name: "service-one-api"}},
	}

	t.Run("should return the matched ignores", func(t *testing.T) {
		ignores := stevedore.Ignores{
			rideServiceFullMatchIgnore,
			appServiceFullMatchIgnore,
			contextNameIgnore,
			contextNameAndAppNameIgnore,
			productionIgnore,
			stagingIgnore,
		}

		rideServiceApp := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "app",
				Namespace: "ns",
				Chart:     "chart/app-dependencies",
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "app-server", "tags": []string{"server"}},
				},
			},
		}

		envComponentsProductionCtx := stevedore.Context{
			Name:              "env-components-production",
			KubernetesContext: "gke://env-components-production",
			Type:              "components",
			EnvironmentType:   "production",
			Environment:       "production",
		}

		predicate := stevedore.NewPredicate(rideServiceApp, envComponentsProductionCtx)

		expectedIgnores := stevedore.IgnoredReleases{{Name: "app"}, {Name: "app"}, {Name: "service-x"}}

		matchedIgnores := ignores.Filter(predicate)

		assert.Equal(t, expectedIgnores, matchedIgnores)
	})

	t.Run("should return empty if there are no matched ignores", func(t *testing.T) {
		ignores := stevedore.Ignores{
			rideServiceFullMatchIgnore,
			appServiceFullMatchIgnore,
			productionIgnore,
			stagingIgnore,
		}

		abcServiceApp := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "abc-service",
				Namespace: "ns",
				Chart:     "chart/abc-service-dependencies",
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "abc-service-server", "tags": []string{"server"}},
				},
			},
		}

		envComponentsProductionCtx := stevedore.Context{
			Name:              "env-components-production",
			KubernetesContext: "gke://env-components-production",
			Type:              "components",
			EnvironmentType:   "test-environment",
			Environment:       "production",
		}

		predicate := stevedore.NewPredicate(abcServiceApp, envComponentsProductionCtx)
		var expectedIgnores stevedore.IgnoredReleases

		matchedIgnores := ignores.Filter(predicate)

		assert.Equal(t, expectedIgnores, matchedIgnores)
	})
}
func TestIgnoreIsValid(t *testing.T) {
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
		{name: "should be invalid for unknown", valid: false, conditions: stevedore.Conditions{"unknown": "invalid"}, error: "Key: 'Ignore.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
		{name: "should be invalid even if one field is not valid", valid: false, conditions: stevedore.Conditions{"applicationName": "x-service", "unknown": "invalid"}, error: "Key: 'Ignore.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("%v", scenario.name), func(t *testing.T) {
			ignore := stevedore.Ignore{Matches: scenario.conditions}

			err := ignore.IsValid()

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
