package stringutils_test

import (
	"github.com/gojek/stevedore/pkg/utils/string"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("should return false if value does not exists", func(t *testing.T) {
		items := []string{"one", "two", "three"}

		exist := stringutils.Contains(items, "four")

		assert.False(t, exist)
	})

	t.Run("should return true if value exists", func(t *testing.T) {
		items := []string{"one", "two", "three"}

		exist := stringutils.Contains(items, "three")

		assert.True(t, exist)
	})
}

func TestUnique(t *testing.T) {
	t.Run("should return unique items from list", func(t *testing.T) {
		items := []string{"one", "two", "three", "one", "three", "four"}

		actual := stringutils.Unique(items)

		expected := []string{"one", "two", "three", "four"}
		assert.ElementsMatch(t, expected, actual)
	})
}

func TestExpand(t *testing.T) {
	t.Run("should expand", func(t *testing.T) {
		values := map[string]string{"firstName": "tom", "lastName": "jerry", "name": "james"}
		mappingFunction := func(key string, match bool) string {
			if !match {
				return values[key]
			}
			return key
		}
		expected := "hello tom jerry *name"

		actual := stringutils.Expand("hello ${firstName} ${lastName} ${*name", mappingFunction)

		assert.Equal(t, expected, actual)
	})
}
