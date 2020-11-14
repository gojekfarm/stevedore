package stevedore_test

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/pkg/config"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"

	"github.com/gojek/stevedore/pkg/internal/mocks/plugin"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestReleaseSpecificationEnrichWith(t *testing.T) {
	t.Run("should return enriched releaseSpecification with final merged values given the overrides", func(t *testing.T) {
		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"key1": "baseValueForK1",
					"key2": "baseValueForK2",
					"key3": "baseValueForK3",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
				},
			},
		}

		overrides := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore", "environmentType": "staging"}, Values: stevedore.Values{"key2": "firstOverrideForK2"}},
				{Matches: stevedore.Conditions{"applicationName": "x-abc-service", "environmentType": "staging"}, Values: stevedore.Values{"key1": "secondOverrideForK1"}},
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
				{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key3": "fourthOverrideForK3"}},
			},
		}

		expected := stevedore.ReleaseSpecification{
			Release: stevedore.NewRelease("x-stevedore", "ns", "helm-repo/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
				"key1": "baseValueForK1",
				"key2": "firstOverrideForK2",
				"key3": "fourthOverrideForK3",
			}, stevedore.Substitute{}, stevedore.Overrides{
				Spec: stevedore.OverrideSpecifications{
					{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
					{Matches: stevedore.Conditions{"applicationName": "x-stevedore"}, Values: stevedore.Values{"key3": "fourthOverrideForK3"}},
					{Matches: stevedore.Conditions{"applicationName": "x-stevedore", "environmentType": "staging"}, Values: stevedore.Values{"key2": "firstOverrideForK2"}},
				},
			}),
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
				},
			},
		}

		actual := releaseSpecification.EnrichWith(stevedore.Context{EnvironmentType: "staging"}, overrides)

		assert.Equal(t, expected, actual)
	})
}

func TestReleaseSpecificationReplace(t *testing.T) {
	t.Run("should replace values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		envs := stevedore.Substitute{"NAME": "x-service", "OPTION1": "repl4", "OPTION3": "repl7"}
		config1 := map[string]interface{}{"ENV": "staging", "OPTION": "repl", "OPTION1": "repl1", "OPTION2": "repl3", "OPTION3": "repl5"}
		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(config1, nil)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name":    "${NAME}",
					"envs":    "${ENV}",
					"option1": "${OPTION1}",
					"option2": "${OPTION2}",
					"option3": "${OPTION3}",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
		}
		expected := stevedore.NewReleaseSpecification(
			stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name":    "x-service",
					"envs":    "staging",
					"option1": "repl4",
					"option2": "repl3",
					"option3": "repl7",
				},
			}, stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
			map[string]interface{}{"ENV": "staging", "NAME": "x-service", "OPTION": "repl", "OPTION1": "repl4", "OPTION2": "repl3", "OPTION3": "repl7"})

		actual, err := releaseSpecification.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		unexported := cmpopts.IgnoreUnexported(stevedore.ReleaseSpecification{}, stevedore.Release{})
		if !cmp.Equal(expected, actual, unexported) {
			assert.Fail(t, cmp.Diff(expected, actual, unexported), "expected to be equal")
		}
	})

	t.Run("should not fetch variables from store if there are no variables for replacement", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name": "x-service",
					"env":  "staging",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
		}
		expected := stevedore.NewReleaseSpecification(
			stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name": "x-service",
					"env":  "staging",
				},
			}, stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
			nil)

		actual, err := releaseSpecification.Replace(stevedore.Context{Environment: "staging"}, stevedore.Substitute{}, config.Providers{})

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

		envs := stevedore.Substitute{"FAME": "x-service"}

		configFromProvider := map[string]interface{}{"ENV": "staging", "OPTION": "repl"}

		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(configFromProvider, nil)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name": "${NAME}",
					"game": "${GAME}",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
		}

		actual, err := releaseSpecification.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 2 variable(s):\n\t1. ${GAME}\n\t2. ${NAME}", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Equal(t, releaseSpecification, actual)
	})

	t.Run("should fail when unable to fetch group secrets", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("unable to fetch group secret at this time"))

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name": "${NAME}",
					"game": "${GAME}",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
		}

		actual, err := releaseSpecification.Replace(stevedore.Context{Environment: "staging"}, stevedore.Substitute{}, configProviders)

		if assert.NotNil(t, err) {
			assert.Equal(t, "error in fetching from provider: unable to fetch group secret at this time", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Equal(t, releaseSpecification, actual)
	})

	t.Run("should fail when unable fetch releaseSpecification configurations", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("unable to fetch configuration at this time"))

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name": "${NAME}",
					"game": "${GAME}",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
				},
			},
		}

		actual, err := releaseSpecification.Replace(stevedore.Context{Environment: "staging"}, stevedore.Substitute{}, configProviders)

		if assert.NotNil(t, err) {
			assert.Equal(t, "error in fetching from provider: unable to fetch configuration at this time", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Equal(t, releaseSpecification, actual)
	})
}

func TestReleaseSpecificationSubstitutedVariables(t *testing.T) {
	t.Run("should return substituted variables", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		envs := stevedore.Substitute{"NAME": "x-service", "OPTION1": "repl4", "OPTION3": "repl7"}

		configFromProvider := map[string]interface{}{
			"ENV":     "staging",
			"OPTION1": "repl1",
			"USED":    "no",
			"OPTION":  "repl",
			"OPTION2": "repl3",
			"OPTION3": "repl5",
		}

		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(configFromProvider, nil)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "helm-repo/x-stevedore-dependencies",
				Values: stevedore.Values{
					"name":    "${NAME}",
					"env":     "${ENV}",
					"option1": "${OPTION1}",
					"option2": "${OPTION2}",
					"option3": "${OPTION3}",
				},
			},
			Configs: stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "ns"},
				},
			},
		}

		expected := stevedore.Substitute{"OPTION1": "repl4", "OPTION2": "repl3", "OPTION3": "repl7", "ENV": "staging", "NAME": "x-service"}

		actual, err := releaseSpecification.Replace(stevedore.Context{Environment: "staging"}, envs, configProviders)

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		variables := actual.SubstitutedVariables()
		assert.Equal(t, expected, variables)
	})
}

func TestReleaseSpecificationHasBuildStep(t *testing.T) {
	t.Run("should return true if chart name is not provided and chartSpec is provided and dependencies are not empty", func(t *testing.T) {
		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{Dependencies: stevedore.Dependencies{{Name: "example"}}}},
		}

		ok := releaseSpecification.HasBuildStep()

		assert.True(t, ok)
	})

	t.Run("should return false if chart name is not provided and chartSpec is provided and dependencies are empty", func(t *testing.T) {
		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}},
		}

		ok := releaseSpecification.HasBuildStep()

		assert.False(t, ok)
	})

	t.Run("should return false if chart name is provided and chartSpec is not provided", func(t *testing.T) {
		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{ChartSpec: stevedore.ChartSpec{}, Chart: "example-service"},
		}

		ok := releaseSpecification.HasBuildStep()

		assert.False(t, ok)
	})
}

func TestReleaseSpecificationContainsDependency(t *testing.T) {
	releaseSpecification := stevedore.ReleaseSpecification{
		Release: stevedore.Release{
			ChartSpec: stevedore.ChartSpec{
				Dependencies: stevedore.Dependencies{
					{Name: "dep-1", Alias: "alias-1"},
					{Name: "dep-1", Alias: "alias-2"},
					{Name: "dep-2"},
				},
			},
		},
	}
	t.Run("should return true if releaseSpecification contains given chart name as dependency", func(t *testing.T) {
		expected := stevedore.Dependencies{
			{Name: "dep-1", Alias: "alias-1"},
			{Name: "dep-1", Alias: "alias-2"},
		}

		actual, ok := releaseSpecification.ContainsDependency("dep-1")

		assert.True(t, ok)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return false if releaseSpecification does not contain given chart name as dependency", func(t *testing.T) {
		expected := stevedore.Dependencies{}

		actual, ok := releaseSpecification.ContainsDependency("dep-3")

		assert.False(t, ok)
		assert.Equal(t, expected, actual)
	})
}
