package stevedore

import (
	"bytes"
	"fmt"
	"sort"
)

// Context is the Configuration context to which operation will be performed against.
// It wraps information needed at a kubernetes cluster level
type Context struct {
	Name              string            `yaml:"name" validate:"required"`
	KubernetesContext string            `yaml:"kubernetesContext" validate:"required"`
	KubeConfigFile    string            `yaml:"kubeConfigFile"`
	Labels            map[string]string `yaml:"labels,omitempty"`
}

// IsValid validates the context and returns error if any
func (ctx Context) IsValid() error {
	return Validate(ctx)
}

// Map converts Context to a map[string]string
func (ctx Context) Map() (map[string]string, error) {
	result := ctx.Conditions()
	return result, nil
}

// String prints the details of context
func (ctx Context) String() string {
	buff := bytes.NewBufferString("\nContext Details:")
	buff.WriteString("\n------------------")
	buff.WriteString(fmt.Sprintf("\nName: %s", ctx.Name))
	buff.WriteString(fmt.Sprintf("\nKubernetes Context: %s", ctx.KubernetesContext))
	buff.WriteString(fmt.Sprintf("\nKubeConfig File: %s", ctx.KubeConfigFile))
	buff.WriteString("\nLabels:")
	keys := make([]string, 0, len(ctx.Labels))
	for key := range ctx.Labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := ctx.Labels[key]
		buff.WriteString(fmt.Sprintf("\n  %s: %s", key, value))
	}
	buff.WriteString("\n------------------")
	return buff.String()
}

// Conditions returns Conditions
func (ctx Context) Conditions() Conditions {
	conditions := Conditions{}
	conditions[ConditionContextName] = ctx.Name
	for key, value := range ctx.Labels {
		conditions[key] = value
	}
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
