package helm_test

import (
	"testing"

	"github.com/databus23/helm-diff/manifest"
	"github.com/gojek/stevedore/pkg/helm"
	"github.com/stretchr/testify/assert"
)

func TestUpstallResponseGenerateSummary(t *testing.T) {
	t.Run("should generate summary", func(t *testing.T) {
		existingSpecs := map[string]*manifest.MappingResult{
			"app, app-redis, Service (v1)": {
				Name:    "app, app-redis, Service (v1)",
				Kind:    "Service",
				Content: "manifest content for the redis",
			},
			"app, app-db-keeper, StatefulSet (apps)": {
				Name:    "app, app-db-keeper, StatefulSet (apps)",
				Kind:    "StatefulSet",
				Content: "manifest content for the db",
			},
			"app, app-cm, ConfigMap (v1)": {
				Name:    "app, app-cm, ConfigMap (v1)",
				Kind:    "ConfigMap",
				Content: "app cm",
			},
		}
		newSpecs := map[string]*manifest.MappingResult{
			"app, app-redis, Service (v1)": {
				Name:    "app, app-redis, Service (v1)",
				Kind:    "Service",
				Content: "manifest content for the redis",
			},
			"app, app-db-keeper, StatefulSet (apps)": {
				Name:    "app, app-db-keeper, StatefulSet (apps)",
				Kind:    "StatefulSet",
				Content: "updated manifest content for the db",
			},
			"app, app-db-exporter, ConfigMap (v1)": {
				Name:    "app, app-db-exporter, ConfigMap (v1)",
				Kind:    "ConfigMap",
				Content: "manifest content for the db",
			},
		}
		expectedSummary := helm.Summary{
			Added: helm.Resources{
				{Name: "app-db-exporter", Kind: "ConfigMap"},
			},
			Modified: helm.Resources{
				{Name: "app-db-keeper", Kind: "StatefulSet"},
			},
			Destroyed: helm.Resources{
				{Name: "app-cm", Kind: "ConfigMap"},
			},
		}
		response := helm.UpstallResponse{ExistingSpecs: existingSpecs, NewSpecs: newSpecs}

		summary := response.Summary()

		assert.Equal(t, expectedSummary, summary)
	})

	t.Run("Summarize Additions", func(t *testing.T) {
		t.Run("should show all resources as Added for new helm release", func(t *testing.T) {
			expectedSummary := helm.Summary{
				Added: helm.Resources{
					{Name: "app-redis", Kind: "Service"},
					{Name: "app-db-keeper", Kind: "StatefulSet"},
				},
				Modified:  helm.Resources{},
				Destroyed: helm.Resources{},
			}
			newSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
				"app, app-db-keeper, StatefulSet (apps)": {
					Name:    "app, app-db-keeper, StatefulSet (apps)",
					Kind:    "StatefulSet",
					Content: "manifest content for the db",
				},
			}
			response := helm.UpstallResponse{ExistingSpecs: nil, NewSpecs: newSpecs}

			summary := response.Summary()

			assert.ElementsMatch(t, expectedSummary.Added, summary.Added)
			assert.Empty(t, summary.Modified)
			assert.Empty(t, summary.Destroyed)
		})

		t.Run("should show newly added resources as Added when updating helm release", func(t *testing.T) {
			expectedSummary := helm.Summary{
				Added: helm.Resources{
					{Name: "app-db-keeper", Kind: "StatefulSet"},
				},
				Modified:  helm.Resources{},
				Destroyed: helm.Resources{},
			}
			newSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
				"app, app-db-keeper, StatefulSet (apps)": {
					Name:    "app, app-db-keeper, StatefulSet (apps)",
					Kind:    "StatefulSet",
					Content: "manifest content for the db",
				},
			}
			existingSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
			}
			response := helm.UpstallResponse{ExistingSpecs: existingSpecs, NewSpecs: newSpecs}

			summary := response.Summary()

			assert.ElementsMatch(t, expectedSummary.Added, summary.Added)
			assert.Empty(t, summary.Modified)
			assert.Empty(t, summary.Destroyed)
		})
	})

	t.Run("Summarize Modifications", func(t *testing.T) {
		t.Run("should show updated resources as Modified when updating helm release", func(t *testing.T) {
			existingSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
				"app, app-db-keeper, StatefulSet (apps)": {
					Name:    "app, app-db-keeper, StatefulSet (apps)",
					Kind:    "StatefulSet",
					Content: "manifest content for the db",
				},
				"app, app-db-exporter, ConfigMap (v1)": {
					Name:    "app, app-db-exporter, ConfigMap (v1)",
					Kind:    "ConfigMap",
					Content: "manifest content for the db",
				},
			}
			newSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
				"app, app-db-keeper, StatefulSet (apps)": {
					Name:    "app, app-db-keeper, StatefulSet (apps)",
					Kind:    "StatefulSet",
					Content: "updated manifest content for the db",
				},
				"app, app-db-exporter, ConfigMap (v1)": {
					Name:    "app, app-db-exporter, ConfigMap (v1)",
					Kind:    "ConfigMap",
					Content: "updated manifest content for the db",
				},
			}
			expectedSummary := helm.Summary{
				Added: helm.Resources{},
				Modified: helm.Resources{
					{Name: "app-db-keeper", Kind: "StatefulSet"},
					{Name: "app-db-exporter", Kind: "ConfigMap"},
				},
				Destroyed: helm.Resources{},
			}
			response := helm.UpstallResponse{ExistingSpecs: existingSpecs, NewSpecs: newSpecs}

			summary := response.Summary()

			assert.ElementsMatch(t, expectedSummary.Modified, summary.Modified)
			assert.Empty(t, summary.Added)
			assert.Empty(t, summary.Destroyed)
		})
	})

	t.Run("Summarize Deletions", func(t *testing.T) {
		t.Run("should show all resources as Destroyed on deleting helm release", func(t *testing.T) {
			expectedSummary := helm.Summary{
				Destroyed: helm.Resources{
					{Name: "app-redis", Kind: "Service"},
					{Name: "app-db-keeper", Kind: "StatefulSet"},
				},
				Modified: helm.Resources{},
				Added:    helm.Resources{},
			}
			existingSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
				"app, app-db-keeper, StatefulSet (apps)": {
					Name:    "app, app-db-keeper, StatefulSet (apps)",
					Kind:    "StatefulSet",
					Content: "manifest content for the db",
				},
			}
			response := helm.UpstallResponse{ExistingSpecs: existingSpecs, NewSpecs: nil}

			summary := response.Summary()

			assert.ElementsMatch(t, expectedSummary.Destroyed, summary.Destroyed)
			assert.Empty(t, summary.Modified)
			assert.Empty(t, summary.Added)
		})

		t.Run("should show deleted resources as Destroyed when updating helm release", func(t *testing.T) {
			expectedSummary := helm.Summary{
				Destroyed: helm.Resources{
					{Name: "app-db-keeper", Kind: "StatefulSet"},
				},
				Modified: helm.Resources{},
				Added:    helm.Resources{},
			}
			newSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
			}
			existingSpecs := map[string]*manifest.MappingResult{
				"app, app-redis, Service (v1)": {
					Name:    "app, app-redis, Service (v1)",
					Kind:    "Service",
					Content: "manifest content for the redis",
				},
				"app, app-db-keeper, StatefulSet (apps)": {
					Name:    "app, app-db-keeper, StatefulSet (apps)",
					Kind:    "StatefulSet",
					Content: "manifest content for the db",
				},
			}
			response := helm.UpstallResponse{NewSpecs: newSpecs, ExistingSpecs: existingSpecs}

			summary := response.Summary()

			assert.ElementsMatch(t, expectedSummary.Destroyed, summary.Destroyed)
			assert.Empty(t, summary.Modified)
			assert.Empty(t, summary.Added)
		})
	})
}
