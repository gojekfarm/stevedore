package stevedore

import (
	"context"
	"fmt"
	"os"

	"github.com/chartmuseum/helm-push/pkg/chartmuseum"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// ChartManager will help in Build/Archive/Upload Chart
type ChartManager interface {
	Build(ctx context.Context, chartPath string) error
	Load(ctx context.Context, chartPath string) (*chart.Chart, error)
	Archive(ctx context.Context, ch *chart.Chart, outDir string) (string, error)
	UploadChart(ctx context.Context, name, url string) error
}

// DefaultChartManager will help in Build/Archive/Upload Chart
type DefaultChartManager struct {
}

// Build rebuilds a local charts directory from a lockfile.
//
// If the lockfile is not present, this will run a Manager.Update()
func (cm DefaultChartManager) Build(ctx context.Context, chartPath string) error {
	var settings = environment.EnvSettings{}
	manager := downloader.Manager{
		Out:        os.Stderr,
		ChartPath:  chartPath,
		HelmHome:   helmHome(),
		Keyring:    keyring(),
		SkipUpdate: false,
		Getters:    getter.All(settings),
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("chart build operation aborted")
	default:
		return manager.Build()
	}
}

// Load loads from a directory.
//
// This loads charts only from directories.
func (cm DefaultChartManager) Load(ctx context.Context, chartPath string) (*chart.Chart, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("chart load operation aborted")
	default:
		return chartutil.LoadDir(chartPath)
	}
}

// Archive creates an archived chart to the given directory.
func (cm DefaultChartManager) Archive(ctx context.Context, ch *chart.Chart, outDir string) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("chart archive operation aborted")
	default:
		return chartutil.Save(ch, outDir)
	}
}

// UploadChart uploads a chart package to given url
func (cm DefaultChartManager) UploadChart(ctx context.Context, name, url string) error {
	client, err := chartmuseum.NewClient(chartmuseum.URL(url))
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("chart upload operation aborted")
	default:
		_, err = client.UploadChartPackage(name, true)
		return err
	}
}

func helmHome() helmpath.Home {
	return helmpath.Home(os.ExpandEnv("$HOME/.helm"))
}

func keyring() string {
	return os.ExpandEnv("$HOME/.gnupg/pubring.gpg")
}
