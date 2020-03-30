package merger_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gojek/stevedore/pkg/merger"
)

func TestMerge(t *testing.T) {
	t.Run("should merge 2 hashes and overrides", func(t *testing.T) {
		aMap := map[string]interface{}{
			"key1": "valueA1",
			"key2": "valueA2",
		}

		bMap := map[string]interface{}{
			"key1": "valueB1",
			"key3": "valueB3",
		}

		expectedMap := map[string]interface{}{
			"key1": "valueB1",
			"key2": "valueA2",
			"key3": "valueB3",
		}
		finalMap, err := merger.Merge(aMap, bMap)

		assert.Nil(t, err)

		if !reflect.DeepEqual(finalMap, expectedMap) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", finalMap, expectedMap)
		}
	})

	t.Run("should merge 2 hashes but override array", func(t *testing.T) {
		aMap := map[string]interface{}{
			"key1": []int{1},
			"key2": "valueA2",
		}

		bMap := map[string]interface{}{
			"key1": []int{2, 3},
			"key2": "valueB2",
			"key3": "valueB3",
		}

		expectedMap := map[string]interface{}{
			"key1": []int{2, 3},
			"key2": "valueB2",
			"key3": "valueB3",
		}
		finalMap, err := merger.Merge(aMap, bMap)

		assert.Nil(t, err)

		if !reflect.DeepEqual(finalMap, expectedMap) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", finalMap, expectedMap)
		}
	})

	t.Run("should test deep merge", func(t *testing.T) {
		aMap := map[string]interface{}{
			"key1": map[string]interface{}{
				"nKey1": "nValueA1",
				"nKey3": []int{1},
			},
			"key2": "valueA2",
		}

		bMap := map[string]interface{}{
			"key1": map[string]interface{}{
				"nKey1": "nValueB1",
				"nKey2": "nValueB2",
			},
			"key2": "valueB2",
		}

		expectedMap := map[string]interface{}{
			"key1": map[string]interface{}{
				"nKey1": "nValueB1",
				"nKey2": "nValueB2",
				"nKey3": []int{1},
			},
			"key2": "valueB2",
		}
		finalMap, err := merger.Merge(aMap, bMap)

		assert.Nil(t, err)

		if !reflect.DeepEqual(finalMap, expectedMap) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", finalMap, expectedMap)
		}
	})

	t.Run("should merge 3 hashes", func(t *testing.T) {
		aMap := map[string]interface{}{
			"key1": "valueA1",
			"key2": "valueA2",
		}

		bMap := map[string]interface{}{
			"key1": "valueB1",
			"key3": "valueB3",
		}

		cMap := map[string]interface{}{
			"key3": "valueC3",
		}

		expectedMap := map[string]interface{}{
			"key1": "valueB1",
			"key2": "valueA2",
			"key3": "valueC3",
		}
		finalMap, err := merger.Merge(aMap, bMap, cMap)

		assert.Nil(t, err)

		if !reflect.DeepEqual(finalMap, expectedMap) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", finalMap, expectedMap)
		}
	})
}
