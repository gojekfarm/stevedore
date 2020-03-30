package stevedore

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestComponentEnrichValues(t *testing.T) {
	t.Run("should return enriched release with final merged values given the overrides", func(t *testing.T) {
		release := Release{
			Name:       "x-stevedore",
			Namespace:  "ns",
			Chart:      "chart/x-stevedore-dependencies",
			Privileged: true,
			Values: Values{
				"key1": "baseValueForK1",
				"key2": "baseValueForK2",
			},
		}

		overrides := Overrides{
			Spec: OverrideSpecifications{
				{Values: Values{"key2": "firstOverrideForK2"}},
				{Values: Values{"key1": "secondOverrideForK1"}},
				{Values: Values{"key2": "thirdOverrideForK2"}},
				{Values: Values{"key3": "fourthOverrideForK3"}},
			},
		}

		expected := Release{
			Name:       "x-stevedore",
			Namespace:  "ns",
			Chart:      "chart/x-stevedore-dependencies",
			Privileged: true,
			Values: Values{
				"key1": "secondOverrideForK1",
				"key2": "thirdOverrideForK2",
				"key3": "fourthOverrideForK3",
			},
			usedSubstitute: Substitute{},
			overrides:      overrides,
		}

		actual := release.EnrichValues(overrides)

		equals := cmp.Equal(expected, actual, cmp.AllowUnexported(Release{}))
		if !equals {
			assert.Fail(t, cmp.Diff(expected, actual, cmp.AllowUnexported(Release{})))
		}
	})
}

func TestComponentOverrides(t *testing.T) {
	t.Run("should return overrides used for enrich", func(t *testing.T) {
		release := Release{
			Name:      "x-stevedore",
			Namespace: "ns",
			Chart:     "chart/x-stevedore-dependencies",
			Values: Values{
				"key1": "baseValueForK1",
				"key2": "baseValueForK2",
			},
		}

		expected := Overrides{
			Spec: OverrideSpecifications{
				{Values: Values{"key2": "firstOverrideForK2"}},
				{Values: Values{"key1": "secondOverrideForK1"}},
				{Values: Values{"key2": "thirdOverrideForK2"}},
				{Values: Values{"key3": "fourthOverrideForK3"}},
			},
		}

		actual := release.EnrichValues(expected)

		assert.Equal(t, expected, actual.Overrides())
	})
}

func TestComponentReplace(t *testing.T) {
	t.Run("should replace values", func(t *testing.T) {
		release := Release{
			Name:       "x-stevedore",
			Namespace:  "ns",
			Chart:      "chart/x-stevedore-dependencies",
			Privileged: true,
			Values: Values{
				"name": "${NAME}",
			},
		}

		expected := Release{
			Name:       "x-stevedore",
			Namespace:  "ns",
			Chart:      "chart/x-stevedore-dependencies",
			Privileged: true,
			Values: Values{
				"name": "x-service",
			},
		}

		actual, err := release.Replace(Substitute{"NAME": "x-service"})

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		if !cmp.Equal(expected, actual, cmpopts.IgnoreUnexported(Release{})) {
			assert.Fail(t, cmp.Diff(expected, actual, cmpopts.IgnoreUnexported(Release{})))
		}
	})

	t.Run("should fail when unable to replace values", func(t *testing.T) {
		release := Release{
			Name:      "x-stevedore",
			Namespace: "ns",
			Chart:     "chart/x-stevedore-dependencies",
			Values: Values{
				"name": "${NAME}",
				"game": "${GAME}",
			},
		}
		actual, err := release.Replace(Substitute{"FAME": "x-service"})

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 2 variable(s):\n\t1. ${GAME}\n\t2. ${NAME}", err.Error())
		}
		assert.NotNil(t, actual)
		assert.True(t, cmp.Equal(release, actual, cmpopts.IgnoreUnexported(Release{})))
	})
}

func TestComponentTillerNamespace(t *testing.T) {
	t.Run("should return namespace if not priviliged", func(t *testing.T) {
		release := Release{
			Namespace:  "ns",
			Privileged: false,
		}

		assert.Equal(t, "ns", release.TillerNamespace())
	})

	t.Run("should return kube-system if priviliged", func(t *testing.T) {
		release := Release{
			Namespace:  "ns",
			Privileged: true,
		}

		assert.Equal(t, "kube-system", release.TillerNamespace())
	})
}

func TestComponentHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		release := Release{ChartSpec: ChartSpec{Dependencies: Dependencies{{Name: "example"}}}}

		ok := release.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		release := Release{ChartSpec: ChartSpec{}}

		ok := release.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		release := Release{ChartSpec: ChartSpec{}, Chart: "example-service"}

		ok := release.HasBuildStep()

		assert.False(t, ok)
	})
}
