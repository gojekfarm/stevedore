package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchersContains(t *testing.T) {
	t.Run("should return true if context is available", func(t *testing.T) {
		matchers := Matchers{
			{
				ConditionContextName:     "one",
				ConditionEnvironment:     "country_staging",
				ConditionEnvironmentType: "staging",
				ConditionContextType:     "readonly",
			},
			{
				ConditionContextName:     "two",
				ConditionEnvironment:     "country_staging",
				ConditionEnvironmentType: "staging",
				ConditionContextType:     "readonly",
			},
		}
		context := Context{
			Name:              "one",
			Type:              "readonly",
			Environment:       "country_staging",
			EnvironmentType:   "staging",
			KubernetesContext: "kube-context-one",
		}

		result := matchers.Contains(context)

		assert.True(t, result)
	})

	t.Run("should return false if context is not available", func(t *testing.T) {
		matchers := Matchers{
			{
				ConditionContextName:     "two",
				ConditionEnvironment:     "country_staging",
				ConditionEnvironmentType: "staging",
				ConditionContextType:     "readonly",
			},
			{
				ConditionContextName:     "three",
				ConditionEnvironment:     "country_staging",
				ConditionEnvironmentType: "staging",
				ConditionContextType:     "readonly",
			},
		}
		context := Context{
			Name:              "one",
			Type:              "readonly",
			Environment:       "country_staging",
			EnvironmentType:   "staging",
			KubernetesContext: "kube-context-one",
		}

		result := matchers.Contains(context)

		assert.False(t, result)
	})
}
