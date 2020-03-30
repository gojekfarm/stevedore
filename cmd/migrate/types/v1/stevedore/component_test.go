package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComponentRelease(t *testing.T) {
	t.Run("should be able to convert component to release", func(t *testing.T) {
		component := Component{
			Name:         "name",
			Namespace:    "namespace",
			Chart:        "chart",
			ChartVersion: "v1.0.0",
			ChartSpec: stevedore.ChartSpec{
				Name: "dependencies",
				Dependencies: stevedore.Dependencies{{
					Name:         "name",
					Alias:        "alias",
					Version:      "v1.0.1",
					Repository:   "repo",
					Condition:    "condition",
					Tags:         []string{},
					Enabled:      true,
					ImportValues: []interface{}{"one", 2},
				}},
			},
			CurrentReleaseVersion: 1,
			Values: stevedore.Values{
				"ONE": 1,
				"2":   "TWO",
			},
			Privileged: true,
		}

		expected := stevedore.Release{
			Name:         "name",
			Namespace:    "namespace",
			Chart:        "chart",
			ChartVersion: "v1.0.0",
			ChartSpec: stevedore.ChartSpec{
				Name: "dependencies",
				Dependencies: stevedore.Dependencies{{
					Name:         "name",
					Alias:        "alias",
					Version:      "v1.0.1",
					Repository:   "repo",
					Condition:    "condition",
					Tags:         []string{},
					Enabled:      true,
					ImportValues: []interface{}{"one", 2},
				}},
			},
			CurrentReleaseVersion: 1,
			Values: stevedore.Values{
				"ONE": 1,
				"2":   "TWO",
			},
			Privileged: true,
		}

		actual := component.Release()

		if !cmp.Equal(actual, expected, cmp.AllowUnexported(stevedore.Release{})) {
			assert.Fail(t, cmp.Diff(actual, expected, cmp.AllowUnexported(stevedore.Release{})))
		}
	})
}
