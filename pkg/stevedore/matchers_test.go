package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchersContains(t *testing.T) {
	t.Run("should return true if context is available", func(t *testing.T) {
		matchers := Matchers{
			{
				ConditionContextName: "one",
				"environment":        "country_staging",
				"environmentType":    "staging",
				"contextType":        "readonly",
			},
			{
				ConditionContextName: "two",
				"environment":        "country_staging",
				"environmentType":    "staging",
				"contextType":        "readonly",
			},
		}
		context := Context{
			Name:              "one",
			KubernetesContext: "kube-context-one",
			Labels: Conditions{
				"contextType":     "readonly",
				"environment":     "country_staging",
				"environmentType": "staging",
			},
		}

		result := matchers.Contains(context)

		assert.True(t, result)
	})

	t.Run("should return false if context is not available", func(t *testing.T) {
		matchers := Matchers{
			{
				ConditionContextName: "two",
				"environment":        "country_staging",
				"environmentType":    "staging",
				"contextType":        "readonly",
			},
			{
				ConditionContextName: "three",
				"environment":        "country_staging",
				"environmentType":    "staging",
				"contextType":        "readonly",
			},
		}
		context := Context{
			Name:              "one",
			KubernetesContext: "kube-context-one",
			Labels: Conditions{
				"type":            "readonly",
				"environment":     "country_staging",
				"environmentType": "staging",
			},
		}

		result := matchers.Contains(context)

		assert.False(t, result)
	})
}
