package stevedore

import (
	"context"
	"fmt"
	"github.com/chartmuseum/helm-push/pkg/helm"
	"github.com/gojek/stevedore/pkg/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"
	"os"
	"testing"
)

func TestNewChartBuilder(t *testing.T) {
	t.Run("should build chart without any error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		entry := repo.Entry{URL: "url"}
		helmRepo := helm.Repo{Entry: &entry}
		helmChart := chart.Chart{}
		chartName := "chartName"
		tempDirName := "temp"

		chartManager := mocks.NewMockChartManager(ctrl)
		fileUtils := mocks.NewMockFileUtils(ctrl)

		fileUtils.EXPECT().TempDir("postgresql").Return(tempDirName, nil)
		fileUtils.EXPECT().WriteFile("temp/requirements.yaml", gomock.Any(), gomock.Any()).DoAndReturn(
			func(filename string, data []byte, perm os.FileMode) error {
				expected := map[string]interface{}{"dependencies": []interface{}{map[interface{}]interface{}{"name": "postgres", "alias": "db", "repository": "https://localhost", "version": "1.0.2"}}}
				result := map[string]interface{}{}
				err := yaml.Unmarshal(data, &result)
				assert.NoError(t, err)
				if !cmp.Equal(expected, result) {
					assert.Fail(t, cmp.Diff(expected, result))
				}
				return nil
			},
		)
		fileUtils.EXPECT().WriteFile("temp/Chart.yaml", gomock.Any(), gomock.Any()).DoAndReturn(
			func(filename string, data []byte, perm os.FileMode) error {
				expected := map[string]interface{}{"name": "postgresql", "version": "0.0.1", "appVersion": "1.0.0"}
				result := map[string]interface{}{}
				err := yaml.Unmarshal(data, &result)
				assert.NoError(t, err)
				if !cmp.Equal(expected, result) {
					assert.Fail(t, cmp.Diff(expected, result))
				}
				return nil
			},
		)
		fileUtils.EXPECT().RemoveAll(tempDirName).Return(nil)

		chartManager.EXPECT().Build(context.TODO(), tempDirName).Return(nil)
		chartManager.EXPECT().Load(context.TODO(), tempDirName).Return(&helmChart, nil)
		chartManager.EXPECT().Archive(context.TODO(), &helmChart, tempDirName).Return(chartName, nil)
		chartManager.EXPECT().UploadChart(context.TODO(), chartName, "url").Return(nil)

		chartBuilder := NewChartBuilder(helmRepo, chartManager, fileUtils)
		dependencies := Dependencies{{Name: "postgres", Alias: "db", Repository: "https://localhost", Version: "1.0.2"}}

		err := chartBuilder.Build(context.TODO(), "postgresql", "0.0.1", "1.0.0", dependencies)

		assert.Nil(t, err)
	})

	t.Run("should handle error returned when creating temp dir", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		entry := repo.Entry{URL: "url"}
		helmRepo := helm.Repo{Entry: &entry}
		tempDirName := "temp"

		chartManager := mocks.NewMockChartManager(ctrl)
		fileUtils := mocks.NewMockFileUtils(ctrl)

		fileUtils.EXPECT().TempDir("postgresql").Return(tempDirName, fmt.Errorf("unable to create temp dir"))

		chartBuilder := NewChartBuilder(helmRepo, chartManager, fileUtils)
		dependencies := Dependencies{{Name: "postgres", Alias: "db", Repository: "https://localhost", Version: "1.0.2"}}

		err := chartBuilder.Build(context.TODO(), "postgresql", "0.0.1", "1.0.0", dependencies)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to create temp dir", err.Error())
		}
	})

	t.Run("should handle error returned when building chart", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		entry := repo.Entry{URL: "url"}
		helmRepo := helm.Repo{Entry: &entry}
		tempDirName := "temp"

		chartManager := mocks.NewMockChartManager(ctrl)
		fileUtils := mocks.NewMockFileUtils(ctrl)

		fileUtils.EXPECT().TempDir("postgresql").Return(tempDirName, nil)
		fileUtils.EXPECT().WriteFile("temp/requirements.yaml", gomock.Any(), gomock.Any()).DoAndReturn(
			func(filename string, data []byte, perm os.FileMode) error {
				expected := map[string]interface{}{"dependencies": []interface{}{map[interface{}]interface{}{"name": "postgres", "alias": "db", "repository": "https://localhost", "version": "1.0.2"}}}
				result := map[string]interface{}{}
				err := yaml.Unmarshal(data, &result)
				assert.NoError(t, err)
				if !cmp.Equal(expected, result) {
					assert.Fail(t, cmp.Diff(expected, result))
				}
				return nil
			},
		)
		fileUtils.EXPECT().WriteFile("temp/Chart.yaml", gomock.Any(), gomock.Any()).DoAndReturn(
			func(filename string, data []byte, perm os.FileMode) error {
				expected := map[string]interface{}{"name": "postgresql", "version": "0.0.1", "appVersion": "1.0.0"}
				result := map[string]interface{}{}
				err := yaml.Unmarshal(data, &result)
				assert.NoError(t, err)
				if !cmp.Equal(expected, result) {
					assert.Fail(t, cmp.Diff(expected, result))
				}
				return nil
			},
		)
		fileUtils.EXPECT().RemoveAll(tempDirName).Return(nil)

		chartManager.EXPECT().Build(context.TODO(), tempDirName).Return(fmt.Errorf("unable to build"))

		chartBuilder := NewChartBuilder(helmRepo, chartManager, fileUtils)
		dependencies := Dependencies{{Name: "postgres", Alias: "db", Repository: "https://localhost", Version: "1.0.2"}}

		err := chartBuilder.Build(context.TODO(), "postgresql", "0.0.1", "1.0.0", dependencies)

		if assert.Error(t, err) {
			assert.Equal(t, "unable to build", err.Error())
		}
	})
}
