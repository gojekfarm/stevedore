package stevedore_test

import (
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestSubstitutePerform(t *testing.T) {
	t.Run("should replace all placeholders", func(t *testing.T) {
		substitute := stevedore.Substitute{
			"NAME":   "x-service",
			"VALUE":  "value1",
			"NUMBER": 18121078944,
			"BOOL":   true,
			"FLOAT":  1.8121078944,
		}
		expected := "name is 'x-service' and age is $UNTOUCHED and value is 'value1' and number is 18121078944 and bool is true and float is 1.8121078944"

		actual, err := substitute.Perform("name is ${NAME} and age is $UNTOUCHED and value is ${VALUE} and number is ${NUMBER} and bool is ${BOOL} and float is ${FLOAT}")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("should not replace placeholders", func(t *testing.T) {
		input := "name is ${NAME}, type is ${TYPE} with url ${URL}"
		substitute := stevedore.Substitute{"NAME": "x-service"}

		actual, err := substitute.Perform(input)

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 2 variable(s):\n\t1. ${TYPE}\n\t2. ${URL}", err.Error())
		}

		assert.NotNil(t, actual)
		assert.Equal(t, input, actual)
	})

	t.Run("should not replace placeholders with $variable", func(t *testing.T) {
		input := "name is ${NAME}, type is $TYPE with url ${URL}"
		substitute := stevedore.Substitute{"NAME": "x-service"}

		actual, err := substitute.Perform(input)

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 1 variable(s):\n\t1. ${URL}", err.Error())
		}

		assert.NotNil(t, actual)
		assert.Equal(t, input, actual)
	})
}

func TestSubstituteMerge(t *testing.T) {
	t.Run("Should merge the given Substitute in-place", func(t *testing.T) {
		substitute := stevedore.Substitute{"NAME": "x-service"}
		expected := stevedore.Substitute{"NAME": "y-service", "TYPE": "worker", "ENV": "staging"}

		result, err := substitute.Merge(stevedore.Substitute{"TYPE": "server"}, stevedore.Substitute{"TYPE": "worker", "NAME": "y-service"}, stevedore.Substitute{"ENV": "staging"})

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
		assert.NotEqual(t, substitute, result)
	})
}
