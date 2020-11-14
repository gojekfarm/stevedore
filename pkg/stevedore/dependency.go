package stevedore

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sort"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
)

// Dependency represent a chart details
type Dependency struct {
	// Name is the name of the dependency.
	//
	// This must match the name in the dependency's Chart.yaml.
	// Required: true
	Name string `json:"name" yaml:"name" validate:"required"`
	// Alias to be used
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`
	// Version of the chart to be used
	// Required: true
	Version string `json:"version" yaml:"version" validate:"required"`
	// Helm repository URL where the chart is pulled from
	// Required: true
	Repository string `json:"repository" yaml:"repository" validate:"required"`
	// A yaml path that resolves to a boolean, used for enabling/disabling charts (e.g. subchart1.enabled )
	Condition string `json:"condition,omitempty" yaml:"condition,omitempty"`
	// Tags can be used to group charts for enabling/disabling together
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Enabled bool determines if chart should be loaded
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	// ImportValues holds the mapping of source values to parent key to be imported. Each item can be a
	// string or pair of child/parent sublist items.
	ImportValues []interface{} `json:"import-values,omitempty" yaml:"import-values,omitempty"`
}

// ChartUtilDependency converts dependency to chartutil.Dependency
func (dependency Dependency) ChartUtilDependency() chart.Dependency {
	return chart.Dependency{
		Name:         dependency.Name,
		Alias:        dependency.Alias,
		Version:      dependency.Version,
		Repository:   dependency.Repository,
		Condition:    dependency.Condition,
		Tags:         dependency.Tags,
		Enabled:      dependency.Enabled,
		ImportValues: dependency.ImportValues,
	}
}

// Dependencies is the collection of charts
type Dependencies []Dependency

// NewDependencies creates a list of stevedore dependencies from helm chart dependencies
func NewDependencies(chartDependencies []*chart.Dependency) Dependencies {
	dependencies := make(Dependencies, 0, len(chartDependencies))
	for _, chartDependency := range chartDependencies {
		dependency := Dependency{
			Name:         chartDependency.Name,
			Alias:        chartDependency.Alias,
			Version:      chartDependency.Version,
			Repository:   chartDependency.Repository,
			Condition:    chartDependency.Condition,
			Tags:         chartDependency.Tags,
			Enabled:      chartDependency.Enabled,
			ImportValues: chartDependency.ImportValues,
		}
		dependencies = append(dependencies, dependency)
	}
	return dependencies
}

// CheckSum will give the SHA256 based on the dependencies (sorted by alias)
func (dependencies Dependencies) CheckSum() (string, error) {
	sort.Slice(dependencies, func(i, j int) bool {
		if dependencies[i].Name == dependencies[j].Name {
			return dependencies[i].Alias < dependencies[j].Alias
		}
		return dependencies[i].Name < dependencies[j].Name
	})

	buffer := bytes.Buffer{}
	err := yaml.NewEncoder(&buffer).Encode(dependencies)
	if err != nil {
		return "", err
	}

	sum := fmt.Sprintf("%x", sha256.Sum256(buffer.Bytes()))
	return sum[:8], nil
}

// Contains will check and return dependencies matching chartName
func (dependencies Dependencies) Contains(chartName string) (Dependencies, bool) {
	matchedDependencies := Dependencies{}
	for _, dependency := range dependencies {
		if dependency.Name == chartName {
			matchedDependencies = append(matchedDependencies, dependency)
		}
	}
	if len(matchedDependencies) == 0 {
		return matchedDependencies, false
	}
	return matchedDependencies, true
}
