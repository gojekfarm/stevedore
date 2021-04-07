package helm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/databus23/helm-diff/diff"
	"github.com/databus23/helm-diff/manifest"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Client is an abstraction through which helm can be interacted with
type Client interface {
	Upstall(ctx context.Context, releaseName, chartName, chartVersion string, plannedReleaseVersion int32, namespace, values string, dryRun bool, timeout int64, atomic bool) (UpstallResponse, error)
}

// DefaultClient is an implementation of helm.Client
type DefaultClient struct {
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		_ = log.Output(2, fmt.Sprintf(format, v...))
	}
}

var settings = cli.New()

// Upstall can install a new release or upgrade if already present
func (c *DefaultClient) Upstall(ctx context.Context, releaseName, chartName, chartVersion string, plannedReleaseVersion int32, namespace, values string, dryRun bool, timeout int64, atomic bool) (UpstallResponse, error) {
	cfg := &action.Configuration{}
	helmDriver := os.Getenv("HELM_DRIVER")
	configFlags := genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		Context:     &settings.KubeContext,
		BearerToken: &settings.KubeToken,
		APIServer:   &settings.KubeAPIServer,
		KubeConfig:  &settings.KubeConfig,
	}
	if err := cfg.Init(&configFlags, namespace, helmDriver, debug); err != nil {
		log.Fatal(err)
	}
	histClient := action.NewHistory(cfg)
	histClient.Max = 1
	if _, err := histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
		releaseValue, err := install(cfg, releaseName, namespace, chartName, chartVersion, values, dryRun, atomic)
		if err != nil {
			return UpstallResponse{}, fmt.Errorf("error installing: %v", err)
		}
		existingSpecs := make(map[string]*manifest.MappingResult)
		newSpecs := manifest.Parse(releaseValue.Manifest, namespace)
		var buffer strings.Builder
		hasDiff := diff.Manifests(existingSpecs, newSpecs, []string{}, true, 5, &buffer)
		return UpstallResponse{
			ExistingSpecs:         existingSpecs,
			NewSpecs:              newSpecs,
			HasDiff:               hasDiff,
			Diff:                  buffer.String(),
			ChartVersion:          chartVersion,
			CurrentReleaseVersion: int32(releaseValue.Version),
		}, nil
	}
	client := action.NewGet(cfg)
	newRelease, err := upgrade(cfg, releaseName, namespace, chartName, chartVersion, values, dryRun, atomic)

	if err != nil {
		return UpstallResponse{}, fmt.Errorf("error upgrading: %v", err)
	}

	existingRelease, err := client.Run(releaseName)
	if err != nil {
		return UpstallResponse{}, fmt.Errorf("error getting current release: %v", err)
	}
	existingSpecs := manifest.Parse(existingRelease.Manifest, namespace)
	newSpecs := manifest.Parse(newRelease.Manifest, namespace)
	var buffer strings.Builder
	hasDiff := diff.Manifests(existingSpecs, newSpecs, []string{}, true, 5, &buffer)
	return UpstallResponse{
		ExistingSpecs:         existingSpecs,
		NewSpecs:              newSpecs,
		HasDiff:               hasDiff,
		Diff:                  buffer.String(),
		ChartVersion:          chartVersion,
		CurrentReleaseVersion: int32(newRelease.Version),
	}, nil
}

func install(cfg *action.Configuration, releaseName, namespace, chartName, chartVersion, values string, dryRun, atomic bool) (*release.Release, error) {
	client := action.NewInstall(cfg)
	client.DryRun = dryRun
	client.ReleaseName = releaseName
	client.Namespace = namespace
	client.Atomic = atomic
	client.ChartPathOptions.Version = chartVersion

	settings := cli.New()
	cp, err := client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return nil, err
	}

	p := getter.All(settings)

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		warning("This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	valuesMap := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(values), &valuesMap); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", values)
	}
	return client.Run(chartRequested, valuesMap)
}

func upgrade(cfg *action.Configuration, releaseName, namespace, chartName, chartVersion, values string, dryRun, atomic bool) (*release.Release, error) {
	client := action.NewUpgrade(cfg)
	client.DryRun = dryRun
	client.Namespace = namespace
	client.Atomic = atomic
	client.ChartPathOptions.Version = chartVersion

	cp, err := client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		warning("This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			return nil, err
		}
	}

	valuesMap := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(values), &valuesMap); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", values)
	}
	return client.Run(releaseName, chartRequested, valuesMap)
}

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
