package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationReleaseSpecifications(t *testing.T) {
	t.Run("should be able to convert to ReleaseSpecifications", func(t *testing.T) {
		applications := Applications{{
			Component: Component{
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
			},
			Configs: stevedore.Configs{
				"ONE": 1,
				"2":   "TWO",
			},
			DependsOn: []string{"1", "two"},
			Mounts: stevedore.Configs{
				"1":   "ONE",
				"TWO": 2,
			},
		}}

		expected := stevedore.ReleaseSpecifications{{
			Release: stevedore.Release{
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
			},
			Configs: stevedore.Configs{
				"ONE": 1,
				"2":   "TWO",
			},
			DependsOn: []string{"1", "two"},
			Mounts: stevedore.Configs{
				"1":   "ONE",
				"TWO": 2,
			},
		}}

		actual := applications.ReleaseSpecifications()

		if !cmp.Equal(actual, expected, cmp.AllowUnexported(stevedore.Release{}, stevedore.ReleaseSpecification{})) {
			assert.Fail(t, cmp.Diff(actual, expected, cmp.AllowUnexported(stevedore.Release{}, stevedore.ReleaseSpecification{})))
		}
	})
}
