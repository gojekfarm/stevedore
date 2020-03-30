package stevedore_test

import (
	"fmt"
	"github.com/gojek/stevedore/client/provider"
	mockPlugin "github.com/gojek/stevedore/pkg/internal/mocks/plugin"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"

	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestReleaseSpecificationsReplace(t *testing.T) {
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

		releaseSpecifications := stevedore.ReleaseSpecifications{{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "company/x-stevedore-dependencies",
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
		}}
		expected := stevedore.ReleaseSpecifications{stevedore.NewReleaseSpecification(
			stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "company/x-stevedore-dependencies",
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
			},
			map[string]interface{}{"ENV": "staging", "NAME": "x-service", "OPTION": "repl"}),
		}

		actual, err := releaseSpecifications.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		unexported := cmpopts.IgnoreUnexported(stevedore.ReleaseSpecification{}, stevedore.Release{})
		if !cmp.Equal(expected, actual, unexported) {
			assert.Fail(t, cmp.Diff(expected, actual, unexported), "expected to be equal")
		}
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

		releaseSpecifications := stevedore.ReleaseSpecifications{
			{
				Release: stevedore.Release{
					Name:      "x-stevedore",
					Namespace: "ns",
					Chart:     "company/x-stevedore-dependencies",
					Values: stevedore.Values{
						"name":   "${X_NAME}",
						"game":   "${X_GAME}",
						"option": "${OPTION}",
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
					Chart:     "company/y-stevedore-dependencies",
					Values: stevedore.Values{
						"name": "${Y_NAME}",
						"game": "${Y_GAME}",
					},
				},
				Configs: stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "y-stevedore", "tags": []string{"server"}},
					},
				},
			},
		}

		actual, err := releaseSpecifications.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 4 variable(s):\n\t1. ${X_GAME}\n\t2. ${X_NAME}\n\t3. ${Y_GAME}\n\t4. ${Y_NAME}", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Equal(t, releaseSpecifications, actual)
	})

	t.Run("should fail with valid error when unable to contact store", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		envs := stevedore.Substitute{}
		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("unable to fetch configuration at this time"))

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		releaseSpecifications := stevedore.ReleaseSpecifications{
			{
				Release: stevedore.Release{
					Name:      "x-stevedore",
					Namespace: "ns",
					Chart:     "company/x-stevedore-dependencies",
					Values: stevedore.Values{
						"name":   "${X_NAME}",
						"game":   "${X_GAME}",
						"option": "${OPTION}",
					},
				},
				Configs: stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
					},
				},
			},
			{
				Release: stevedore.Release{
					Name:      "y-stevedore",
					Namespace: "ns",
					Chart:     "company/y-stevedore-dependencies",
					Values: stevedore.Values{
						"name": "${Y_NAME}",
						"game": "${Y_GAME}",
					},
				},
				Configs: stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "y-stevedore", "tags": []string{"server"}},
					},
				},
			},
		}

		actual, err := releaseSpecifications.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		if assert.NotNil(t, err) {
			assert.Equal(t, "error in fetching from provider: unable to fetch configuration at this time", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Equal(t, releaseSpecifications, actual)
	})
}

func TestReleaseSpecificationsEnrichWith(t *testing.T) {
	t.Run("should return enriched releaseSpecifications with final merged values given the overrides", func(t *testing.T) {

		releaseSpecifications := stevedore.ReleaseSpecifications{
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{
					Name:      "x-stevedore",
					Namespace: "ns",
					Chart:     "company/x-stevedore-dependencies",
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
					Chart:     "company/x-abc-service-dependencies",
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
		}

		overrides := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key2": "firstOverrideForK2"}},
				{Matches: stevedore.Conditions{"applicationName": "x-abc-service", "environmentType": "staging"}, Values: stevedore.Values{"key1": "secondOverrideForK1"}},
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore", "environment": "some-specific-staging-env"}, Values: stevedore.Values{"key3": "fourthOverrideForK3"}},
				{Matches: stevedore.Conditions{"environment": "some-specific-staging-env"}, Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
			},
		}

		expected := stevedore.ReleaseSpecifications{
			stevedore.ReleaseSpecification{
				Release: stevedore.NewRelease(
					"x-stevedore", "ns", "company/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
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
				Release: stevedore.NewRelease("x-abc-service", "ns", "company/x-abc-service-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
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
		}

		actual := releaseSpecifications.EnrichWith(stevedore.Context{EnvironmentType: "staging", Environment: "some-specific-staging-env"}, overrides)

		assert.Equal(t, expected, actual)
	})
}

func TestReleaseSpecificationsNamespaces(t *testing.T) {
	t.Run("should return namespaces from all releaseSpecifications uniquely", func(t *testing.T) {
		releaseSpecifications := stevedore.ReleaseSpecifications{
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{Name: "y-stevedore", Namespace: "ns"},
			},
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{Name: "z-stevedore", Namespace: "default"},
			},
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{Name: "z-stevedore", Namespace: "default"},
			},
		}

		actual := releaseSpecifications.Namespaces()

		expected := []string{"ns", "default"}
		assert.ElementsMatch(t, expected, actual)
	})
}

func TestReleaseSpecificationsHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		releaseSpecifications := stevedore.ReleaseSpecifications{
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{Dependencies: stevedore.Dependencies{{Name: "example"}}}},
			},
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
			},
		}

		ok := releaseSpecifications.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		releaseSpecifications := stevedore.ReleaseSpecifications{
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
			},
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
			},
		}

		ok := releaseSpecifications.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		releaseSpecifications := stevedore.ReleaseSpecifications{
			stevedore.ReleaseSpecification{
				Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
			},
		}

		ok := releaseSpecifications.HasBuildStep()

		assert.False(t, ok)
	})
}
