package helm

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/provenance"
	"os"
	"path/filepath"
	"strings"

	// load the gcp plugin (required to authenticate against GKE clusters)
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/util/homedir"
	"k8s.io/helm/pkg/helm/helmpath"
)

// HomePath returns the home path to HELM
func HomePath() string {
	if homePath := os.Getenv("HELM_HOME"); homePath != "" {
		return homePath
	}
	return filepath.Join(homedir.HomeDir(), ".helm")
}

// ChartDownloader represent interface to download chart
type ChartDownloader interface {
	DownloadTo(ref, version, dest string) (string, *provenance.Verification, error)
}

// ChartVerifier represents interface to verify chart
type ChartVerifier interface {
	VerifyChart(path string, keyring string) (*provenance.Verification, error)
}

// DefaultChartVerifier represents default chart verifier
type DefaultChartVerifier struct{}

// VerifyChart verifies whether chart exists in the given path and keyring
func (d DefaultChartVerifier) VerifyChart(path string, keyring string) (*provenance.Verification, error) {
	return downloader.VerifyChart(path, keyring)
}

// LocateChartPath makes sure a chart archive is present given a chart name and version
func LocateChartPath(fs afero.Fs, name, version string, verify bool, keyring string, chartDownloader ChartDownloader, verifier ChartVerifier) (string, error) {
	name = strings.TrimSpace(name)
	version = strings.TrimSpace(version)

	if fi, err := fs.Stat(name); err == nil {
		abs, err := filepath.Abs(name)
		if err != nil {
			return abs, err
		}
		if verify {
			if fi.IsDir() {
				return "", errors.New("cannot verify a directory")
			}
			if _, err := verifier.VerifyChart(abs, keyring); err != nil {
				return "", err
			}
		}
		return abs, nil
	}
	if filepath.IsAbs(name) || strings.HasPrefix(name, ".") {
		return name, fmt.Errorf("path %q not found", name)
	}

	chartRepo := filepath.Join(helmpath.Home(HomePath()).Repository(), name)
	if _, err := fs.Stat(chartRepo); err == nil {
		return filepath.Abs(chartRepo)
	}

	filename, _, err := chartDownloader.DownloadTo(name, version, helmpath.Home(HomePath()).Archive())
	if err == nil {
		abs, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return abs, nil
	}

	return filename, err
}
