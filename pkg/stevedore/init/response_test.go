package init

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponsesString(t *testing.T) {
	t.Run("should return string representing underlying responses", func(t *testing.T) {
		responses := Responses{
			{Message: "message 1", Namespace: "namespace 1"},
			{Message: "message 2", Namespace: "namespace 2"},
		}
		expected := "Stevedore initialised in below namespace(s):\nnamespace 1: message 1\nnamespace 2: message 2\n"

		actual := responses.String()

		assert.Equal(t, expected, actual)
	})
}
