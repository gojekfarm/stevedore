package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionsWeight(t *testing.T) {
	t.Run("should return weight for the condition", func(t *testing.T) {
		conditions := Conditions{
			"environmentType":        "environment-type",
			"contextType":            "context-type",
			"environment":            "environment",
			ConditionApplicationName: "application-name",
			ConditionContextName:     "context-name",
		}

		labels := Labels{
			{Name: "environmentType"},
			{Name: "environment"},
			{Name: "contextType"},
			{Name: "contextName"},
			{Name: "applicationName"},
		}
		weight := conditions.Weight(labels)

		assert.Equal(t, 31, weight)
	})
}

func TestConditionsConvert(t *testing.T) {
	t.Run("should convert the condition based on the context", func(t *testing.T) {
		conditions := Conditions{
			"environmentType":        "environment-type",
			"contextType":            "context-type",
			"environment":            "environment",
			ConditionApplicationName: "application-name",
			ConditionContextName:     "context-name",
		}

		expected := Conditions{
			"environmentType":        "pre-production",
			"contextType":            "readonly",
			"environment":            "staging",
			ConditionApplicationName: "application-name",
			ConditionContextName:     "new-context",
		}

		actual := conditions.Convert(Context{
			Name: "new-context",
			Labels: Conditions{
				"environment":     "staging",
				"environmentType": "pre-production",
				"contextType":     "readonly",
			},
		})

		assert.Equal(t, expected, actual)
	})
}
