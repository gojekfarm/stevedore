package stevedore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextsExists(t *testing.T) {
	t.Run("should return true if context name already exists", func(t *testing.T) {
		existingCtx := Context{Name: "components", KubernetesContext: "components"}
		contexts := Contexts{
			Context{Name: "services", KubernetesContext: "services"},
			existingCtx,
		}
		index, exists := contexts.Find("components")

		assert.True(t, exists)
		assert.Equal(t, 1, index)
	})

	t.Run("should return false if context name already exists", func(t *testing.T) {
		contexts := Contexts{
			Context{Name: "components", KubernetesContext: "components"},
			Context{Name: "services", KubernetesContext: "services"},
		}

		index, exists := contexts.Find("not-in-list")

		assert.False(t, exists)
		assert.Equal(t, -1, index)
	})
}

func TestContextValid(t *testing.T) {
	type scenario struct {
		context      Context
		errorMessage string
		name         string
	}

	emptyScenarios := []scenario{
		{
			name:         "name is not provided",
			context:      Context{KubernetesContext: "components"},
			errorMessage: "Key: 'Context.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name:         "kubernetes context is not provided",
			context:      Context{Name: "components"},
			errorMessage: "Key: 'Context.KubernetesContext' Error:Field validation for 'KubernetesContext' failed on the 'required' tag",
		},
	}

	for _, s := range emptyScenarios {
		t.Run(fmt.Sprintf("should return error if %s", s.name), func(t *testing.T) {
			err := s.context.IsValid()

			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), s.errorMessage)
			}
		})
	}
}

func TestContextString(t *testing.T) {
	t.Run("it should return context as string", func(t *testing.T) {
		context := Context{
			Name:              "components",
			KubernetesContext: "components",
			KubeConfigFile:    "~/.kube/configs/minikube",
			Labels: Conditions{
				"contextType":     "services",
				"environmentType": "staging",
				"environment":     "env",
			},
		}
		expected := `
Context Details:
------------------
Name: components
Kubernetes Context: components
KubeConfig File: ~/.kube/configs/minikube
Labels:
  contextType: services
  environment: env
  environmentType: staging
------------------`

		content := context.String()

		assert.Equal(t, expected, content)
	})
}

func TestContextConditions(t *testing.T) {
	t.Run("it should return conditions", func(t *testing.T) {
		context := Context{
			Name:              "components",
			KubernetesContext: "components",
			Labels: Conditions{
				"contextType":     "services",
				"environmentType": "staging",
				"environment":     "env",
			},
		}
		expected := Conditions{ConditionContextName: "components", "contextType": "services", "environmentType": "staging", "environment": "env"}
		conditions := context.Conditions()

		assert.Equal(t, expected, conditions)
	})
}

func TestContextMap(t *testing.T) {
	context := Context{
		Name: "components",
		Labels: Conditions{
			"type":              "services",
			"environmentType":   "staging",
			"kubernetesContext": "components",
			"environment":       "env",
		},
		KubernetesContext: "components",
	}
	expected := map[string]string{
		"contextName":       "components",
		"type":              "services",
		"environmentType":   "staging",
		"kubernetesContext": "components",
		"environment":       "env",
	}

	actual, err := context.Map()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
