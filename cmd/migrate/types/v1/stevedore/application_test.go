package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationReleaseSpecification(t *testing.T) {
	t.Run("should be able to convert to ReleaseSpecification", func(t *testing.T) {
		application := Application{
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
		}

		expected := stevedore.ReleaseSpecification{
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
		}

		actual := application.ReleaseSpecification()

		if !cmp.Equal(actual, expected, cmp.AllowUnexported(stevedore.Release{}, stevedore.ReleaseSpecification{})) {
			assert.Fail(t, cmp.Diff(actual, expected, cmp.AllowUnexported(stevedore.Release{}, stevedore.ReleaseSpecification{})))
		}
	})
}
