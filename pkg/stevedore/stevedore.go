package stevedore

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	chartMuseumHelm "github.com/chartmuseum/helm-push/pkg/helm"
	"github.com/gojek/stevedore/pkg/utils"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/gojek/stevedore/log"
	"github.com/gojek/stevedore/pkg/helm"
)

// Opts represents options to stevedore release
type Opts struct {
	DryRun   bool
	Parallel bool
	Filter   bool
	Timeout  int64
}

// Stevedore installs or upgrades helm releases
type Stevedore struct {
	helm.Client
	Opts
	Upstaller
	DependencyBuilder
}

// CreateResponse will take the manifests and helmClients and produce response based on given Opts
func CreateResponse(ctx context.Context, manifestFiles ManifestFiles, opts Opts, helmRepoName string, helmTimeout int64) (Responses, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("request aborted abruptly by client")
	default:
		dependencyBuilder, err := CreateDependencyBuilder(manifestFiles, helmRepoName)
		if err != nil {
			return nil, err
		}
		s := Stevedore{Client: &helm.DefaultClient{}, Opts: opts, Upstaller: HelmUpstaller{}, DependencyBuilder: dependencyBuilder}
		responses, _ := s.Do(ctx, manifestFiles, helmTimeout)
		return responses, nil
	}
}

// CreateDependencyBuilder for the given chart repo name
func CreateDependencyBuilder(manifestFiles ManifestFiles, chartRepoName string) (DependencyBuilder, error) {
	if !manifestFiles.HasBuildStep() {
		return NoopDependencyBuilder{}, nil
	}

	matchingRepo, err := chartMuseumHelm.GetRepoByName(chartRepoName)
	if err != nil {
		return NoopDependencyBuilder{}, fmt.Errorf("%s", noChartRepoError(chartRepoName))
	}

	chartBuilder := NewChartBuilder(*matchingRepo, DefaultChartManager{}, utils.NewOsFileUtils())
	settings := cli.EnvSettings{}
	entry := repo.Entry{Name: matchingRepo.Name, URL: matchingRepo.URL, Username: matchingRepo.Username, Password: matchingRepo.Password, CertFile: matchingRepo.CertFile, KeyFile: matchingRepo.KeyFile, CAFile: matchingRepo.CAFile}
	chartRepository, err := repo.NewChartRepository(&entry, getter.All(&settings))

	if err != nil {
		return NoopDependencyBuilder{}, err
	}
	return NewDependencyBuilder(*matchingRepo, &http.Client{}, chartBuilder, *chartRepository)
}

func noChartRepoError(repoName string) string {
	buff := bytes.NewBufferString(fmt.Sprintf("unable to find helm chart repo with name '%s'\n", repoName))
	buff.WriteString(fmt.Sprintf("to add the repo '%s', run the following command\n", repoName))
	buff.WriteString(fmt.Sprintf("\n$ helm repo add %s <your chart muesum url>\n", repoName))
	buff.WriteString(fmt.Sprintf("\nto use different helm repo other than '%s', specify it with '-r <repo name>'\n", repoName))

	return buff.String()
}

// Do install or upgrade releases
func (s Stevedore) Do(ctx context.Context, manifestFiles ManifestFiles, helmTimeout int64) (Responses, error) {
	responseCh := make(chan Response)
	go s.createResponses(ctx, manifestFiles, responseCh, helmTimeout)

	var acc Responses
	for response := range responseCh {
		if !s.Filter || (response.HasDiff || response.Err != nil) {
			acc = append(acc, response)
		}

		if response.Err != nil && !s.DryRun && !s.Parallel {
			return acc, response.Err
		}
	}

	return acc, nil
}

func (s Stevedore) createResponses(ctx context.Context, manifestFiles ManifestFiles, responseCh chan<- Response, helmTimeout int64) {
	var wg sync.WaitGroup
	proceed := make(chan bool)
	s.response(ctx, manifestFiles, &wg, responseCh, proceed, helmTimeout)
	wg.Wait()
	close(responseCh)
}

func (s Stevedore) buildChartIfNeeded(ctx context.Context, releaseSpecification ReleaseSpecification) (ReleaseSpecification, error) {
	if s.DependencyBuilder == nil {
		return releaseSpecification, nil
	}
	builtApplication, isChartBuilt, err := s.DependencyBuilder.BuildChart(ctx, releaseSpecification)
	if err != nil {
		return releaseSpecification, fmt.Errorf("failed to Build chart with error: %s", err.Error())
	}
	if isChartBuilt {
		err = s.DependencyBuilder.UpdateRepo()
		if err != nil {
			return releaseSpecification, fmt.Errorf("error updating helm repo with error: %s", err.Error())
		}
	}
	return builtApplication, nil
}

func (s Stevedore) response(ctx context.Context, manifestFiles ManifestFiles, wg *sync.WaitGroup, responseCh chan<- Response, proceed chan bool, helmTimeout int64) {
	for _, request := range manifestFiles {
		for _, releaseSpecification := range request.Manifest.Spec {
			releaseSpecification, err := s.buildChartIfNeeded(ctx, releaseSpecification)
			if err != nil {
				log.Error(err.Error())
				responseCh <- Response{
					request.File,
					releaseSpecification.Release.Name,
					releaseSpecification.Release.Chart,
					releaseSpecification.Release.ChartVersion,
					releaseSpecification.Release.CurrentReleaseVersion,
					helm.UpstallResponse{},
					fmt.Errorf("%v", err.Error()),
				}
				continue
			}
			wg.Add(1)
			go s.Upstaller.Upstall(ctx, s.Client, releaseSpecification, request.File, responseCh, proceed, wg, s.Opts, helmTimeout)
			if !<-proceed {
				return
			}
		}
	}
}
