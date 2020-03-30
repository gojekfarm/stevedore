package manifest_test

import (
	"testing"

	"github.com/gojek/stevedore/cmd/manifest"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestInfoFilterBy(t *testing.T) {
	t.Run("should return only the components that exist in responses", func(t *testing.T) {
		nginxBResult := stevedore.ManifestFile{
			File: "nginx-b.yaml",
			Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "production"}},
				Spec: stevedore.ReleaseSpecifications{
					{
						Release: stevedore.Release{
							Name:      "nginx-b-1",
							Namespace: "default",
						},
					},
					{
						Release: stevedore.Release{
							Name:      "nginx-b-2",
							Namespace: "default",
						},
					},
				},
			},
		}

		results := stevedore.ManifestFiles{
			{
				File: "nginx-a.yaml",
				Manifest: stevedore.Manifest{
					DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "staging"}},
					Spec: stevedore.ReleaseSpecifications{
						{
							Release: stevedore.Release{
								Name:      "nginx-a-1",
								Namespace: "default",
							},
						},
						{
							Release: stevedore.Release{
								Name:      "nginx-a-2",
								Namespace: "default",
							},
						},
					},
				},
			},
			nginxBResult,
			{
				File: "nginx-c.yaml",
				Manifest: stevedore.Manifest{
					DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "production"}},
					Spec: stevedore.ReleaseSpecifications{
						{
							Release: stevedore.Release{
								Name:      "nginx-c-1",
								Namespace: "namespace-sample",
							},
						},
					},
				},
			},
		}

		expected := manifest.Info{
			ManifestFiles: stevedore.ManifestFiles{
				{
					File: "nginx-a.yaml",
					Manifest: stevedore.Manifest{
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "staging"}},
						Spec: stevedore.ReleaseSpecifications{
							{
								Release: stevedore.Release{
									Name:         "nginx-a-1",
									Namespace:    "default",
									ChartVersion: "1.2.0",
								},
							},
						},
					},
				},
				nginxBResult,
			},
		}

		info := manifest.Info{ManifestFiles: results}

		responses := stevedore.Responses{
			{ReleaseName: "nginx-a-1", ChartVersion: "1.2.0"}, {ReleaseName: "nginx-b-1"}, {ReleaseName: "nginx-b-2"},
		}
		actual := info.FilterBy(responses)

		assert.Equal(t, expected, actual)
	})
}
