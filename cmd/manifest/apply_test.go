package manifest_test

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/cmd/internal/mocks/mockManifest"
	"github.com/gojek/stevedore/pkg/config"
	pkgManifest "github.com/gojek/stevedore/pkg/manifest"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"

	"github.com/gojek/stevedore/cmd/internal/mocks"
	"github.com/gojek/stevedore/cmd/internal/mocks/mockPlugin"
	"github.com/gojek/stevedore/cmd/internal/mocks/mockProvider"
	"github.com/gojek/stevedore/cmd/manifest"
	"github.com/gojek/stevedore/pkg/file"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestNewManifest(t *testing.T) {
	t.Run("should respect ignore, override order and return enriched info", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		context := map[string]string{
			provider.ManifestPathFile:      "/mock/services/service-one.yaml",
			provider.ManifestRecursiveFlag: "true",
		}
		mockManifestImpl := pkgManifest.ProviderImpl{
			Provider: mockManifestProvider,
			Context:  context,
		}

		stevedoreContext := stevedore.Context{
			Name:              "components-staging",
			KubernetesContext: "components",
			Labels: stevedore.Conditions{
				"environment":     "env",
				"environmentType": "staging",
				"contextType":     "components",
			},
		}
		labels := stevedore.Labels{
			{Name: "environmentType"},
			{Name: "environment"},
			{Name: "contextType"},
			{Name: "contextName"},
			{Name: "applicationName"},
		}

		manifestFileName := "/mock/configs/x-stevedore.yaml"
		manifestFiles := stevedore.ManifestFiles{
			stevedore.ManifestFile{File: manifestFileName, Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "x-stevedore",
							Namespace: "ns",
							Chart:     "company/x-stevedore-dependencies",
							Values: stevedore.Values{
								"COMPONENT": "${COMPONENT}",
								"ENV":       "${ENV}",
							},
						},
						Configs: stevedore.Configs{
							"store": []map[string]interface{}{
								{"name": "x-stevedore", "tags": []string{"server"}},
								{"name": "repo-ns"},
							},
						},
					},
				},
			}},
			stevedore.ManifestFile{File: "/mock/configs/y-stevedore.yaml", Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "y-stevedore",
							Namespace: "ns",
							Chart:     "company/y-stevedore-dependencies",
							Values:    stevedore.Values{},
						},
					},
				},
			}},
			stevedore.ManifestFile{File: "/mock/configs/z-stevedore.yaml", Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "env"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "z-stevedore",
							Namespace: "ns",
							Chart:     "company/y-stevedore-dependencies",
							Values:    stevedore.Values{},
						},
					},
				},
			}},
		}

		ignores := stevedore.Ignores{
			stevedore.Ignore{
				Matches: stevedore.Conditions{
					"contextName": "components-staging",
				},
				Releases: stevedore.IgnoredReleases{
					stevedore.IgnoredRelease{Name: "y-stevedore"},
				},
			},
		}

		overrides := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"environmentType": "staging",
					},
					Values: stevedore.Values{
						"ENVIRONMENT_TYPE": "${ENVIRONMENT_TYPE}",
						"ENVIRONMENT":      "environment",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"environment": "env",
					},
					Values: stevedore.Values{
						"ENVIRONMENT":  "${ENVIRONMENT}",
						"CONTEXT_TYPE": "context_type",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"contextType": "components",
					},
					Values: stevedore.Values{
						"CONTEXT_TYPE": "${CONTEXT_TYPE}",
						"CONTEXT_NAME": "context_name",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"contextName": "components-staging",
					},
					Values: stevedore.Values{
						"CONTEXT_NAME":     "${CONTEXT_NAME}",
						"APPLICATION_NAME": "application_name",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"applicationName": "x-stevedore",
					},
					Values: stevedore.Values{
						"APPLICATION_NAME": "${APPLICATION_NAME}",
					},
				},
			},
		}

		storeValues := map[string]interface{}{
			"COMPONENT":        "x-component",
			"ENVIRONMENT_TYPE": "test",
			"ENVIRONMENT":      "test-environment",
			"CONTEXT_TYPE":     "components",
			"CONTEXT_NAME":     "test-components",
			"APPLICATION_NAME": "x-application",
			"ENV":              "using group",
		}

		contextAsMap, _ := stevedoreContext.Map()
		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProviderOptions := []map[string]interface{}{
			{"name": "x-stevedore", "tags": []string{"server"}},
			{"name": "repo-ns"},
		}

		mockConfigProvider.EXPECT().Fetch(contextAsMap, mockConfigProviderOptions).Return(storeValues, nil)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()
		envValues := map[string]interface{}{"ENV": "using env"}

		mockContextProvider.EXPECT().Context().Return(stevedoreContext, nil)
		mockContextProvider.EXPECT().Labels().Return(labels, nil)
		mockIgnoreProvider.EXPECT().Ignores().Return(ignores, nil)
		mockOverrideProvider.EXPECT().Overrides().Return(overrides, nil)
		mockManifestProvider.EXPECT().Manifests(context).Return(manifestFiles, nil)
		mockEnvProvider.EXPECT().Envs().Return(provider.EnvsFiles{}, nil)

		mockEnvironment.EXPECT().Fetch().Return(envValues)
		mockReporter.EXPECT().ReportContext(stevedoreContext)

		mockReporter.EXPECT().ReportIgnores(ignores)
		mockReporter.EXPECT().ReportOverrides(overrides)
		mockReporter.EXPECT().ReportEnvs(gomock.Any())
		mockReporter.EXPECT().ReportSkipped(gomock.Any())
		mockReporter.EXPECT().ReportManifest(gomock.Any())

		expected := manifest.Info{
			Context: stevedore.Context{
				Name:              "components-staging",
				KubernetesContext: "components",
				Labels: stevedore.Conditions{
					"contextType":     "components",
					"environment":     "env",
					"environmentType": "staging",
				},
			},
			Ignored: stevedore.IgnoredReleases{
				{Name: "z-stevedore", Reason: "Not applicable for the context 'components-staging'"},
				{Name: "y-stevedore", Reason: ""},
			},
			ManifestFiles: stevedore.ManifestFiles{
				stevedore.ManifestFile{
					File: manifestFileName,
					Manifest: stevedore.Manifest{
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.NewReleaseSpecification(
								stevedore.NewRelease("x-stevedore", "ns", "company/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
									"COMPONENT":        "x-component",
									"ENVIRONMENT_TYPE": "test",
									"ENVIRONMENT":      "test-environment",
									"CONTEXT_TYPE":     "components",
									"CONTEXT_NAME":     "test-components",
									"APPLICATION_NAME": "x-application",
									"ENV":              "using env",
								}, stevedore.Substitute{}, stevedore.Overrides{}), stevedore.Configs{
									"store": []map[string]interface{}{
										{"name": "x-stevedore", "tags": []string{"server"}},
										{"name": "repo-ns"},
									},
								}, nil),
						},
					},
				},
			},
		}

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			configProviders,
		)

		assert.Nil(t, err)
		if assert.NotNil(t, manifests) {
			ignoreTypes := cmpopts.IgnoreTypes(stevedore.Substitute{}, stevedore.Overrides{})
			info := *manifests

			isEqual := cmp.Equal(expected, info, ignoreTypes)
			if !isEqual {
				assert.Fail(t, cmp.Diff(expected, info, ignoreTypes))
			}
		}
	})

	t.Run("should not fetch values from store and env if no placeholders are available in values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		context := map[string]string{
			provider.ManifestPathFile:      "/mock/services/service-one.yaml",
			provider.ManifestRecursiveFlag: "true",
		}
		mockManifestImpl := pkgManifest.ProviderImpl{
			Provider: mockManifestProvider,
			Context:  context,
		}

		stevedoreContext := stevedore.Context{
			Name:              "components-staging",
			KubernetesContext: "components",
			Labels: stevedore.Conditions{
				"environment":     "env",
				"environmentType": "staging",
				"contextType":     "components",
			},
		}
		labels := stevedore.Labels{
			{Name: "environmentType"},
			{Name: "environment"},
			{Name: "contextType"},
			{Name: "contextName"},
			{Name: "applicationName"},
		}

		manifestFileName := "/mock/configs/x-stevedore.yaml"
		manifestFiles := stevedore.ManifestFiles{
			stevedore.ManifestFile{File: manifestFileName, Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "x-stevedore",
							Namespace: "ns",
							Chart:     "company/x-stevedore-dependencies",
							Values: stevedore.Values{
								"COMPONENT": "x-component",
							},
						},
						Configs: stevedore.Configs{
							"store": []map[string]interface{}{
								{"name": "x-stevedore", "tags": []string{"server"}},
								{"name": "repo-ns"},
							},
						},
					},
				},
			}},
			stevedore.ManifestFile{File: "/mock/configs/y-stevedore.yaml", Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "y-stevedore",
							Namespace: "ns",
							Chart:     "company/y-stevedore-dependencies",
							Values:    stevedore.Values{},
						},
					},
				},
			}},
		}

		ignores := stevedore.Ignores{
			stevedore.Ignore{
				Matches: stevedore.Conditions{
					"contextName": "components-staging",
				},
				Releases: stevedore.IgnoredReleases{
					stevedore.IgnoredRelease{Name: "y-stevedore"},
				},
			},
		}

		overrides := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"environmentType": "staging",
					},
					Values: stevedore.Values{
						"ENVIRONMENT_TYPE": "test",
						"ENVIRONMENT":      "environment",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"environment": "env",
					},
					Values: stevedore.Values{
						"ENVIRONMENT":  "test-environment",
						"CONTEXT_TYPE": "context_type",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"contextType": "components",
					},
					Values: stevedore.Values{
						"CONTEXT_TYPE": "components",
						"CONTEXT_NAME": "context_name",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"contextName": "components-staging",
					},
					Values: stevedore.Values{
						"CONTEXT_NAME":     "test-components",
						"APPLICATION_NAME": "application_name",
					},
				},
				stevedore.OverrideSpecification{
					Matches: stevedore.Conditions{
						"applicationName": "x-stevedore",
					},
					Values: stevedore.Values{
						"APPLICATION_NAME": "x-application",
					},
				},
			},
		}

		mockContextProvider.EXPECT().Context().Return(stevedoreContext, nil)
		mockContextProvider.EXPECT().Labels().Return(labels, nil)
		mockIgnoreProvider.EXPECT().Ignores().Return(ignores, nil)
		mockOverrideProvider.EXPECT().Overrides().Return(overrides, nil)
		mockManifestProvider.EXPECT().Manifests(context).Return(manifestFiles, nil)
		mockEnvProvider.EXPECT().Envs().Return(provider.EnvsFiles{}, nil)
		mockEnvironment.EXPECT().Fetch().Return(stevedore.Substitute{})
		mockReporter.EXPECT().ReportContext(stevedoreContext)
		mockReporter.EXPECT().ReportIgnores(ignores)
		mockReporter.EXPECT().ReportOverrides(overrides)
		mockReporter.EXPECT().ReportEnvs(gomock.Any())
		mockReporter.EXPECT().ReportSkipped(gomock.Any())
		mockReporter.EXPECT().ReportManifest(gomock.Any())

		expected := manifest.Info{
			Context: stevedore.Context{
				Name:              "components-staging",
				KubernetesContext: "components",
				Labels: stevedore.Conditions{
					"contextType":     "components",
					"environment":     "env",
					"environmentType": "staging",
				},
			},
			Ignored: stevedore.IgnoredReleases{{Name: "y-stevedore"}},
			ManifestFiles: stevedore.ManifestFiles{
				stevedore.ManifestFile{
					File: manifestFileName,
					Manifest: stevedore.Manifest{
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.NewReleaseSpecification(
								stevedore.NewRelease("x-stevedore", "ns", "company/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{
									"COMPONENT":        "x-component",
									"ENVIRONMENT_TYPE": "test",
									"ENVIRONMENT":      "test-environment",
									"CONTEXT_TYPE":     "components",
									"CONTEXT_NAME":     "test-components",
									"APPLICATION_NAME": "x-application",
								}, stevedore.Substitute{}, stevedore.Overrides{}),
								stevedore.Configs{
									"store": []map[string]interface{}{
										{"name": "x-stevedore", "tags": []string{"server"}},
										{"name": "repo-ns"},
									},
								}, nil),
						},
					},
				},
			},
		}

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			config.Providers{},
		)

		assert.Nil(t, err)
		if assert.NotNil(t, manifests) {
			ignoreTypes := cmpopts.IgnoreTypes(stevedore.Substitute{}, stevedore.Overrides{})
			info := *manifests

			isEqual := cmp.Equal(expected, info, ignoreTypes)
			if !isEqual {
				assert.Fail(t, cmp.Diff(expected, info, ignoreTypes))
			}
		}
	})

	t.Run("should return error if values are not replaced", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		mockManifestImpl := pkgManifest.ProviderImpl{Provider: mockManifestProvider, Context: map[string]string{}}

		stevedoreContext := stevedore.Context{
			Name:              "components-staging",
			KubernetesContext: "components",
			Labels: stevedore.Conditions{
				"environment":     "env",
				"environmentType": "staging",
				"contextType":     "components",
			},
		}
		labels := stevedore.Labels{
			{Name: "environmentType"},
			{Name: "environment"},
			{Name: "contextType"},
			{Name: "contextName"},
			{Name: "applicationName"},
		}

		manifestFileName := "/mock/configs/x-stevedore.yaml"
		manifestFiles := stevedore.ManifestFiles{
			stevedore.ManifestFile{File: manifestFileName, Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "x-stevedore",
							Namespace: "ns",
							Chart:     "company/x-stevedore-dependencies",
							Values: stevedore.Values{
								"COMPONENT": "${COMPONENT}",
							},
						},
						Configs: stevedore.Configs{
							"store": []map[string]interface{}{
								{"name": "x-stevedore", "tags": []string{"server"}},
								{"name": "repo-ns"},
							},
						},
					},
				},
			}},
			stevedore.ManifestFile{File: "/mock/configs/y-stevedore.yaml", Manifest: stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "components-staging"}},
				Spec: stevedore.ReleaseSpecifications{
					stevedore.ReleaseSpecification{
						Release: stevedore.Release{
							Name:      "y-stevedore",
							Namespace: "ns",
							Chart:     "company/y-stevedore-dependencies",
							Values:    stevedore.Values{},
						},
					},
				},
			}},
		}

		ignores := stevedore.Ignores{
			stevedore.Ignore{
				Matches: stevedore.Conditions{
					"contextName": "components-staging",
				},
				Releases: stevedore.IgnoredReleases{
					stevedore.IgnoredRelease{Name: "y-stevedore"},
				},
			},
		}

		overrides := stevedore.Overrides{}
		storeValues := map[string]interface{}{}

		mockContextProvider.EXPECT().Context().Return(stevedoreContext, nil)
		mockContextProvider.EXPECT().Labels().Return(labels, nil)
		mockIgnoreProvider.EXPECT().Ignores().Return(ignores, nil)
		mockOverrideProvider.EXPECT().Overrides().Return(overrides, nil)
		mockManifestProvider.EXPECT().Manifests(gomock.Any()).Return(manifestFiles, nil)
		mockEnvProvider.EXPECT().Envs().Return(provider.EnvsFiles{}, nil)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})
		mockReporter.EXPECT().ReportContext(stevedoreContext)
		mockReporter.EXPECT().ReportIgnores(ignores)
		mockReporter.EXPECT().ReportOverrides(overrides)
		mockReporter.EXPECT().ReportEnvs(gomock.Any())

		mockConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		mockConfigProviderOptions := []map[string]interface{}{
			{"name": "x-stevedore", "tags": []string{"server"}},
			{"name": "repo-ns"},
		}
		contextAsMap, _ := stevedoreContext.Map()
		mockConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		mockConfigProvider.EXPECT().Fetch(contextAsMap, mockConfigProviderOptions).Return(storeValues, nil)

		plugins := provider.Plugins{"store": provider.ClientPlugin{PluginImpl: mockConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			configProviders,
		)

		if assert.NotNil(t, err) {
			assert.Equal(t, file.Errors{file.Error{Reason: stevedore.SubstituteError{"${COMPONENT}"}, Filename: "/mock/configs/x-stevedore.yaml"}}, err)
		}

		assert.Nil(t, manifests)
	})

	t.Run("should return error if there are error when fetching context ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		mockManifestImpl := pkgManifest.ProviderImpl{Provider: mockManifestProvider, Context: map[string]string{}}

		mockContextProvider.EXPECT().Context().Return(stevedore.Context{}, fmt.Errorf("unable to get context at this time"))

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			config.Providers{},
		)

		if assert.NotNil(t, err) {
			assert.Equal(t, "unable to get context at this time", err.Error())
		}

		assert.Nil(t, manifests)
	})

	t.Run("should return error if there are error when fetching ignores ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		mockManifestImpl := pkgManifest.ProviderImpl{Provider: mockManifestProvider, Context: map[string]string{}}

		mockContextProvider.EXPECT().Context().Return(stevedore.Context{}, nil)
		mockContextProvider.EXPECT().Labels().Return(stevedore.Labels{}, nil)
		mockEnvProvider.EXPECT().Envs().Return(provider.EnvsFiles{}, nil)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})
		mockIgnoreProvider.EXPECT().Ignores().Return(nil, fmt.Errorf("unable to get ignores at this time"))
		mockReporter.EXPECT().ReportContext(gomock.Any())
		mockReporter.EXPECT().ReportEnvs(gomock.Any())

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			config.Providers{},
		)

		if assert.NotNil(t, err) {
			assert.Equal(t, "unable to get ignores at this time", err.Error())
		}

		assert.Nil(t, manifests)
	})

	t.Run("should return error if there are error when fetching overrides ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		mockManifestImpl := pkgManifest.ProviderImpl{Provider: mockManifestProvider, Context: map[string]string{}}

		mockContextProvider.EXPECT().Context().Return(stevedore.Context{}, nil)
		mockContextProvider.EXPECT().Labels().Return(stevedore.Labels{}, nil)
		mockIgnoreProvider.EXPECT().Ignores().Return(stevedore.Ignores{}, nil)
		mockEnvProvider.EXPECT().Envs().Return(provider.EnvsFiles{}, nil)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})
		mockOverrideProvider.EXPECT().Overrides().Return(stevedore.Overrides{}, fmt.Errorf("unable to get overrides at this time"))
		mockReporter.EXPECT().ReportEnvs(gomock.Any())
		mockReporter.EXPECT().ReportContext(gomock.Any())
		mockReporter.EXPECT().ReportIgnores(gomock.Any())

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			config.Providers{},
		)

		if assert.NotNil(t, err) {
			assert.Equal(t, "unable to get overrides at this time", err.Error())
		}

		assert.Nil(t, manifests)
	})

	t.Run("should return error if there are error when fetching manifests ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockContextProvider := mockProvider.NewMockContextProvider(ctrl)
		mockIgnoreProvider := mockProvider.NewMockIgnoreProvider(ctrl)
		mockManifestProvider := mockManifest.NewMockProvider(ctrl)
		mockOverrideProvider := mockProvider.NewMockOverrideProvider(ctrl)
		mockEnvProvider := mockProvider.NewMockEnvProvider(ctrl)
		mockReporter := mockManifest.NewMockReporter(ctrl)
		mockManifestImpl := pkgManifest.ProviderImpl{Provider: mockManifestProvider, Context: map[string]string{}}

		mockContextProvider.EXPECT().Context().Return(stevedore.Context{}, nil)
		mockContextProvider.EXPECT().Labels().Return(stevedore.Labels{}, nil)
		mockIgnoreProvider.EXPECT().Ignores().Return(stevedore.Ignores{}, nil)
		mockOverrideProvider.EXPECT().Overrides().Return(stevedore.Overrides{}, nil)
		mockEnvProvider.EXPECT().Envs().Return(provider.EnvsFiles{}, nil)
		mockManifestProvider.EXPECT().Manifests(gomock.Any()).Return(nil, fmt.Errorf("unable to get manifests at this time"))
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})
		mockReporter.EXPECT().ReportContext(gomock.Any())
		mockReporter.EXPECT().ReportIgnores(gomock.Any())
		mockReporter.EXPECT().ReportEnvs(gomock.Any())
		mockReporter.EXPECT().ReportOverrides(gomock.Any())

		manifests, err := manifest.NewManifests(
			mockEnvironment,
			mockContextProvider,
			mockManifestImpl,
			mockOverrideProvider,
			mockIgnoreProvider,
			mockEnvProvider,
			mockReporter,
			config.Providers{},
		)

		if assert.NotNil(t, err) {
			assert.Equal(t, "unable to get manifests at this time", err.Error())
		}

		assert.Nil(t, manifests)
	})
}
