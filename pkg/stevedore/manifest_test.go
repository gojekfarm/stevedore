package stevedore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gojek/stevedore/client/provider"
	mockPlugin "github.com/gojek/stevedore/pkg/internal/mocks/plugin"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestNewManifest(t *testing.T) {
	t.Run("Should Validate and return Stevedore Manifest", func(t *testing.T) {
		configString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chartSpec:
      name: x-stevedore-dependencies
      dependencies:
        - name: redis
          alias: cache-redis
          version: 1.9.4
          repository: 'http://helm-charts/'
    values:
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

		expectedStevedoreManifest := stevedore.Manifest{
			Kind:    stevedore.KindStevedoreManifest,
			Version: stevedore.ManifestCurrentVersion,
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{
						Name:      "x-stevedore",
						Namespace: "ns",
						ChartSpec: stevedore.ChartSpec{
							Name: "x-stevedore-dependencies",
							Dependencies: stevedore.Dependencies{
								stevedore.Dependency{
									Name:       "redis",
									Alias:      "cache-redis",
									Version:    "1.9.4",
									Repository: "http://helm-charts/"},
							},
						},
					},
					Configs: stevedore.Configs{
						"store": []interface{}{
							map[interface{}]interface{}{
								"name": "x-stevedore", "tags": []interface{}{"server"},
							},
						},
					},
				},
			},
		}

		actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

		if assert.NoError(t, err) {
			assert.Equal(t, expectedStevedoreManifest, *actualStevedoreManifest)
		}
	})

	t.Run("should not enforce for configs", func(t *testing.T) {
		configString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chartSpec:
      name: x-stevedore-dependencies
      dependencies:
        - name: redis
          alias: cache-redis
          version: 1.9.4
          repository: 'http://helm-charts/'
    values:`

		expectedStevedoreManifest := stevedore.Manifest{
			Kind:    stevedore.KindStevedoreManifest,
			Version: stevedore.ManifestCurrentVersion,
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{
						Name:      "x-stevedore",
						Namespace: "ns",
						ChartSpec: stevedore.ChartSpec{
							Name: "x-stevedore-dependencies",
							Dependencies: stevedore.Dependencies{
								stevedore.Dependency{
									Name:       "redis",
									Alias:      "cache-redis",
									Version:    "1.9.4",
									Repository: "http://helm-charts/"},
							},
						},
					},
				},
			},
		}

		actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

		if assert.NoError(t, err) {
			assert.Equal(t, expectedStevedoreManifest, *actualStevedoreManifest)
		}
	})

	t.Run("Validate should fail for type mismatch errors", func(t *testing.T) {

		t.Run("When deployTo is not an condition", func(t *testing.T) {
			configString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - release-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chart: chart/x-stevedore-dependencies
    values:
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Nil(t, actualStevedoreManifest)
		})

		t.Run("When applications is not an array", func(t *testing.T) {
			configString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
  name: x-stevedore
  release:
    name: x-stevedore
    namespace: ns
    chart: chart/x-stevedore-dependencies
    values:
  store:
    name: x-stevedore
    tag: server`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Nil(t, actualStevedoreManifest)
		})

		t.Run("When applications release values is not hash", func(t *testing.T) {
			configString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chart: chart/x-stevedore-dependencies
    values: [ "blah"]
  store:
    name: x-stevedore
    tag: server`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Nil(t, actualStevedoreManifest)
		})
	})

	t.Run("Validate should fail for field not present errors", func(t *testing.T) {

		t.Run("When deployTo is not present", func(t *testing.T) {
			configString := `
spec:
- release:
    name: x-stevedore
    namespace: ns
    chart: chart/x-stevedore-dependencies
    values:
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Key: 'Manifest.DeployTo' Error:Field validation for 'DeployTo' failed on the 'required' tag")
			assert.Nil(t, actualStevedoreManifest)
		})

		t.Run("When deployTo is empty", func(t *testing.T) {
			configString := `
deployTo:
spec:
- release:
    name: x-stevedore
    namespace: ns
    chart: chart/x-stevedore-dependencies
    values:
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Key: 'Manifest.DeployTo' Error:Field validation for 'DeployTo' failed on the 'required' tag")
			assert.Nil(t, actualStevedoreManifest)
		})

		t.Run("When applications is empty", func(t *testing.T) {
			configString := `
deployTo:
  - contextName: components-staging
spec:`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Key: 'Manifest.Spec' Error:Field validation for 'Spec' failed on the 'required' tag")
			assert.Nil(t, actualStevedoreManifest)
		})

		t.Run("When applications components is not present", func(t *testing.T) {
			configString := `
deployTo:
  - contextName: components-staging
spec:
- configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Key: 'Manifest.Spec[0].Release.Name' Error:Field validation for 'Name' failed on the 'required' tag")
			assert.Contains(t, err.Error(), "Key: 'Manifest.Spec[0].Release.Namespace' Error:Field validation for 'Namespace' failed on the 'required' tag")
			assert.Contains(t, err.Error(), "Key: 'Manifest.Spec[0].Release.Chart' Error:Field validation for 'Chart' failed on the 'EitherChartOrChartSpec' tag")
			assert.Nil(t, actualStevedoreManifest)
		})

		t.Run("When release name is not present", func(t *testing.T) {
			configString := `
deployTo:
  - contextName: components-staging
spec:
- release:
    namespace: ns
    chart: chart/x-stevedore-dependencies
    values: {}`

			actualStevedoreManifest, err := stevedore.NewManifest(strings.NewReader(configString))

			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "Key: 'Manifest.Spec[0].Release.Name' Error:Field validation for 'Name' failed on the 'required' tag")
			}
			assert.Nil(t, actualStevedoreManifest)
		})
	})
}

func TestManifestEnrichWith(t *testing.T) {
	t.Run("should return enriched manifest with final merged values given the overrides", func(t *testing.T) {

		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{
						Name:      "x-stevedore",
						Namespace: "ns",
						Chart:     "chart/x-stevedore-dependencies",
						Values: stevedore.Values{
							"key1": "x-stevedore-baseValueForK1",
							"key2": "x-stevedore-baseValueForK2",
							"key3": "x-stevedore-baseValueForK3",
						},
					},
					Configs: stevedore.Configs{
						"store": []map[string]interface{}{
							{"name": "x-stevedore", "tags": []string{"server"}},
						},
					},
				},
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{
						Name:      "x-abc-service",
						Namespace: "ns",
						Chart:     "chart/x-abc-service-dependencies",
						Values: stevedore.Values{
							"key1": "x-abc-service-baseValueForK1",
							"key2": "x-abc-service-baseValueForK2",
							"key4": "x-abc-service-baseValueForK3",
						},
					},
					Configs: stevedore.Configs{
						"store": []map[string]interface{}{
							{"name": "x-stevedore", "tags": []string{"server"}},
						},
					},
				},
			},
		}

		overrides := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key2": "firstOverrideForK2"}},
				{Matches: stevedore.Conditions{"applicationName": "x-abc-service", "environmentType": "staging"}, Values: stevedore.Values{"key1": "secondOverrideForK1"}},
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore", "environment": "some-specific-staging-env"}, Values: stevedore.Values{"key3": "fourthOverrideForK3"}},
				{Matches: stevedore.Conditions{"environment": "some-specific-staging-env"}, Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
			},
		}

		expected := stevedore.Manifest{
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.NewRelease("x-stevedore", "ns", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
						"key1": "x-stevedore-baseValueForK1",
						"key2": "firstOverrideForK2",
						"key3": "fourthOverrideForK3",
					}, stevedore.Substitute{}, stevedore.Overrides{
						Spec: stevedore.OverrideSpecifications{
							{Matches: stevedore.Conditions{"environment": "some-specific-staging-env"}, Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
							{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key2": "firstOverrideForK2"}},
							{Matches: stevedore.Conditions{"applicationName": "x-stevedore", "environment": "some-specific-staging-env"}, Values: stevedore.Values{"key3": "fourthOverrideForK3"}},
						},
					}),
					Configs: stevedore.Configs{
						"store": []map[string]interface{}{
							{"name": "x-stevedore", "tags": []string{"server"}},
						},
					},
				},
				stevedore.ReleaseSpecification{
					Release: stevedore.NewRelease("x-abc-service", "ns", "chart/x-abc-service-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
						"key1": "secondOverrideForK1",
						"key2": "thirdOverrideForK2",
						"key4": "x-abc-service-baseValueForK3",
					}, stevedore.Substitute{}, stevedore.Overrides{
						Spec: stevedore.OverrideSpecifications{
							{Matches: stevedore.Conditions{"environment": "some-specific-staging-env"}, Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
							{Matches: stevedore.Conditions{"applicationName": "x-abc-service", "environmentType": "staging"}, Values: stevedore.Values{"key1": "secondOverrideForK1"}},
						},
					}),
					Configs: stevedore.Configs{
						"store": []map[string]interface{}{
							{"name": "x-stevedore", "tags": []string{"server"}},
						},
					},
				},
			},
		}

		actual := manifest.EnrichWith(stevedore.Context{EnvironmentType: "staging", Environment: "some-specific-staging-env"}, overrides)

		assert.Equal(t, expected, actual)
	})
}

func TestManifestReplace(t *testing.T) {
	t.Run("should replace values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		envs := stevedore.Substitute{"NAME": "x-service"}

		config := map[string]interface{}{"ENV": "staging", "OPTION": "repl"}
		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(config, nil)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{{
				Release: stevedore.Release{
					Name:      "x-stevedore",
					Namespace: "ns",
					Chart:     "chart/x-stevedore-dependencies",
					Values: stevedore.Values{
						"name":   "${NAME}",
						"env":    "${ENV}",
						"option": "${OPTION}",
					},
				},
				Configs: stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				},
			}},
		}
		expected := stevedore.Manifest{
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{stevedore.NewReleaseSpecification(
				stevedore.Release{
					Name:      "x-stevedore",
					Namespace: "ns",
					Chart:     "chart/x-stevedore-dependencies",
					Values: stevedore.Values{
						"name":   "x-service",
						"env":    "staging",
						"option": "repl",
					},
				},
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, map[string]interface{}{"ENV": "staging", "NAME": "x-service", "OPTION": "repl"}),
			},
		}

		actual, err := manifest.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.True(t, cmp.Equal(expected, actual, cmpopts.IgnoreUnexported(stevedore.ReleaseSpecification{}, stevedore.Release{})))
	})

	t.Run("should fail when unable to replace values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		envs := stevedore.Substitute{"FAME": "service"}

		config := map[string]interface{}{"ENV": "staging", "OPTION": "repl"}
		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(config, nil).MaxTimes(2)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{stevedore.Conditions{
				stevedore.ConditionContextName: "components-staging"},
			},
			Spec: stevedore.ReleaseSpecifications{
				{
					Release: stevedore.Release{
						Name:      "x-stevedore",
						Namespace: "ns",
						Chart:     "chart/x-stevedore-dependencies",
						Values: stevedore.Values{
							"name": "${X_NAME}",
							"game": "${X_GAME}",
						},
					},
					Configs: stevedore.Configs{
						"store": []map[string]interface{}{
							{"name": "x-stevedore", "tags": []string{"server"}},
							{"name": "store-ns"},
						},
					},
				},
				{
					Release: stevedore.Release{
						Name:      "y-stevedore",
						Namespace: "ns",
						Chart:     "chart/y-stevedore-dependencies",
						Values: stevedore.Values{
							"name": "${Y_NAME}",
							"game": "${Y_GAME}",
						},
					},
					Configs: stevedore.Configs{
						"store": []map[string]interface{}{
							{"name": "y-stevedore", "tags": []string{"worker"}},
						},
					},
				},
			},
		}

		actual, err := manifest.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 4 variable(s):\n\t1. ${X_GAME}\n\t2. ${X_NAME}\n\t3. ${Y_GAME}\n\t4. ${Y_NAME}", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Equal(t, manifest, actual)
	})
}

func TestManifestIsApplicableFor(t *testing.T) {

	t.Run("should return true", func(t *testing.T) {
		manifest := stevedore.Manifest{DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "staging"}}}

		isApplicable := manifest.IsApplicableFor(stevedore.Context{Name: "staging"})

		assert.True(t, isApplicable)
	})

	t.Run("should return false", func(t *testing.T) {
		manifest := stevedore.Manifest{DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "staging"}}}

		isApplicable := manifest.IsApplicableFor(stevedore.Context{Name: "uat"})

		assert.False(t, isApplicable)
	})
}

func TestManifestNamespaces(t *testing.T) {
	t.Run("should return namespaces from all manifests uniquely", func(t *testing.T) {
		manifests := stevedore.Manifests{
			stevedore.Manifest{
				DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "env"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "x-stevedore", Namespace: "default"},
					},
				},
			},
			stevedore.Manifest{
				DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "y-stevedore", Namespace: "ns"},
					},
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "z-stevedore", Namespace: "default"},
					},
				},
			},
		}

		actual := manifests.Namespaces()

		expected := []string{"default", "ns"}
		assert.ElementsMatch(t, expected, actual)
	})
}

func TestManifestHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		manifest := stevedore.Manifest{
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{Dependencies: stevedore.Dependencies{{Name: "example"}}}},
				},
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
				},
			},
		}

		ok := manifest.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		manifest := stevedore.Manifest{
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
				},
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
				},
			},
		}

		ok := manifest.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		manifest := stevedore.Manifest{
			Spec: stevedore.ReleaseSpecifications{
				stevedore.ReleaseSpecification{
					Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
				},
			},
		}

		ok := manifest.HasBuildStep()
		assert.False(t, ok)
	})
}

func TestManifestsFilter(t *testing.T) {
	t.Run("should return true and filtered manifests and ignoreComponents", func(t *testing.T) {
		manifests := stevedore.Manifests{
			stevedore.Manifest{
				DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "env"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "x-stevedore"},
					},
				},
			},
			stevedore.Manifest{
				DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "y-stevedore"},
					},
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "z-stevedore"},
					},
				},
			},
		}

		expectedManifests := stevedore.Manifests{
			stevedore.Manifest{
				DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "env"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "x-stevedore"},
					},
				},
			},
		}
		expectedIgnoredReleases := stevedore.IgnoredReleases{
			stevedore.IgnoredRelease{
				Name:   "y-stevedore",
				Reason: "Not applicable for the context 'env'",
			},
			stevedore.IgnoredRelease{
				Name:   "z-stevedore",
				Reason: "Not applicable for the context 'env'",
			},
		}

		filteredManifests, ignoreComponents, result := manifests.Filter(stevedore.Context{Name: "env"})

		assert.True(t, result)
		assert.Equal(t, expectedIgnoredReleases, ignoreComponents)
		assert.Equal(t, expectedManifests, filteredManifests)
	})

	t.Run("should return false and empty manifests and ignoreComponents", func(t *testing.T) {
		manifests := stevedore.Manifests{
			stevedore.Manifest{
				DeployTo: stevedore.Matchers{stevedore.Conditions{stevedore.ConditionContextName: "env"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{Name: "x-stevedore"},
					},
				},
			},
		}

		filteredManifests, ignoreComponents, result := manifests.Filter(stevedore.Context{Name: "env"})

		assert.False(t, result)
		assert.Equal(t, stevedore.IgnoredReleases{}, ignoreComponents)
		assert.Equal(t, manifests, filteredManifests)
	})
}

func TestManifestsHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		manifest := stevedore.Manifests{
			stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{Dependencies: stevedore.Dependencies{{Name: "example"}}}},
					},
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
					},
				},
			},
			stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
					},
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
					},
				},
			},
		}

		ok := manifest.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		manifest := stevedore.Manifests{
			stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
					},
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
					},
				},
			},
			stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
					},
				},
			},
		}

		ok := manifest.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		manifest := stevedore.Manifests{
			stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
					},
				},
			},
		}

		ok := manifest.HasBuildStep()
		assert.False(t, ok)
	})
}

func TestManifestFormat(t *testing.T) {
	t.Run("it should be able to formatAs as yaml", func(t *testing.T) {
		manifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chartSpec:
      name: x-stevedore-dependencies
      dependencies:
        - name: redis
          alias: cache-redis
          version: 1.9.4
          repository: 'http://helm-charts/'
        - name: db
          alias: postgres-cluster
          version: 2.0.1
          repository: 'http://helm-charts/'
    values:
      redis:
        metrics:
          enabled: false
        create: true
      db:
        postgres:
          application:
            username: user
        create: false
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

		manifest, err := stevedore.NewManifest(strings.NewReader(manifestString))
		if !assert.NoError(t, err) {
			return
		}

		expected := `kind: StevedoreManifest
version: "2"
deployTo:
- contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chartSpec:
      name: x-stevedore-dependencies
      dependencies:
      - name: redis
        alias: cache-redis
        version: 1.9.4
        repository: http://helm-charts/
      - name: db
        alias: postgres-cluster
        version: 2.0.1
        repository: http://helm-charts/
    values:
      db:
        create: false
        postgres:
          application:
            username: user
      redis:
        create: true
        metrics:
          enabled: false
  configs:
    store:
    - name: x-stevedore
      tags:
      - server
`
		actual := fmt.Sprintf("%y", manifest)

		assert.Equal(t, expected, actual)
	})

	t.Run("it should be able to formatAs as json", func(t *testing.T) {
		manifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chartSpec:
      name: x-stevedore-dependencies
      dependencies:
        - name: redis
          alias: cache-redis
          version: 1.9.4
          repository: 'http://helm-charts/'
        - name: db
          alias: postgres-cluster
          version: 2.0.1
          repository: 'http://helm-charts/'
    values:
      redis:
        metrics:
          enabled: false
        create: true
      db:
        postgres:
          application:
            username: user
        create: false
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

		manifest, err := stevedore.NewManifest(strings.NewReader(manifestString))
		if !assert.NoError(t, err) {
			return
		}

		actual, err := stevedore.NewManifest(strings.NewReader(fmt.Sprintf("%j", manifest)))
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, manifest, actual)
	})

	t.Run("it should be able to formatAs as json", func(t *testing.T) {
		manifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components-staging
spec:
- release:
    name: x-stevedore
    namespace: ns
    chartSpec:
      name: x-stevedore-dependencies
      dependencies:
        - name: redis
          alias: cache-redis
          version: 1.9.4
          repository: 'http://helm-charts/'
        - name: db
          alias: postgres-cluster
          version: 2.0.1
          repository: 'http://helm-charts/'
    values:
      redis:
        metrics:
          enabled: false
        create: true
      db:
        postgres:
          application:
            username: user
        create: false
  configs:
    store:
      - name: x-stevedore
        tags: ["server"]`

		manifest, err := stevedore.NewManifest(strings.NewReader(manifestString))
		if !assert.NoError(t, err) {
			return
		}

		actual, err := stevedore.NewManifest(strings.NewReader(fmt.Sprintf("%#j", manifest)))
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, manifest, actual)
	})
}
