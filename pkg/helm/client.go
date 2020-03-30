package helm

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/helmpath"
	"math"
	"os"
	"strings"
	"time"

	"github.com/databus23/helm-diff/diff"
	"github.com/databus23/helm-diff/manifest"
	"github.com/gojek/stevedore/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	helmEnv "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// Client is an abstraction through which helm can be interacted with
type Client interface {
	Upstall(ctx context.Context, releaseName, chartName, chartVersion string, plannedReleaseVersion int32, namespace, values string, dryRun bool, timeout int64) (UpstallResponse, error)
	Close()
}

// DefaultClient is an implementation of helm.Client
type DefaultClient struct {
	helm.Interface
	tillerTunnel *kube.Tunnel
}

// Upstall can install a new release or upgrade if already present
func (c *DefaultClient) Upstall(ctx context.Context, releaseName, chartName, chartVersion string, plannedReleaseVersion int32, namespace, values string, dryRun bool, timeout int64) (UpstallResponse, error) {
	keyring := ""
	chartDownloader := downloader.ChartDownloader{
		HelmHome: helmpath.Home(HomePath()),
		Out:      os.Stdout,
		Keyring:  keyring,
		Getters:  getter.All(helmEnv.EnvSettings{}),
	}
	chartVerifier := DefaultChartVerifier{}
	chartPath, err := LocateChartPath(afero.NewOsFs(), chartName, chartVersion, false, keyring, &chartDownloader, chartVerifier)
	if err != nil {
		return UpstallResponse{}, fmt.Errorf("[Upstall] error when locating chart %s: %v", chartName, err)
	}

	loadedChart, err := SafeLoadChart(chartPath)

	if err != nil {
		return UpstallResponse{}, fmt.Errorf("[Upstall] error when loading chart from path %s: %v", chartPath, err)
	}

	chartVersion = loadedChart.Metadata.Version

	var existingSpecs map[string]*manifest.MappingResult
	var newSpecs map[string]*manifest.MappingResult

	select {
	case <-ctx.Done():
		return UpstallResponse{}, fmt.Errorf("[Upstall] fetching of release details was aborted abruptly")
	default:
	}

	existingReleaseResponse, err := c.ReleaseContent(releaseName)
	var newInstall bool
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("release: %q not found", releaseName)) {
		newInstall = true
	} else if err != nil {
		return UpstallResponse{}, fmt.Errorf("[Upstall] error when checking for existing release: %v", err)
	}

	if newInstall {
		select {
		case <-ctx.Done():
			return UpstallResponse{}, fmt.Errorf("[Upstall] release %s installation aborted abruptly", releaseName)
		default:
		}
		installReleaseResponse, err := c.InstallReleaseFromChart(
			loadedChart,
			namespace,
			helm.ReleaseName(releaseName),
			helm.InstallDryRun(dryRun),
			helm.ValueOverrides([]byte(values)),
			helm.InstallTimeout(timeout),
		)

		if err != nil {
			return UpstallResponse{}, fmt.Errorf("[Upstall] error when installing release %s: %v", releaseName, err)
		}

		newRelease := installReleaseResponse.Release

		existingSpecs = make(map[string]*manifest.MappingResult)
		newSpecs = manifest.ParseRelease(newRelease, true)
	} else {
		if plannedReleaseVersion == 0 {
			plannedReleaseVersion = existingReleaseResponse.GetRelease().GetVersion()
		} else {
			currentReleaseVersion := existingReleaseResponse.GetRelease().GetVersion()
			if plannedReleaseVersion != currentReleaseVersion {
				errMsg := fmt.Sprintf("Release version mismatch: Planned Version - %d, Current version - %d", plannedReleaseVersion, currentReleaseVersion)
				return UpstallResponse{}, fmt.Errorf("[Upstall] error when installing release %s: %v", releaseName, errMsg)
			}
		}

		select {
		case <-ctx.Done():
			return UpstallResponse{}, fmt.Errorf("[Upstall] release %s upgrade aborted abruptly", releaseName)
		default:
		}
		upgradeReleaseResponse, err := c.UpdateReleaseFromChart(
			releaseName,
			loadedChart,
			helm.UpgradeDryRun(dryRun),
			helm.UpdateValueOverrides([]byte(values)),
			helm.UpgradeTimeout(timeout),
		)

		if err != nil {
			return UpstallResponse{}, fmt.Errorf("[Upstall] error when upgrading release %s: %v", releaseName, err)
		}

		existingRelease := existingReleaseResponse.Release
		newRelease := upgradeReleaseResponse.Release

		existingSpecs = manifest.ParseRelease(existingRelease, true)
		newSpecs = manifest.ParseRelease(newRelease, true)
	}

	var buffer strings.Builder
	hasDiff := diff.DiffManifests(existingSpecs, newSpecs, []string{}, 5, &buffer)
	return UpstallResponse{existingSpecs, newSpecs, hasDiff, buffer.String(), chartVersion, plannedReleaseVersion}, nil
}

// Close the tunnel created by the client
func (c *DefaultClient) Close() {
	c.tillerTunnel.Close()
}

// NewHelmClient creates a helm client by loading kube context and port-forwarding to the
// tiller of the given tillerNamespace
func NewHelmClient(tillerNamespace string, client kubernetes.Interface, clientConfig *rest.Config) (Client, error) {
	tillerTunnel, err := portforwarder.New(tillerNamespace, client, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error port-forwarding tiller pod in namespace %s due to %v", tillerNamespace, err)
	}

	tillerHost := fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)
	options := []helm.Option{
		helm.Host(tillerHost),
		helm.ConnectTimeout(int64(10)),
	}

	return &DefaultClient{
		Interface:    helm.NewClient(options...),
		tillerTunnel: tillerTunnel,
	}, nil
}

const maxLoadRetries = 3

// SafeLoadChart will load a chart. If loading fails, it will retry upto maxLoadRetries exponentially backing off for
// each failure
func SafeLoadChart(chartPath string) (*chart.Chart, error) {
	var loadedChart *chart.Chart
	var err = fmt.Errorf("chart %s not loaded", chartPath)

	for i := 0; i < maxLoadRetries; i++ {
		loadedChart, err = chartutil.Load(chartPath)
		if err == nil {
			return loadedChart, nil
		}

		log.Debug(fmt.Sprintf("[safeLoadChart] unable to load chart %s due to %v. retrying...", chartPath, err))
		time.Sleep(time.Duration(math.Pow(2, float64(i))))
	}

	return nil, err
}
