package stevedore

import (
	"bytes"
	"context"
	"log"
	"path/filepath"

	"github.com/chartmuseum/helm-push/pkg/helm"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
)

// ChartBuilder provides necessary methods
type ChartBuilder interface {
	Build(ctx context.Context, chartName, version, appVersion string, dependencies Dependencies) error
}

// DefaultChartBuilder will Build a chart and upload to given repo
type DefaultChartBuilder struct {
	repo      helm.Repo
	manager   ChartManager
	fileUtils FileUtils
}

// NewChartBuilder create DefaultChartBuilder
func NewChartBuilder(repo helm.Repo, manager ChartManager, fileUtils FileUtils) DefaultChartBuilder {
	return DefaultChartBuilder{repo: repo, manager: manager, fileUtils: fileUtils}
}

// Build build and uploads the chart and reports error if any
func (cb DefaultChartBuilder) Build(ctx context.Context, chartName, version, appVersion string, dependencies Dependencies) error {
	tempDir, err := cb.fileUtils.TempDir(chartName)
	if err != nil {
		return err
	}
	defer func() {
		err = cb.fileUtils.RemoveAll(tempDir)
		if err != nil {
			log.Printf("error while deleting %s: %v", tempDir, err)
		}
	}()

	ch, err := cb.createChart(ctx, tempDir, chartName, version, appVersion, dependencies)
	if err != nil {
		return err
	}
	name, err := cb.manager.Archive(ctx, ch, tempDir)
	if err != nil {
		return err
	}
	err = cb.manager.UploadChart(ctx, name, cb.repo.URL)
	return err
}

func (cb DefaultChartBuilder) createChart(ctx context.Context, tempDir, chartName, version, appVersion string, dependencies Dependencies) (*chart.Chart, error) {
	err := cb.writeChartFiles(tempDir, chartName, version, appVersion, dependencies)
	if err != nil {
		return nil, err
	}
	err = cb.manager.Build(ctx, tempDir)
	if err != nil {
		return nil, err
	}
	return cb.manager.Load(ctx, tempDir)
}

func (cb DefaultChartBuilder) writeChartFiles(tempDir, chartName, version, appVersion string, dependencies Dependencies) error {
	return cb.writeChartYaml(tempDir, chartName, version, appVersion, dependencies)
}

func (cb DefaultChartBuilder) writeChartYaml(tempDir, chartName, withVersion, appVersion string, dependencies Dependencies) error {
	chartFile := filepath.Join(tempDir, "Chart.yaml")
	buffer := bytes.Buffer{}
	chartContent := map[string]interface{}{"name": chartName, "version": withVersion, "appVersion": appVersion, "dependencies": dependencies}
	err := yaml.NewEncoder(&buffer).Encode(chartContent)

	if err != nil {
		return err
	}
	if err := cb.fileUtils.WriteFile(chartFile, buffer.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
	return nil
}
