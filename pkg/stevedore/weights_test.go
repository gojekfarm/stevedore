package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightsSum(t *testing.T) {
	t.Run("should assign 2^n as weight to each knownCriteria", func(t *testing.T) {
		weights := NewWeights([]string{"one", "two", "four"})

		assert.Equal(t, 1, weights.Sum([]string{"one"}))
		assert.Equal(t, 3, weights.Sum([]string{"one", "two"}))
		assert.Equal(t, 7, weights.Sum([]string{"four", "one", "two"}))
	})
}
