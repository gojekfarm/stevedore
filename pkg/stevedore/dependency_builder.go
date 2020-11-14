package stevedore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"

	pkgHttp "github.com/gojek/stevedore/pkg/http"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/blang/semver"
	"github.com/chartmuseum/helm-push/pkg/helm"
	"github.com/gojek/stevedore/log"
)

var initialVersion = "0.0.1"

// DependencyBuilder will Build the dependency chart and push to Repo
type DependencyBuilder interface {
	Build(ctx context.Context, manifests ManifestFiles) (ManifestFiles, error)
	UpdateRepo() error
	BuildChart(ctx context.Context, releaseSpecification ReleaseSpecification) (ReleaseSpecification, bool, error)
}

// NoopDependencyBuilder will Build the dependency chart and push to it nowhere
type NoopDependencyBuilder struct{}

// Build will Build the dependency chart for given manifests and push to Repo
func (db NoopDependencyBuilder) Build(ctx context.Context, manifests ManifestFiles) (ManifestFiles, error) {
	return manifests, nil
}

// UpdateRepo updates the chart repository
func (db NoopDependencyBuilder) UpdateRepo() error {
	return nil
}

// BuildChart will Build the dependency chart for given releaseSpecification and push to Repo
func (db NoopDependencyBuilder) BuildChart(ctx context.Context, releaseSpecification ReleaseSpecification) (ReleaseSpecification, bool, error) {
	return releaseSpecification, false, nil
}

// DefaultDependencyBuilder will Build the dependency chart and push to Repo
type DefaultDependencyBuilder struct {
	repo            helm.Repo
	client          pkgHttp.Client
	chartBuilder    ChartBuilder
	chartRepository ChartRepository
}

// NewDependencyBuilder will create a DependencyBuilder
func NewDependencyBuilder(repo helm.Repo, client pkgHttp.Client, chartBuilder DefaultChartBuilder, chartRepository repo.ChartRepository) (DependencyBuilder, error) {
	return DefaultDependencyBuilder{repo: repo, client: client, chartBuilder: chartBuilder, chartRepository: &chartRepository}, nil
}

// NewDefaultDependencyBuilder will create a DependencyBuilder (only for test)
func NewDefaultDependencyBuilder(repo helm.Repo, client pkgHttp.Client, chartBuilder ChartBuilder, chartRepository ChartRepository) DependencyBuilder {
	return DefaultDependencyBuilder{repo: repo, client: client, chartBuilder: chartBuilder, chartRepository: chartRepository}
}

// Build will Build the dependency chart for given manifests and push to Repo
func (db DefaultDependencyBuilder) Build(ctx context.Context, manifests ManifestFiles) (ManifestFiles, error) {
	var result ManifestFiles
	for _, manifest := range manifests {
		newManifest := manifest
		newManifest.Spec = ReleaseSpecifications{}
		for _, releaseSpecification := range manifest.Spec {
			releaseSpecification, _, err := db.BuildChart(ctx, releaseSpecification)
			if err != nil {
				return nil, err
			}
			newManifest.Spec = append(newManifest.Spec, releaseSpecification)
		}
		result = append(result, newManifest)
	}
	return result, nil
}

// UpdateRepo updates the chart repository
func (db DefaultDependencyBuilder) UpdateRepo() error {
	_, err := db.chartRepository.DownloadIndexFile() //TODO
	return err
}

// BuildChart will Build the dependency chart for given releaseSpecification and push to Repo
func (db DefaultDependencyBuilder) BuildChart(ctx context.Context, releaseSpecification ReleaseSpecification) (ReleaseSpecification, bool, error) {
	if db.chartBuilder == nil {
		return releaseSpecification, false, fmt.Errorf("unable to build chart, reason chart builder is not defined")
	}

	if !releaseSpecification.Release.HasBuildStep() {
		return releaseSpecification, false, nil
	}

	chartSpec := releaseSpecification.Release.ChartSpec
	dependencies := chartSpec.Dependencies
	if len(dependencies) != 0 {
		chartName := releaseSpecification.Release.ChartSpec.Name
		shouldBuild, version, appVersion, err := db.shouldBuild(chartName, dependencies)
		if err != nil {
			return ReleaseSpecification{}, false, err
		}
		if shouldBuild {
			err := db.chartBuilder.Build(ctx, chartName, version, appVersion, dependencies)
			if err != nil {
				return ReleaseSpecification{}, false, err
			}
		}
		releaseSpecification.Release.Chart = fmt.Sprintf("%s/%s", db.repo.Name, chartName)
		releaseSpecification.Release.ChartVersion = version
	}
	return releaseSpecification, true, nil
}

type chartInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	AppVersion string `json:"appVersion"`
}

func (db DefaultDependencyBuilder) shouldBuild(chartName string, dependencies Dependencies) (bool, string, string, error) {
	if db.client == nil {
		return false, "", "", fmt.Errorf("unable to fetch existing chart details, reason: http client not defined")
	}

	if db.repo.Entry == nil {
		return false, "", "", fmt.Errorf("unable to fetch existing chart details, reason: repo not defined")
	}

	dependenciesCheckSum, err := dependencies.CheckSum()
	if err != nil {
		return false, "", "", err
	}

	parsedURL, err := url.Parse(db.repo.URL)
	if err != nil {
		return false, "", "", err
	}

	parsedURL.Path = path.Join(parsedURL.Path, "api", "charts", chartName)
	response, err := db.client.Get(parsedURL.String())
	if err != nil {
		return false, "", "", err
	}

	if response.StatusCode == 404 {
		log.Debug(fmt.Sprintf("no chart info found for %s", chartName))
		return true, initialVersion, dependenciesCheckSum, nil
	}

	latestCh, err := getLatestChart(response)
	if err != nil {
		log.Debug(fmt.Sprintf("unable to get latest version for %s", chartName))
		return false, "", "", err
	}

	appVersion := latestCh.AppVersion
	if appVersion == dependenciesCheckSum {
		return false, latestCh.Version, appVersion, nil
	}

	version, err := semver.New(latestCh.Version)
	if err != nil {
		log.Debug(fmt.Sprintf("unable to get semver for %s with version %s", chartName, latestCh.Version))
		return false, "", "", err
	}
	return true, nextVersion(version), dependenciesCheckSum, nil
}

func getLatestChart(response *http.Response) (chartInfo, error) {
	var charts []chartInfo
	if response.Body == nil {
		return chartInfo{}, fmt.Errorf("invalid response from chart repository")
	}

	err := json.NewDecoder(response.Body).Decode(&charts)
	if err != nil {
		return chartInfo{}, fmt.Errorf("unable to parse chart information %v", err.Error())
	}
	sort.Slice(charts, func(i, j int) bool {
		return charts[i].Version > charts[j].Version
	})
	return charts[0], nil
}

func nextVersion(version *semver.Version) string {
	major := version.Major
	minor := version.Minor
	patch := version.Patch
	if patch < 10 {
		patch++
	} else if minor < 10 {
		patch = 0
		minor++
	} else {
		major++
		minor = 0
		patch = 0
	}
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
