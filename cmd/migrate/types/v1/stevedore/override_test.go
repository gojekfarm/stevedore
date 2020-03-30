package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOverridesConvert(t *testing.T) {
	t.Run("should be able to convert", func(t *testing.T) {
		overrides := Overrides{
			{
				Matches: stevedore.Conditions{
					"one": "1",
					"two": "2",
				},
				Values: stevedore.Values{
					"HOST": "host",
					"NAME": "name",
				},
			},
		}

		actual := overrides.Convert()
		expected := stevedore.Overrides{
			Kind:    stevedore.KindStevedoreOverride,
			Version: stevedore.OverrideCurrentVersion,
			Spec: stevedore.OverrideSpecifications{
				{
					Matches: stevedore.Conditions{
						"one": "1",
						"two": "2",
					},
					Values: stevedore.Values{
						"HOST": "host",
						"NAME": "name",
					},
				},
			},
		}

		if !cmp.Equal(expected, actual) {
			assert.Fail(t, cmp.Diff(expected, actual))
		}
	})
}
