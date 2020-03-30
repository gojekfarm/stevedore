package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManifest_matchers(t *testing.T) {
	t.Run("should get the optimal contexts", func(t *testing.T) {
		contexts := stevedore.Contexts{
			{
				Name:            "one",
				Type:            "readonly",
				Environment:     "country-staging",
				EnvironmentType: "staging",
			},
			{
				Name:            "two",
				Type:            "readonly",
				Environment:     "country-staging",
				EnvironmentType: "staging",
			},
			{
				Name:            "three",
				Type:            "services",
				Environment:     "country-staging",
				EnvironmentType: "staging",
			},
			{
				Name:            "four",
				Type:            "readonly",
				Environment:     "another-country-staging",
				EnvironmentType: "staging",
			},
			{
				Name:            "five",
				Type:            "readonly",
				Environment:     "country-staging",
				EnvironmentType: "production",
			},
			{
				Name:            "six",
				Type:            "components",
				Environment:     "country-uat",
				EnvironmentType: "uat",
			},
			{
				Name:            "seven",
				Type:            "unknown",
				Environment:     "country-uat",
				EnvironmentType: "uat",
			},
		}
		manifest := Manifest{Environments: stevedore.Environments{"one", "two", "three", "four", "five", "six", "eight"}}
		expected := stevedore.Matchers{
			{stevedore.ConditionContextType: "readonly"},
			{stevedore.ConditionContextType: "services"},
			{stevedore.ConditionContextName: "six"},
			{stevedore.ConditionContextName: "eight"},
		}

		actual := manifest.matchers(contexts, true)
		if !assert.ElementsMatch(t, expected, actual) {
			assert.Fail(t, cmp.Diff(expected, actual))
		}
	})
}
