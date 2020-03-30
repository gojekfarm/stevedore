package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoresFilter(t *testing.T) {

	rideServiceFullMatchIgnore := Ignore{
		Matches: Conditions{
			"applicationName": "app",
			"contextType":     "components",
			"contextName":     "components-production",
			"environmentType": "production",
			"environment":     "production",
		},
		Releases: IgnoredReleases{{Name: "app"}},
	}

	contextNameIgnore := Ignore{
		Matches: Conditions{
			"contextName": "components-production",
		},
		Releases: IgnoredReleases{{Name: "app"}},
	}

	contextNameAndAppNameIgnore := Ignore{
		Matches: Conditions{
			"contextName":     "components-production",
			"applicationName": "service-one",
		},
		Releases: IgnoredReleases{{Name: "service-one"}},
	}

	appServiceFullMatchIgnore := Ignore{
		Matches: Conditions{
			"applicationName": "app-service",
			"contextType":     "components",
			"contextName":     "components-production",
			"environmentType": "production",
			"environment":     "production",
		},
		Releases: IgnoredReleases{{Name: "app-service"}},
	}

	productionIgnore := Ignore{
		Matches: Conditions{
			"environmentType": "production",
		},
		Releases: IgnoredReleases{{Name: "service-x-service"}},
	}

	stagingIgnore := Ignore{
		Matches: Conditions{
			"environmentType": "staging",
		},
		Releases: IgnoredReleases{{Name: "service-a-api"}},
	}

	t.Run("should return manifest only if the release specification is not ignored", func(t *testing.T) {
		ignores := Ignores{
			rideServiceFullMatchIgnore,
			appServiceFullMatchIgnore,
			contextNameIgnore,
			contextNameAndAppNameIgnore,
			productionIgnore,
			stagingIgnore,
		}

		rideServiceApp := ReleaseSpecification{
			Release: Release{
				Name:      "app",
				Namespace: "ns",
				Chart:     "chart/app-dependencies",
			},
			Configs: Configs{
				"store": []map[string]interface{}{
					{"name": "app-server", "tags": []string{"server"}},
				},
			},
		}

		serviceThreeApp := ReleaseSpecification{
			Release: Release{
				Name:      "service-three",
				Namespace: "ns",
				Chart:     "chart/service-three-dependencies",
			},
			Configs: Configs{
				"store": []map[string]interface{}{
					{"name": "service-three-server", "tags": []string{"server"}},
				},
			},
		}

		conversationServiceApp := ReleaseSpecification{
			Release: Release{
				Name:      "conversation-service",
				Namespace: "ns",
				Chart:     "chart/conversation-service-dependencies",
			},
			Configs: Configs{
				"store": []map[string]interface{}{
					{"name": "conversation-service-server", "tags": []string{"server"}},
				},
			},
		}

		manifests := ManifestFiles{
			{
				File: "manifest1",
				Manifest: Manifest{
					DeployTo: Matchers{{ConditionContextName: "components-production"}},
					Spec:     ReleaseSpecifications{rideServiceApp},
				},
			},
			{
				File: "manifest2",
				Manifest: Manifest{
					DeployTo: Matchers{{ConditionContextName: "components-production"}},
					Spec:     ReleaseSpecifications{serviceThreeApp},
				},
			},
			{
				File: "manifest3",
				Manifest: Manifest{
					DeployTo: Matchers{{ConditionContextName: "staging"}},
					Spec:     ReleaseSpecifications{conversationServiceApp},
				},
			},
		}

		omponentsProductionCtx := Context{
			Name:              "components-production",
			KubernetesContext: "gke://components-production",
			Type:              "components",
			EnvironmentType:   "production",
			Environment:       "production",
		}

		expectedManifests := ManifestFiles{
			{
				File: "manifest2",
				Manifest: Manifest{
					DeployTo: Matchers{{ConditionContextName: "components-production"}},
					Spec:     ReleaseSpecifications{serviceThreeApp},
				},
			},
		}

		filteredManifests, ignoredComponents := manifests.Filter(ignores, omponentsProductionCtx)
		expectedIgnoredReleases := IgnoredReleases{
			IgnoredRelease{
				Name:   "conversation-service",
				Reason: "Not applicable for the context 'components-production'",
			},
			IgnoredRelease{
				Name: "app",
			},
		}

		assert.Equal(t, expectedManifests, filteredManifests)
		assert.Equal(t, expectedIgnoredReleases, ignoredComponents)
	})

	t.Run("should return manifest with the applications not ignored", func(t *testing.T) {
		ignores := Ignores{
			rideServiceFullMatchIgnore,
			appServiceFullMatchIgnore,
			contextNameIgnore,
			contextNameAndAppNameIgnore,
			productionIgnore,
			stagingIgnore,
		}

		rideServiceApp := ReleaseSpecification{
			Release: Release{
				Name:      "app",
				Namespace: "ns",
				Chart:     "chart/app-dependencies",
			},
			Configs: Configs{
				"store": []map[string]interface{}{
					{"name": "app-server", "tags": []string{"server"}},
				},
			},
		}
		serviceThreeApp := ReleaseSpecification{
			Release: Release{
				Name:      "service-three",
				Namespace: "ns",
				Chart:     "chart/service-three-dependencies",
			},
			Configs: Configs{
				"store": []map[string]interface{}{
					{"name": "service-three-server", "tags": []string{"server"}},
				},
			},
		}

		manifests := ManifestFiles{
			{
				File: "manifest1",
				Manifest: Manifest{
					DeployTo: Matchers{{ConditionContextName: "components-production"}},
					Spec: ReleaseSpecifications{
						rideServiceApp,
						serviceThreeApp,
					},
				},
			},
		}

		omponentsProductionCtx := Context{
			Name:              "components-production",
			KubernetesContext: "gke://components-production",
			Type:              "components",
			EnvironmentType:   "production",
			Environment:       "production",
		}

		expectedManifests := ManifestFiles{
			{
				File: "manifest1",
				Manifest: Manifest{
					DeployTo: Matchers{{ConditionContextName: "components-production"}},
					Spec:     ReleaseSpecifications{serviceThreeApp},
				},
			},
		}

		filteredManifests, ignoreNames := manifests.Filter(ignores, omponentsProductionCtx)
		assert.Equal(t, expectedManifests, filteredManifests)
		assert.Equal(t, IgnoredReleases{{Name: "app"}}, ignoreNames)
	})
}

func TestManifestFileNamespaces(t *testing.T) {
	t.Run("should get namespaces from a manifest file uniquely", func(t *testing.T) {
		manifestFile := ManifestFile{
			File: "manifest1",
			Manifest: Manifest{
				DeployTo: Matchers{{ConditionContextName: "components-production"}},
				Spec: ReleaseSpecifications{
					{
						Release: Release{
							Name:      "app",
							Namespace: "namespace-one",
						},
					},
					{
						Release: Release{
							Name:      "service-three",
							Namespace: "namespace-two",
						},
					},
					{
						Release: Release{
							Name:      "service-four",
							Namespace: "namespace-two",
						},
					},
				},
			},
		}

		actual := manifestFile.Namespaces()

		expected := []string{"namespace-one", "namespace-two"}
		assert.ElementsMatch(t, expected, actual)
	})
}

func TestManifestFileHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		manifestFile := ManifestFile{
			Manifest: Manifest{
				Spec: ReleaseSpecifications{
					ReleaseSpecification{
						Release: Release{ChartSpec: ChartSpec{Dependencies: Dependencies{{Name: "example"}}}},
					},
					ReleaseSpecification{
						Release: Release{ChartSpec: ChartSpec{}},
					},
				},
			},
		}

		ok := manifestFile.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		manifestFile := ManifestFile{
			Manifest: Manifest{
				Spec: ReleaseSpecifications{
					ReleaseSpecification{
						Release: Release{ChartSpec: ChartSpec{}},
					},
					ReleaseSpecification{
						Release: Release{ChartSpec: ChartSpec{}, Chart: "example-service"},
					},
				},
			},
		}

		ok := manifestFile.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		manifestFile := ManifestFile{
			Manifest: Manifest{
				Spec: ReleaseSpecifications{
					ReleaseSpecification{
						Release: Release{ChartSpec: ChartSpec{}, Chart: "example-service"},
					},
				},
			},
		}

		ok := manifestFile.HasBuildStep()
		assert.False(t, ok)
	})
}

func TestManifestFilesAllNamespaces(t *testing.T) {
	manifestFiles := ManifestFiles{
		{
			File: "nginx-a.yaml",
			Manifest: Manifest{
				DeployTo: Matchers{{ConditionContextName: "staging"}},
				Spec: ReleaseSpecifications{
					{
						Release: Release{
							Name:      "nginx-a-1",
							Namespace: "namespce-sample",
						},
					},
					{
						Release: Release{
							Name:      "nginx-a-2",
							Namespace: "default",
						},
					},
				},
			},
		},
		{
			File: "nginx-b.yaml",
			Manifest: Manifest{
				DeployTo: Matchers{{ConditionContextName: "staging"}},
				Spec: ReleaseSpecifications{
					{
						Release: Release{
							Name:      "nginx-b-1",
							Namespace: "namespace-four",
						},
					},
					{
						Release: Release{
							Name:       "nginx-b-2",
							Namespace:  "default",
							Privileged: true,
						},
					},
				},
			},
		},
	}

	expected := []string{"namespce-sample", "namespace-four", "default", "kube-system"}
	assert.ElementsMatch(t, expected, manifestFiles.AllNamespaces())
}

func TestManifestFilesHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		manifestFiles := ManifestFiles{
			ManifestFile{
				Manifest: Manifest{
					Spec: ReleaseSpecifications{
						ReleaseSpecification{
							Release: Release{ChartSpec: ChartSpec{Dependencies: Dependencies{{Name: "example"}}}},
						},
						ReleaseSpecification{
							Release: Release{ChartSpec: ChartSpec{}},
						},
					},
				},
			},
		}

		ok := manifestFiles.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		manifestFiles := ManifestFiles{
			ManifestFile{
				Manifest: Manifest{
					Spec: ReleaseSpecifications{
						ReleaseSpecification{
							Release: Release{ChartSpec: ChartSpec{}},
						},
						ReleaseSpecification{
							Release: Release{ChartSpec: ChartSpec{}, Chart: "example-service"},
						},
					},
				},
			},
		}

		ok := manifestFiles.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		manifestFiles := ManifestFiles{
			ManifestFile{
				Manifest: Manifest{
					Spec: ReleaseSpecifications{
						ReleaseSpecification{
							Release: Release{ChartSpec: ChartSpec{}, Chart: "example-service"},
						},
					},
				},
			},
		}

		ok := manifestFiles.HasBuildStep()
		assert.False(t, ok)
	})
}
