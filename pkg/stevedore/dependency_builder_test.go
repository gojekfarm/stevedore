package stevedore_test

import (
	"bytes"
	"context"
	"github.com/chartmuseum/helm-push/pkg/helm"
	chartMocks "github.com/gojek/stevedore/pkg/internal/mocks/chart"
	httpMocks "github.com/gojek/stevedore/pkg/internal/mocks/http"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"k8s.io/helm/pkg/repo"
	"net/http"
	"testing"
)

func TestNoopDependencyBuilderBuild(t *testing.T) {
	t.Run("should be able to build the dependencies", func(t *testing.T) {
		dependencies := stevedore.Dependencies{{
			Name:       "postgres",
			Repository: "http://localhost",
			Alias:      "db",
			Version:    "0.0.1",
		}}
		manifestFiles := stevedore.ManifestFiles{{
			File: "example.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						ChartSpec: stevedore.ChartSpec{
							Name:         "example-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}

		builder := stevedore.NoopDependencyBuilder{}

		actual, err := builder.Build(context.TODO(), manifestFiles)

		assert.NoError(t, err)
		assert.Equal(t, manifestFiles, actual)
	})
}

func TestNoopDependencyBuilderBuildChart(t *testing.T) {
	t.Run("should be able to build the chart", func(t *testing.T) {
		dependencies := stevedore.Dependencies{{
			Name:       "postgres",
			Repository: "http://localhost",
			Alias:      "db",
			Version:    "0.0.1",
		}}

		releaseSpecification := stevedore.ReleaseSpecification{
			Release: stevedore.Release{
				ChartSpec: stevedore.ChartSpec{
					Name:         "example-dependencies",
					Dependencies: dependencies,
				},
			},
		}

		builder := stevedore.NoopDependencyBuilder{}

		actual, built, err := builder.BuildChart(context.TODO(), releaseSpecification)

		assert.NoError(t, err)
		assert.False(t, built)
		assert.Equal(t, releaseSpecification, actual)
	})
}

func TestNoopDependencyBuilderUpdateRepo(t *testing.T) {
	t.Run("should be able to update repo", func(t *testing.T) {
		builder := stevedore.NoopDependencyBuilder{}

		err := builder.UpdateRepo()

		assert.NoError(t, err)
	})
}

func TestDefaultDependencyBuilderBuild(t *testing.T) {
	t.Run("should be able to build the dependencies for a new chart with helm repo url having trailing slash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dependencies := stevedore.Dependencies{{
			Name:       "postgres",
			Repository: "http://localhost",
			Alias:      "db",
			Version:    "0.0.1",
		}}

		manifestFiles := stevedore.ManifestFiles{{
			File: "example.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						ChartSpec: stevedore.ChartSpec{
							Name:         "example-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}
		expectedManifests := stevedore.ManifestFiles{{
			File: "example.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						Chart:        "forTest/example-dependencies",
						ChartVersion: "0.0.1",
						ChartSpec: stevedore.ChartSpec{
							Name:         "example-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}
		helmRepo := helm.Repo{Entry: &repo.Entry{Name: "forTest", URL: "http://test.local/"}}

		mockHTTPClient := httpMocks.NewMockClient(ctrl)
		chartBuilder := chartMocks.NewMockChartBuilder(ctrl)

		mockHTTPClient.EXPECT().Get("http://test.local/api/charts/example-dependencies").Return(&http.Response{StatusCode: 404}, nil)
		chartBuilder.EXPECT().Build(context.TODO(), "example-dependencies", "0.0.1", "6a17c442", dependencies)

		dependencyBuilder := stevedore.NewDefaultDependencyBuilder(helmRepo, mockHTTPClient, chartBuilder, nil)

		files, err := dependencyBuilder.Build(context.TODO(), manifestFiles)

		assert.NoError(t, err)
		assert.Equal(t, expectedManifests, files)
	})

	t.Run("should be able to build the dependencies for a new chart with helm repo url not having trailing slash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dependencies := stevedore.Dependencies{{
			Name:       "postgres",
			Repository: "http://localhost",
			Alias:      "db",
			Version:    "0.0.1",
		}}

		manifestFiles := stevedore.ManifestFiles{{
			File: "example.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						ChartSpec: stevedore.ChartSpec{
							Name:         "example-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}
		expectedManifests := stevedore.ManifestFiles{{
			File: "example.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						Chart:        "forTest/example-dependencies",
						ChartVersion: "0.0.1",
						ChartSpec: stevedore.ChartSpec{
							Name:         "example-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}
		helmRepo := helm.Repo{Entry: &repo.Entry{Name: "forTest", URL: "http://test.local"}}

		mockHTTPClient := httpMocks.NewMockClient(ctrl)
		chartBuilder := chartMocks.NewMockChartBuilder(ctrl)

		mockHTTPClient.EXPECT().Get("http://test.local/api/charts/example-dependencies").Return(&http.Response{StatusCode: 404}, nil)
		chartBuilder.EXPECT().Build(context.TODO(), "example-dependencies", "0.0.1", "6a17c442", dependencies)

		dependencyBuilder := stevedore.NewDefaultDependencyBuilder(helmRepo, mockHTTPClient, chartBuilder, nil)

		files, err := dependencyBuilder.Build(context.TODO(), manifestFiles)

		assert.NoError(t, err)
		assert.Equal(t, expectedManifests, files)
	})

	t.Run("should be able to build the dependencies for existing chart", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dependencies := stevedore.Dependencies{{
			Name:       "postgres",
			Repository: "http://localhost",
			Alias:      "db",
			Version:    "0.0.1",
		}}

		manifestFiles := stevedore.ManifestFiles{{
			File: "example-a.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						ChartSpec: stevedore.ChartSpec{
							Name:         "a-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}, {
			File: "example-b.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						ChartSpec: stevedore.ChartSpec{
							Name:         "b-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}, {
			File: "example-c.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						ChartSpec: stevedore.ChartSpec{
							Name:         "c-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}
		expectedManifests := stevedore.ManifestFiles{{
			File: "example-a.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						Chart:        "forTest/a-dependencies",
						ChartVersion: "0.0.3",
						ChartSpec: stevedore.ChartSpec{
							Name:         "a-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}, {
			File: "example-b.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						Chart:        "forTest/b-dependencies",
						ChartVersion: "0.2.0",
						ChartSpec: stevedore.ChartSpec{
							Name:         "b-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}, {
			File: "example-c.yaml",
			Manifest: stevedore.Manifest{
				Spec: stevedore.ReleaseSpecifications{{
					Release: stevedore.Release{
						Chart:        "forTest/c-dependencies",
						ChartVersion: "1.0.0",
						ChartSpec: stevedore.ChartSpec{
							Name:         "c-dependencies",
							Dependencies: dependencies,
						},
					},
				}},
			},
		}}

		knownDependencyChartA := ioutil.NopCloser(bytes.NewReader([]byte(`[
{"name": "a-dependencies",
  "version": "0.0.2",
  "appVersion": "9.0.0"}]`)))
		knownDependencyChartB := ioutil.NopCloser(bytes.NewReader([]byte(`[
{"name": "b-dependencies",
  "version": "0.1.10",
  "appVersion": "9.0.0"}]`)))
		knownDependencyChartC := ioutil.NopCloser(bytes.NewReader([]byte(`[
{"name": "c-dependencies",
  "version": "0.10.10",
  "appVersion": "9.0.0"}]`)))
		helmRepo := helm.Repo{Entry: &repo.Entry{Name: "forTest", URL: "http://test.local/stable"}}

		mockHTTPClient := httpMocks.NewMockClient(ctrl)
		chartBuilder := chartMocks.NewMockChartBuilder(ctrl)

		mockHTTPClient.EXPECT().Get("http://test.local/stable/api/charts/a-dependencies").Return(&http.Response{Body: knownDependencyChartA}, nil)
		mockHTTPClient.EXPECT().Get("http://test.local/stable/api/charts/b-dependencies").Return(&http.Response{Body: knownDependencyChartB}, nil)
		mockHTTPClient.EXPECT().Get("http://test.local/stable/api/charts/c-dependencies").Return(&http.Response{Body: knownDependencyChartC}, nil)
		chartBuilder.EXPECT().Build(context.TODO(), "a-dependencies", "0.0.3", "6a17c442", dependencies)
		chartBuilder.EXPECT().Build(context.TODO(), "b-dependencies", "0.2.0", "6a17c442", dependencies)
		chartBuilder.EXPECT().Build(context.TODO(), "c-dependencies", "1.0.0", "6a17c442", dependencies)

		dependencyBuilder := stevedore.NewDefaultDependencyBuilder(helmRepo, mockHTTPClient, chartBuilder, nil)

		files, err := dependencyBuilder.Build(context.TODO(), manifestFiles)

		assert.NoError(t, err)
		assert.Equal(t, expectedManifests, files)
	})
}
