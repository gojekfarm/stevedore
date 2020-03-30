package stevedore

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvConvert(t *testing.T) {
	t.Run("should be able to convert", func(t *testing.T) {
		env := Env{
			Matches: stevedore.Conditions{
				"ONE": "1",
				"TWO": "2",
			},
			Values: stevedore.Substitute{
				"HOST": "host",
				"NAME": "name",
			},
		}

		expected := stevedore.EnvSpecification{
			Matches: stevedore.Conditions{
				"ONE": "1",
				"TWO": "2",
			},
			Values: stevedore.Substitute{
				"HOST": "host",
				"NAME": "name",
			},
		}

		actual := env.Convert()

		if !cmp.Equal(expected, actual) {
			assert.Fail(t, cmp.Diff(expected, actual))
		}
	})
}
