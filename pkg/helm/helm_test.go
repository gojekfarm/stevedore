package helm_test

import (
	"context"
	"fmt"
	"github.com/gojek/stevedore/pkg/internal/mocks"
	helmMock "github.com/gojek/stevedore/pkg/internal/mocks/helm"
	osMock "github.com/gojek/stevedore/pkg/internal/mocks/os"
	"github.com/golang/mock/gomock"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/provenance"
	"os"
	"testing"

	"github.com/gojek/stevedore/pkg/helm"
	"github.com/stretchr/testify/assert"
	helmPkg "k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
)

func TestLocateChartPath(t *testing.T) {
	t.Run("should return the chart path when existing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFileInfo := osMock.NewMockFileInfo(ctrl)
		mockChartDownloader := helmMock.NewMockChartDownloader(ctrl)
		mockChartVerifier := helmMock.NewMockChartVerifier(ctrl)

		mockFs.EXPECT().Stat("/tmp/example-chart").Return(mockFileInfo, nil)

		chartPath, err := helm.LocateChartPath(mockFs, "/tmp/example-chart", "0.0.1", false, "", mockChartDownloader, mockChartVerifier)

		assert.NoError(t, err)
		assert.Equal(t, "/tmp/example-chart", chartPath)
	})

	t.Run("should return error if path given is directory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFileInfo := osMock.NewMockFileInfo(ctrl)
		mockChartDownloader := helmMock.NewMockChartDownloader(ctrl)
		mockChartVerifier := helmMock.NewMockChartVerifier(ctrl)

		mockFs.EXPECT().Stat("/tmp/example-chart").Return(mockFileInfo, nil)
		mockFileInfo.EXPECT().IsDir().Return(true)

		_, err := helm.LocateChartPath(mockFs, "/tmp/example-chart", "0.0.1", true, "", mockChartDownloader, mockChartVerifier)

		assert.Error(t, err)
		assert.Equal(t, "cannot verify a directory", err.Error())
	})

	t.Run("should return error if it has dot prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockChartDownloader := helmMock.NewMockChartDownloader(ctrl)
		mockChartVerifier := helmMock.NewMockChartVerifier(ctrl)

		chartPath := "./tmp/example-chart"
		mockFs.EXPECT().Stat(chartPath).Return(nil, fmt.Errorf("unable to get file stats"))

		_, err := helm.LocateChartPath(mockFs, chartPath, "0.0.1", true, "", mockChartDownloader, mockChartVerifier)

		if assert.Error(t, err) {
			assert.Equal(t, "path \"./tmp/example-chart\" not found", err.Error())
		}
	})

	t.Run("should return error if verification of chart fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFileInfo := osMock.NewMockFileInfo(ctrl)
		mockChartDownloader := helmMock.NewMockChartDownloader(ctrl)
		mockChartVerifier := helmMock.NewMockChartVerifier(ctrl)

		chartName := "example-chart"
		chartPath := fmt.Sprintf("./tmp/%s", chartName)
		mockFs.EXPECT().Stat(chartPath).Return(mockFileInfo, nil)
		mockFileInfo.EXPECT().IsDir().Return(false)
		expectedError := fmt.Errorf("chart verification failed")
		mockChartVerifier.EXPECT().VerifyChart(
			gomock.Any(),
			gomock.Any()).DoAndReturn(func(path string, keyring string) (*provenance.Verification, error) {
			assert.Contains(t, path, chartName)
			return nil, expectedError
		})

		_, err := helm.LocateChartPath(mockFs, chartPath, "0.0.1", true, "", mockChartDownloader, mockChartVerifier)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should append helm path and return abs if chart is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFileInfo := osMock.NewMockFileInfo(ctrl)
		mockChartDownloader := helmMock.NewMockChartDownloader(ctrl)
		mockChartVerifier := helmMock.NewMockChartVerifier(ctrl)

		chartName := "example-chart"
		version := "0.0.1"
		mockFs.EXPECT().Stat(gomock.Any()).Return(mockFileInfo, fmt.Errorf("failed to get stats"))
		mockFs.EXPECT().Stat(gomock.Any()).Return(mockFileInfo, nil)

		actualChartPath, err := helm.LocateChartPath(mockFs, chartName, version, true, "", mockChartDownloader, mockChartVerifier)

		assert.NoError(t, err)
		assert.Contains(t, actualChartPath, chartName)
	})

	t.Run("should download chart and return chart path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFileInfo := osMock.NewMockFileInfo(ctrl)
		mockChartDownloader := helmMock.NewMockChartDownloader(ctrl)
		mockChartVerifier := helmMock.NewMockChartVerifier(ctrl)

		chartName := "example-chart"
		version := "0.0.1"
		mockFs.EXPECT().Stat(gomock.Any()).Return(mockFileInfo, fmt.Errorf("failed to get stats")).MaxTimes(2)

		mockChartDownloader.EXPECT().DownloadTo(chartName, version, helmpath.Home(helm.HomePath()).Archive()).Return(chartName, nil, nil)

		actualChartPath, err := helm.LocateChartPath(mockFs, chartName, version, true, "", mockChartDownloader, mockChartVerifier)

		assert.NoError(t, err)
		assert.Contains(t, actualChartPath, chartName)
	})
}

func TestUpstallWithFakeClient(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.SkipNow()
	}
	releaseName := "postgres-through-test"
	chartName := "stable/postgresql"
	chartVersion := ""
	var currentReleaseVersion int32
	namespace := "default"
	t.Run("should not install if dry run is enabled", func(t *testing.T) {
		helmClient := &helmPkg.FakeClient{}
		client := &helm.DefaultClient{
			Interface: helmClient,
		}
		valuesYAML := ``

		upstallResponse, err := client.Upstall(context.TODO(), releaseName, chartName, chartVersion, currentReleaseVersion, namespace, valuesYAML, true, 60)
		assert.NoError(t, err)
		assert.NotEqual(t, "", upstallResponse.ChartVersion)

		response, _ := helmClient.ListReleases()
		assert.Equal(t, int64(0), response.Count)
	})

	t.Run("should install if dry run is disabled", func(t *testing.T) {
		helmClient := &helmPkg.FakeClient{}
		client := &helm.DefaultClient{
			Interface: helmClient,
		}

		valuesYAML := `image:
  pullPolicy: IfNotPresent`

		_, err := client.Upstall(context.TODO(), releaseName, chartName, chartVersion, currentReleaseVersion, namespace, valuesYAML, false, 60)

		assert.NoError(t, err)
		listResponse, _ := helmClient.ListReleases()
		assert.Equal(t, int64(1), listResponse.Count)

		responseContent, _ := helmClient.ReleaseContent(releaseName)
		assert.Equal(t, int32(1), responseContent.GetRelease().GetVersion())
		assert.Equal(t, valuesYAML, responseContent.GetRelease().GetConfig().GetRaw())
	})

	t.Run("should upgrade with given values if release is already installed", func(t *testing.T) {
		valuesYAML := `image:
  pullPolicy: IfNotPresent`

		helmClient := &helmPkg.FakeClient{
			Rels: []*release.Release{{
				Name:      releaseName,
				Namespace: namespace,
				Version:   1,
				Config:    &chart.Config{Raw: valuesYAML},
			}},
		}
		client := &helm.DefaultClient{
			Interface: helmClient,
		}

		updatedValuesYAML := `image:
  pullPolicy: Always`

		_, err := client.Upstall(context.TODO(), releaseName, chartName, chartVersion, currentReleaseVersion, namespace, updatedValuesYAML, false, 60)

		assert.NoError(t, err)
		response, _ := helmClient.ReleaseContent(releaseName)

		assert.Equal(t, int32(2), response.GetRelease().GetVersion())
		assert.Equal(t, updatedValuesYAML, response.GetRelease().GetConfig().GetRaw())
	})
}
