package stevedore

import (
	"bytes"
	"context"
	"fmt"
	chartMuseumHelm "github.com/chartmuseum/helm-push/pkg/helm"
	"github.com/gojek/stevedore/pkg/utils"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/repo"
	"net/http"
	"sync"

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
	helm.Clients
	Opts
	Upstaller
	DependencyBuilder
}

// CreateResponse will take the manifests and helmClients and produce response based on given Opts
func CreateResponse(ctx context.Context, manifestFiles ManifestFiles, opts Opts, helmRepoName string, createHelmClient func(namespaces []string) (helm.Clients, error), helmTimeout int64) (Responses, error) {
	namespaces := manifestFiles.AllNamespaces()
	helmClients, err := createHelmClient(namespaces)
	if err != nil {
		return nil, fmt.Errorf("error creating helm clients due to %v", err)
	}
	closeHelmClients := func() {
		helmClients.Close()
		log.Debug("Gracefully closed all helm clients")
	}
	defer closeHelmClients()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("request aborted abruptly by client")
	default:
		dependencyBuilder, err := CreateDependencyBuilder(manifestFiles, helmRepoName)
		if err != nil {
			return nil, err
		}
		s := Stevedore{Clients: helmClients, Opts: opts, Upstaller: HelmUpstaller{}, DependencyBuilder: dependencyBuilder}
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
	settings := environment.EnvSettings{}
	entry := repo.Entry{Name: matchingRepo.Name, Cache: matchingRepo.Cache, URL: matchingRepo.URL, Username: matchingRepo.Username, Password: matchingRepo.Password, CertFile: matchingRepo.CertFile, KeyFile: matchingRepo.KeyFile, CAFile: matchingRepo.CAFile}
	chartRepository, err := repo.NewChartRepository(&entry, getter.All(settings))

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
			namespace := releaseSpecification.Release.TillerNamespace()
			if client, ok := s.Clients[namespace]; ok {
				wg.Add(1)
				go s.Upstaller.Upstall(ctx, client, releaseSpecification, request.File, responseCh, proceed, wg, s.Opts, helmTimeout)
				if !<-proceed {
					return
				}
			} else {
				err := fmt.Errorf("unable to retrieve helm client for namespace %s", namespace)
				log.Error(err.Error())
				responseCh <- Response{
					request.File,
					releaseSpecification.Release.Name,
					releaseSpecification.Release.Chart,
					releaseSpecification.Release.ChartVersion,
					releaseSpecification.Release.CurrentReleaseVersion,
					helm.UpstallResponse{},
					err,
				}
				continue
			}
		}
	}
}
