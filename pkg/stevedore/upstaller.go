package stevedore

import (
	"context"
	"fmt"
	"sync"

	"github.com/gojek/stevedore/pkg/helm"
)

// Upstaller will release/upgrade a releaseSpecification
type Upstaller interface {
	Upstall(ctx context.Context, client helm.Client, releaseSpecification ReleaseSpecification, file string,
		responseCh chan<- Response, proceed chan<- bool, wg *sync.WaitGroup, opts Opts, helmTimeout int64, helmAtomic bool)
}

// HelmUpstaller will release/upgrade a releaseSpecification using helmClient
type HelmUpstaller struct{}

// Upstall will release/upgrade a releaseSpecification using helmClient
func (HelmUpstaller) Upstall(ctx context.Context, client helm.Client, releaseSpecification ReleaseSpecification, file string,
	responseCh chan<- Response, proceed chan<- bool, wg *sync.WaitGroup, opts Opts, helmTimeout int64, helmAtomic bool) {
	defer wg.Done()
	var err error
	if opts.Parallel {
		proceed <- true
	} else {
		defer func() {
			proceed <- err == nil
		}()
	}
	var upstallResponse helm.UpstallResponse
	manifestName := releaseSpecification.Release.Name
	namespace := releaseSpecification.Release.Namespace
	chartName := releaseSpecification.Release.Chart
	chartVersion := releaseSpecification.Release.ChartVersion
	currentReleaseVersion := releaseSpecification.Release.CurrentReleaseVersion
	values, _ := releaseSpecification.Release.Values.ToYAML()
	upstallResponse, err = client.Upstall(ctx, manifestName, chartName, chartVersion, currentReleaseVersion, namespace, values, opts.DryRun, helmTimeout, helmAtomic)
	chartVersion = upstallResponse.ChartVersion
	newCurrentReleaseVersion := upstallResponse.CurrentReleaseVersion
	if err != nil {
		responseCh <- Response{
			file,
			manifestName,
			chartName,
			chartVersion,
			newCurrentReleaseVersion,
			upstallResponse,
			fmt.Errorf("error when installing %s due to %v", manifestName, err),
		}
		return
	}

	responseCh <- Response{
		file,
		manifestName,
		chartName,
		chartVersion,
		newCurrentReleaseVersion,
		upstallResponse,
		nil,
	}
}
