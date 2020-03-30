package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIgnoresConvert(t *testing.T) {
	t.Run("should be able to convert", func(t *testing.T) {
		ignores := Ignores{{
			Matches: stevedore.Conditions{
				"ONE": "1",
				"TWO": "2",
			},
			Components: stevedore.IgnoredReleases{
				{
					Name:   "x-service",
					Reason: "temporarily ignored",
				},
			},
		}}

		expected := stevedore.Ignores{{
			Matches: stevedore.Conditions{
				"ONE": "1",
				"TWO": "2",
			},
			Releases: stevedore.IgnoredReleases{
				{
					Name:   "x-service",
					Reason: "temporarily ignored",
				},
			},
		}}

		actual := ignores.Convert()

		if !cmp.Equal(expected, actual) {
			assert.Fail(t, cmp.Diff(expected, actual))
		}
	})
}
