package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConditionWeight(t *testing.T) {
	t.Run("should return the correct weight for given knownCriteria", func(t *testing.T) {
		assert.Equal(t, 1, defaultConditionWeights.Sum([]string{"environmentType"}))
		assert.Equal(t, 2, defaultConditionWeights.Sum([]string{"environment"}))
		assert.Equal(t, 4, defaultConditionWeights.Sum([]string{"contextType"}))
		assert.Equal(t, 8, defaultConditionWeights.Sum([]string{"contextName"}))
		assert.Equal(t, 16, defaultConditionWeights.Sum([]string{"applicationName"}))
	})
}

func TestConditionsWeight(t *testing.T) {
	t.Run("should return weight for the condition", func(t *testing.T) {
		conditions := Conditions{
			ConditionEnvironmentType: "environment-type",
			ConditionContextType:     "context-type",
			ConditionEnvironment:     "environment",
			ConditionApplicationName: "application-name",
			ConditionContextName:     "context-name",
		}

		weight := conditions.Weight()

		assert.Equal(t, 31, weight)
	})
}

func TestConditionsConvert(t *testing.T) {
	t.Run("should convert the condition based on the context", func(t *testing.T) {
		conditions := Conditions{
			ConditionEnvironmentType: "environment-type",
			ConditionContextType:     "context-type",
			ConditionEnvironment:     "environment",
			ConditionApplicationName: "application-name",
			ConditionContextName:     "context-name",
		}

		expected := Conditions{
			ConditionEnvironmentType: "pre-production",
			ConditionContextType:     "readonly",
			ConditionEnvironment:     "staging",
			ConditionApplicationName: "application-name",
			ConditionContextName:     "new-context",
		}

		actual := conditions.Convert(Context{Environment: "staging", Name: "new-context", EnvironmentType: "pre-production", Type: "readonly"})

		assert.Equal(t, expected, actual)
	})
}
