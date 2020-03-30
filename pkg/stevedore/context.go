package stevedore

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
)

// Context is the Configuration context to which operation will be performed against.
// It wraps information needed at a kubernetes cluster level
type Context struct {
	Name              string `yaml:"name" validate:"required"`
	Type              string `yaml:"type" validate:"required,any=services/components/readonly"`
	Environment       string `yaml:"environment" validate:"required"`
	KubernetesContext string `yaml:"kubernetesContext" validate:"required"`
	EnvironmentType   string `yaml:"environmentType" validate:"required"`
	KubeConfigFile    string `yaml:"kubeConfigFile"`
}

// IsValid validates the context and returns error if any
func (ctx Context) IsValid() error {
	return Validate(ctx)
}

// Map converts Context to a map[string]string
func (ctx Context) Map() (map[string]string, error) {
	data, err := yaml.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	mapToreturn := map[string]string{}
	if err := yaml.Unmarshal(data, &mapToreturn); err != nil {
		return nil, err
	}
	return mapToreturn, nil
}

// String prints the details of context
func (ctx Context) String() string {
	buff := bytes.NewBufferString("\nContext Details:")
	buff.WriteString("\n------------------")
	buff.WriteString(fmt.Sprintf("\nName: %s", ctx.Name))
	buff.WriteString(fmt.Sprintf("\nType: %s", ctx.Type))
	buff.WriteString(fmt.Sprintf("\nEnvironment: %s", ctx.Environment))
	buff.WriteString(fmt.Sprintf("\nKubernetes Context: %s", ctx.KubernetesContext))
	buff.WriteString(fmt.Sprintf("\nEnvironment Type: %s", ctx.EnvironmentType))
	buff.WriteString(fmt.Sprintf("\nKubeConfig File: %s", ctx.KubeConfigFile))
	buff.WriteString("\n------------------")
	return buff.String()
}

// Conditions returns Conditions
func (ctx Context) Conditions() Conditions {
	conditions := Conditions{}
	conditions[ConditionEnvironment] = ctx.Environment
	conditions[ConditionEnvironmentType] = ctx.EnvironmentType
	conditions[ConditionContextName] = ctx.Name
	conditions[ConditionContextType] = ctx.Type
	return conditions
}

// Contexts is a collection of context
type Contexts []Context

// Find returns the context if present with its presence as a boolean
func (c Contexts) Find(name string) (int, bool) {
	for i, ctx := range c {
		if ctx.Name == name {
			return i, true
		}
	}

	return -1, false
}
