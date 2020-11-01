package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabels_Weights(t *testing.T) {
	t.Run("should return the weights", func(t *testing.T) {
		labels := Labels{
			{Name: "one"},
			{Name: "two"},
			{Name: "three"},
		}

		expected := Weights{
			values: map[string]int{
				"one":   1,
				"two":   2,
				"three": 4,
			},
		}

		weights := labels.Weights()

		assert.NotEmpty(t, weights)
		assert.Equal(t, expected, weights)
	})

	t.Run("should sort labels based on weight", func(t *testing.T) {
		labels := Labels{
			{Name: "one", Weight: 3},
			{Name: "two", Weight: 2},
			{Name: "three", Weight: 1},
		}

		expected := Weights{
			values: map[string]int{
				"one":   4,
				"two":   2,
				"three": 1,
			},
		}

		weights := labels.Weights()

		assert.NotEmpty(t, weights)
		assert.Equal(t, expected, weights)
	})

	t.Run("should consider the array position if weight is not specified", func(t *testing.T) {
		labels := Labels{
			{Name: "one", Weight: 3},
			{Name: "two"},
			{Name: "three"},
			{Name: "four", Weight: 4},
		}

		expected := Weights{
			values: map[string]int{
				"one":   4,
				"two":   1,
				"three": 2,
				"four":  8,
			},
		}

		weights := labels.Weights()

		assert.NotEmpty(t, weights)
		assert.Equal(t, expected, weights)
	})
}
