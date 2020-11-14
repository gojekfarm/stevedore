package stevedore_test

import (
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentsContains(t *testing.T) {
	t.Run("should return true", func(t *testing.T) {
		environments := stevedore.Environments{"env", "staging", "production"}

		contains := environments.Contains("env")

		assert.True(t, contains)
	})

	t.Run("should return false", func(t *testing.T) {
		environments := stevedore.Environments{"env", "staging", "production"}

		contains := environments.Contains("uat")

		assert.False(t, contains)
	})
}

func TestNewEnvironments(t *testing.T) {
	t.Run("should return environments", func(t *testing.T) {
		environments := stevedore.NewEnvironments([]string{"production", "integration"})

		assert.Equal(t, stevedore.Environments{"production", "integration"}, environments)
	})
}
