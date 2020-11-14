package maputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCommon(t *testing.T) {
	t.Run("giving same value should extract the same value", func(t *testing.T) {
		sample := map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}}
		actual := ExtractCommon(sample, sample)

		assert.Equal(t, sample, actual)
	})

	t.Run("giving different value should extract the common value", func(t *testing.T) {
		sampleOne := map[string]interface{}{"key1": map[string]interface{}{"key2": "value", "key3": "value"}}
		sampleTwo := map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}}
		actual := ExtractCommon(sampleOne, sampleTwo)

		assert.Equal(t, sampleTwo, actual)
	})

	t.Run("giving different value should extract the common value", func(t *testing.T) {
		sampleOne := map[string]interface{}{"key1": map[string]interface{}{"key2": "value", "key3": "value"}}
		sampleTwo := map[string]interface{}{"key1": map[string]interface{}{"key2": "value", "key3": "value3"}}

		expected := map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}}
		actual := ExtractCommon(sampleOne, sampleTwo)

		assert.Equal(t, expected, actual)
	})

}

func TestKeys(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	expectedOutput := []string{"key1", "key2"}

	assert.ElementsMatch(t, expectedOutput, Keys(input))
}
