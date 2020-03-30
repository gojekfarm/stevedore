package file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorError(t *testing.T) {
	t.Run("should have formatted message", func(t *testing.T) {
		err := Error{Reason: fmt.Errorf("reason provided from test"), Filename: "x-stevedore"}

		assert.Equal(t, "Error while processing x-stevedore\nreason provided from test", err.Error())
	})
}

func TestErrorsError(t *testing.T) {
	t.Run("should have formatted message", func(t *testing.T) {
		err := Errors{
			Error{Reason: fmt.Errorf("reason provided from test"), Filename: "x-stevedore"},
			Error{Reason: fmt.Errorf("something wrong with the file"), Filename: "y-stevedore"},
		}

		assert.Equal(t, "Failed with 2 Errors\n\n1. Error while processing x-stevedore\nreason provided from test\n\n2. Error while processing y-stevedore\nsomething wrong with the file\n", err.Error())
	})
}
