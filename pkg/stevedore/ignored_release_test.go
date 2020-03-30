package stevedore_test

import (
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreReleaseFind(t *testing.T) {
	ignoredRideService := stevedore.IgnoredRelease{Name: "service-component", Reason: "temporarily ignored"}
	ignoredReleases := stevedore.IgnoredReleases{ignoredRideService, {Name: "app"}}

	t.Run("should return ignored release and true if found", func(t *testing.T) {
		ignoredReleases, exists := ignoredReleases.Find("service-component")

		assert.True(t, exists)
		assert.Equal(t, ignoredRideService, ignoredReleases)
	})

	t.Run("should return false if not found", func(t *testing.T) {
		_, exists := ignoredReleases.Find("this service doesnt exist")

		assert.False(t, exists)
	})
}

func TestIgnoreReleaseName(t *testing.T) {
	ignoredRideService := stevedore.IgnoredRelease{Name: "service-component", Reason: "temporarily ignored"}
	ignoredReleases := stevedore.IgnoredReleases{ignoredRideService, {Name: "app"}}

	t.Run("should all names from ignored components", func(t *testing.T) {
		ignoredReleases := ignoredReleases.Names()

		assert.Equal(t, ignoredReleases, []string{"service-component", "app"})
	})
}
