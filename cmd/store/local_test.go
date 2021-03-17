package store

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	t.Run("should fetch envs", func(t *testing.T) {
		_ = os.Setenv("TEST_STEVEDORE_KEY_1", "value1")
		_ = os.Setenv("TEST_STEVEDORE_KEY_2", "value2")
		defer func() {
			_ = os.Unsetenv("TEST_STEVEDORE_KEY_1")
			_ = os.Unsetenv("TEST_STEVEDORE_KEY_2")
		}()

		local := Local{}
		envs := local.Fetch()

		assert.Equal(t, envs["TEST_STEVEDORE_KEY_1"], "value1")
		assert.Equal(t, envs["TEST_STEVEDORE_KEY_2"], "value2")
	})

	t.Run("should fetch envs containing = in value", func(t *testing.T) {
		err := os.Setenv("TEST_STEVEDORE_KEY", "val=ue=")
		if err != nil {
			t.Fatalf("Failed to set env %v", err)
		}
		defer func() {
			_ = os.Unsetenv("TEST_STEVEDORE_KEY")
			if err != nil {
				t.Fatalf("Failed to unset env %v", err)
			}
		}()

		local := Local{}
		envs := local.Fetch()

		assert.Equal(t, envs["TEST_STEVEDORE_KEY"], "val=ue=")
	})
}
