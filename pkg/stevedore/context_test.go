package stevedore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextsExists(t *testing.T) {
	t.Run("should return true if context name already exists", func(t *testing.T) {
		existingCtx := Context{Name: "components", KubernetesContext: "components", Environment: "env"}
		contexts := Contexts{
			Context{Name: "services", KubernetesContext: "services", Environment: "env"},
			existingCtx,
		}
		index, exists := contexts.Find("components")

		assert.True(t, exists)
		assert.Equal(t, 1, index)
	})

	t.Run("should return false if context name already exists", func(t *testing.T) {
		contexts := Contexts{
			Context{Name: "components", KubernetesContext: "components", Environment: "env"},
			Context{Name: "services", KubernetesContext: "services", Environment: "env"},
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
		valid        bool
	}

	emptyScenarios := []scenario{
		{
			name:         "name is not provided",
			context:      Context{KubernetesContext: "components", Environment: "env", EnvironmentType: "staging", Type: "components"},
			errorMessage: "Key: 'Context.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name:         "kubernetes context is not provided",
			context:      Context{Name: "components", Environment: "env", EnvironmentType: "staging", Type: "components"},
			errorMessage: "Key: 'Context.KubernetesContext' Error:Field validation for 'KubernetesContext' failed on the 'required' tag",
		},
		{
			name:         "environment is not provided",
			context:      Context{Name: "components", KubernetesContext: "components", EnvironmentType: "staging", Type: "components"},
			errorMessage: "Key: 'Context.Environment' Error:Field validation for 'Environment' failed on the 'required' tag",
		},
		{
			name:         "environment type is not provided",
			context:      Context{Name: "components", KubernetesContext: "components", Environment: "env", Type: "components"},
			errorMessage: "Key: 'Context.EnvironmentType' Error:Field validation for 'EnvironmentType' failed on the 'required' tag",
		},
		{
			name:         "type is not provided",
			context:      Context{Name: "components", KubernetesContext: "components", Environment: "env", EnvironmentType: "staging"},
			errorMessage: "Key: 'Context.Type' Error:Field validation for 'Type' failed on the 'required' tag",
		},
	}

	typeScenarios := []scenario{
		{
			name:         "should return error if type is not components|services|readonly",
			context:      Context{Name: "components", KubernetesContext: "components", Environment: "env", EnvironmentType: "staging", Type: "bla"},
			errorMessage: "Key: 'Context.Type' Error:Field validation for 'Type' failed on the 'any' tag",
			valid:        false,
		},
		{
			name:         "should not return error if type is components",
			context:      Context{Name: "components", KubernetesContext: "components", Environment: "env", EnvironmentType: "staging", Type: "components"},
			errorMessage: "Key: 'Context.Type' Error:Field validation for 'Type' failed on the 'any' tag",
			valid:        true,
		},
		{
			name:         "should not return error if type is services",
			context:      Context{Name: "components", KubernetesContext: "components", Environment: "env", EnvironmentType: "staging", Type: "services"},
			errorMessage: "Key: 'Context.Type' Error:Field validation for 'Type' failed on the 'any' tag",
			valid:        true,
		},
		{
			name:         "should not return error if type is readonly",
			context:      Context{Name: "components", KubernetesContext: "components", Environment: "env", EnvironmentType: "staging", Type: "readonly"},
			errorMessage: "Key: 'Context.Type' Error:Field validation for 'Type' failed on the 'any' tag",
			valid:        true,
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

	for _, s := range typeScenarios {
		t.Run(s.name, func(c *testing.T) {
			err := s.context.IsValid()

			if s.valid {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), s.errorMessage)
				}
			}
		})
	}
}

func TestContextString(t *testing.T) {
	t.Run("it should return context as string", func(t *testing.T) {
		context := Context{Name: "components", Type: "services", EnvironmentType: "staging", KubernetesContext: "components", Environment: "env", KubeConfigFile: "~/.kube/configs/minikube"}
		expected := `
Context Details:
------------------
Name: components
Type: services
Environment: env
Kubernetes Context: components
Environment Type: staging
KubeConfig File: ~/.kube/configs/minikube
------------------`

		content := context.String()

		assert.Equal(t, expected, content)
	})
}

func TestContextConditions(t *testing.T) {
	t.Run("it should return conditions", func(t *testing.T) {
		context := Context{Name: "components", Type: "services", EnvironmentType: "staging", KubernetesContext: "components", Environment: "env"}
		expected := Conditions{ConditionContextName: "components", ConditionContextType: "services", ConditionEnvironmentType: "staging", ConditionEnvironment: "env"}
		conditions := context.Conditions()

		assert.Equal(t, expected, conditions)
	})
}

func TestContextMap(t *testing.T) {
	context := Context{Name: "components", Type: "services", EnvironmentType: "staging", KubernetesContext: "components", Environment: "env"}
	expected := map[string]string{"name": "components", "type": "services", "environmentType": "staging", "kubernetesContext": "components", "environment": "env", "kubeConfigFile": ""}

	actual, err := context.Map()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
