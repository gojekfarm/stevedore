package stevedore_test

import (
	"reflect"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestMatchedOverrideSpecifications(t *testing.T) {

	rideServiceFullMatchOverride := stevedore.OverrideSpecification{
		Matches: stevedore.Conditions{
			"applicationName": "app",
			"contextType":     "components",
			"contextName":     "components-production",
			"environmentType": "production",
			"environment":     "env_production",
		},
		Values: stevedore.Values{
			"key": "app-full-match",
		},
	}

	contextNameOverride := stevedore.OverrideSpecification{
		Matches: stevedore.Conditions{
			"contextName": "components-production",
		},
		Values: stevedore.Values{
			"key": "components-production",
		},
	}

	contextNameAndAppNameOverride := stevedore.OverrideSpecification{
		Matches: stevedore.Conditions{
			"contextName":     "components-production",
			"applicationName": "app",
		},
		Values: stevedore.Values{
			"key": "app-components-production",
		},
	}

	appServiceFullMatchOverride := stevedore.OverrideSpecification{
		Matches: stevedore.Conditions{
			"applicationName": "app-service",
			"contextType":     "components",
			"contextName":     "components-production",
			"environmentType": "production",
			"environment":     "env_production",
		},
		Values: stevedore.Values{
			"key": "app-service-full-match",
		},
	}

	productionOverride := stevedore.OverrideSpecification{
		Matches: stevedore.Conditions{
			"environmentType": "production",
		},
		Values: stevedore.Values{
			"key": "production",
		},
	}

	stagingOverride := stevedore.OverrideSpecification{
		Matches: stevedore.Conditions{
			"environmentType": "staging",
		},
		Values: stevedore.Values{
			"key": "staging",
		},
	}

	t.Run("should return the matched overrideSpecifications", func(t *testing.T) {
		overrideSpecifications := stevedore.OverrideSpecifications{
			rideServiceFullMatchOverride,
			appServiceFullMatchOverride,
			contextNameOverride,
			contextNameAndAppNameOverride,
			productionOverride,
			stagingOverride,
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
			Name:              "components-production",
			KubernetesContext: "gke://components-production",
			Labels: stevedore.Conditions{
				"contextType":     "components",
				"environmentType": "production",
				"environment":     "env_production",
			},
		}

		predicate := stevedore.NewPredicate(rideServiceApp, envComponentsProductionCtx)

		expectedOverrideSpecifications := stevedore.OverrideSpecifications{
			productionOverride,
			contextNameOverride,
			contextNameAndAppNameOverride,
			rideServiceFullMatchOverride,
		}

		labels := stevedore.Labels{
			{Name: "environmentType"},
			{Name: "environment"},
			{Name: "contextType"},
			{Name: "contextName"},
			{Name: "applicationName"},
		}

		matchedOverrideSpecifications := overrideSpecifications.CollateBy(predicate, labels)

		assert.Equal(t, expectedOverrideSpecifications, matchedOverrideSpecifications)
	})

	t.Run("should not return any matched overrideSpecifications", func(t *testing.T) {
		overrideSpecifications := stevedore.OverrideSpecifications{
			rideServiceFullMatchOverride,
			appServiceFullMatchOverride,
			productionOverride,
			stagingOverride,
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
			Name:              "components-production",
			KubernetesContext: "gke://components-production",
			Labels: stevedore.Conditions{
				"type":            "components",
				"environmentType": "test-environment",
				"environment":     "env_production",
			},
		}

		predicate := stevedore.NewPredicate(abcServiceApp, envComponentsProductionCtx)
		expectedOverrideSpecifications := stevedore.OverrideSpecifications{}

		matchedOverrideSpecifications := overrideSpecifications.CollateBy(predicate, stevedore.Labels{})

		assert.Equal(t, expectedOverrideSpecifications, matchedOverrideSpecifications)
	})
}

func TestOverrideSpecificationsMergeValuesInto(t *testing.T) {
	t.Run("should return overridden values merged with base", func(t *testing.T) {
		baseValues := stevedore.Values{
			"key1": "baseValueForK1",
			"key2": "baseValueForK2",
			"key3": "baseValueForK3",
		}

		overrideSpecifications := stevedore.OverrideSpecifications{
			{Values: stevedore.Values{"key2": "firstOverrideForK2"}},
			{Values: stevedore.Values{"key1": "secondOverrideForK1"}},
			{Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
		}

		expected := stevedore.Values{
			"key1": "secondOverrideForK1",
			"key2": "thirdOverrideForK2",
			"key3": "baseValueForK3",
		}

		actual := overrideSpecifications.MergeValuesInto(baseValues)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", actual, expected)
		}
	})
}
