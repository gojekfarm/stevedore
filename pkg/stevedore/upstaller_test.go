package stevedore_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/databus23/helm-diff/manifest"
	"github.com/gojek/stevedore/pkg/helm"
	mocks "github.com/gojek/stevedore/pkg/internal/mocks/helm"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpstaller_upstall(t *testing.T) {
	var currentReleaseVersion int32
	chartVersion := ""
	var timeout int64 = 10
	var atomic = true
	t.Run("should create helm response when helm upstall succeeds", func(t *testing.T) {
		t.Run("with parallel true", func(t *testing.T) {
			opts := stevedore.Opts{Parallel: true, DryRun: true}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			release := stevedore.Release{
				Name:      "postgres",
				Namespace: "default",
				Chart:     "stable/postgresql",
				Values: stevedore.Values{
					"image": "postgresql:10",
				},
			}

			valuesYaml, _ := release.Values.ToYAML()
			upstallResponse := helm.UpstallResponse{
				ExistingSpecs: map[string]*manifest.MappingResult{"key": {}},
				NewSpecs:      map[string]*manifest.MappingResult{"key": {}},
				HasDiff:       false,
				Diff:          "",
			}
			client.EXPECT().Upstall(context.TODO(),
				release.Name,
				release.Chart,
				chartVersion,
				currentReleaseVersion,
				release.Namespace,
				valuesYaml,
				opts.DryRun,
				timeout,
				atomic).Return(upstallResponse, nil)

			responseCh := make(chan stevedore.Response)
			wg := sync.WaitGroup{}
			wg.Add(1)
			proceed := make(chan bool)
			releaseSpecification := stevedore.ReleaseSpecification{
				Release: release,
			}
			go stevedore.HelmUpstaller{}.Upstall(context.TODO(), client, releaseSpecification, "postgres.yaml", responseCh, proceed, &wg, opts, timeout, atomic)

			wgForAssert := sync.WaitGroup{}
			wgForAssert.Add(1)
			go func(group *sync.WaitGroup) {
				assert.True(t, <-proceed)
				wgForAssert.Done()
			}(&wgForAssert)

			go func(wg *sync.WaitGroup) {
				wg.Wait()
				close(responseCh)
			}(&wg)

			var responses stevedore.Responses
			for response := range responseCh {
				responses = append(responses, response)
			}

			expectedResponse := stevedore.Response{
				File:            "postgres.yaml",
				ChartName:       "stable/postgresql",
				ReleaseName:     "postgres",
				UpstallResponse: upstallResponse,
				Err:             nil,
			}
			assert.Len(t, responses, 1)
			assert.NotNil(t, responses)
			assert.Equal(t, responses[0], expectedResponse)
			wgForAssert.Wait()
		})

		t.Run("with parallel false", func(t *testing.T) {
			opts := stevedore.Opts{Parallel: false, DryRun: true}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			release := stevedore.Release{
				Name:      "postgres",
				Namespace: "default",
				Chart:     "stable/postgresql",
				Values: stevedore.Values{
					"image": "postgresql:10",
				},
			}
			valuesYaml, _ := release.Values.ToYAML()
			upstallResponse := helm.UpstallResponse{
				ExistingSpecs: map[string]*manifest.MappingResult{"key": {}},
				NewSpecs:      map[string]*manifest.MappingResult{"key": {}},
				HasDiff:       false,
				Diff:          "",
			}
			client.EXPECT().Upstall(context.TODO(),
				release.Name,
				release.Chart,
				chartVersion,
				currentReleaseVersion,
				release.Namespace,
				valuesYaml,
				opts.DryRun,
				timeout,
				atomic).Return(upstallResponse, nil)

			responseCh := make(chan stevedore.Response)
			wg := sync.WaitGroup{}
			wg.Add(1)
			proceed := make(chan bool)
			releaseSpecification := stevedore.ReleaseSpecification{
				Release: release,
			}
			go stevedore.HelmUpstaller{}.Upstall(context.TODO(), client, releaseSpecification, "postgres.yaml", responseCh, proceed, &wg, opts, timeout, atomic)

			wgForAssert := sync.WaitGroup{}
			wgForAssert.Add(1)
			go func(group *sync.WaitGroup) {
				assert.True(t, <-proceed)
				wgForAssert.Done()
			}(&wgForAssert)

			go func(wg *sync.WaitGroup) {
				wg.Wait()
				close(responseCh)
			}(&wg)

			var responses stevedore.Responses
			for response := range responseCh {
				responses = append(responses, response)
			}

			expectedResponse := stevedore.Response{
				File:            "postgres.yaml",
				ReleaseName:     "postgres",
				ChartName:       "stable/postgresql",
				UpstallResponse: upstallResponse,
				Err:             nil,
			}
			assert.Len(t, responses, 1)
			assert.Equal(t, responses[0], expectedResponse)
			wgForAssert.Wait()
		})
	})

	t.Run("proceed should be populated before response creation when parallel is enabled", func(t *testing.T) {
		opts := stevedore.Opts{Parallel: true, DryRun: true}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockClient(ctrl)
		release := stevedore.Release{
			Name:      "postgres",
			Namespace: "default",
			Chart:     "stable/postgresql",
			Values: stevedore.Values{
				"image": "postgresql:10",
			},
		}
		valuesYaml, _ := release.Values.ToYAML()
		upstallResponse := helm.UpstallResponse{
			ExistingSpecs: map[string]*manifest.MappingResult{"key": {}},
			NewSpecs:      map[string]*manifest.MappingResult{"key": {}},
			HasDiff:       false,
			Diff:          "",
		}
		client.EXPECT().Upstall(context.TODO(),
			release.Name,
			release.Chart,
			chartVersion,
			currentReleaseVersion,
			release.Namespace,
			valuesYaml,
			opts.DryRun,
			timeout,
			atomic).Return(upstallResponse, nil)

		responseCh := make(chan stevedore.Response)
		wg := sync.WaitGroup{}
		wg.Add(1)
		proceed := make(chan bool)
		releaseSpecification := stevedore.ReleaseSpecification{
			Release: release,
		}
		go stevedore.HelmUpstaller{}.Upstall(context.TODO(), client, releaseSpecification, "postgres.yaml", responseCh, proceed, &wg, opts, timeout, atomic)

		<-proceed

		go func(wg *sync.WaitGroup) {
			wg.Wait()
			close(responseCh)
		}(&wg)

		var responses stevedore.Responses
		for response := range responseCh {
			responses = append(responses, response)
		}

		expectedResponse := stevedore.Response{
			File:            "postgres.yaml",
			ChartName:       "stable/postgresql",
			ReleaseName:     "postgres",
			UpstallResponse: upstallResponse,
			Err:             nil,
		}
		assert.Len(t, responses, 1)
		assert.Equal(t, responses[0], expectedResponse)
	})

	t.Run("should add error to response when helm upstall fails", func(t *testing.T) {
		t.Run("with parallel true", func(t *testing.T) {
			opts := stevedore.Opts{Parallel: true, DryRun: true}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			release := stevedore.Release{
				Name:      "postgres",
				Namespace: "default",
				Chart:     "stable/postgresql",
				Values: stevedore.Values{
					"image": "postgresql:10",
				},
			}

			valuesYaml, _ := release.Values.ToYAML()
			client.EXPECT().Upstall(context.TODO(),
				release.Name,
				release.Chart,
				chartVersion,
				currentReleaseVersion,
				release.Namespace,
				valuesYaml,
				opts.DryRun,
				timeout, atomic).Return(helm.UpstallResponse{}, fmt.Errorf("something went wrong"))

			responseCh := make(chan stevedore.Response)
			wg := sync.WaitGroup{}
			wg.Add(1)
			proceed := make(chan bool)
			releaseSpecification := stevedore.ReleaseSpecification{
				Release: release,
			}
			go stevedore.HelmUpstaller{}.Upstall(context.TODO(), client, releaseSpecification, "postgres.yaml", responseCh, proceed, &wg, opts, timeout, atomic)

			wgForAssert := sync.WaitGroup{}
			wgForAssert.Add(1)
			go func(group *sync.WaitGroup) {
				assert.True(t, <-proceed)
				wgForAssert.Done()
			}(&wgForAssert)

			go func(wg *sync.WaitGroup) {
				wg.Wait()
				close(responseCh)
			}(&wg)

			var responses stevedore.Responses
			for response := range responseCh {
				responses = append(responses, response)
			}

			expectedResponse := stevedore.Response{
				File:            "postgres.yaml",
				ChartName:       "stable/postgresql",
				ReleaseName:     "postgres",
				UpstallResponse: helm.UpstallResponse{},
				Err:             fmt.Errorf("error when installing postgres due to something went wrong"),
			}
			assert.Len(t, responses, 1)
			assert.Equal(t, responses[0], expectedResponse)
			wgForAssert.Wait()
		})

		t.Run("with parallel false", func(t *testing.T) {
			opts := stevedore.Opts{Parallel: false, DryRun: true}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			release := stevedore.Release{
				Name:      "postgres",
				Namespace: "default",
				Chart:     "stable/postgresql",
				Values: stevedore.Values{
					"image": "postgresql:10",
				},
			}

			valuesYaml, _ := release.Values.ToYAML()
			client.EXPECT().Upstall(context.TODO(),
				release.Name,
				release.Chart,
				chartVersion,
				currentReleaseVersion,
				release.Namespace,
				valuesYaml,
				opts.DryRun,
				timeout, atomic).Return(helm.UpstallResponse{}, fmt.Errorf("something went wrong"))

			responseCh := make(chan stevedore.Response)
			wg := sync.WaitGroup{}
			wg.Add(1)
			proceed := make(chan bool)
			releaseSpecification := stevedore.ReleaseSpecification{
				Release: release,
			}
			go stevedore.HelmUpstaller{}.Upstall(context.TODO(), client, releaseSpecification, "postgres.yaml", responseCh, proceed, &wg, opts, timeout, atomic)

			wgForAssert := sync.WaitGroup{}
			wgForAssert.Add(1)
			go func(group *sync.WaitGroup) {
				assert.False(t, <-proceed)
				wgForAssert.Done()
			}(&wgForAssert)

			go func(wg *sync.WaitGroup) {
				wg.Wait()
				close(responseCh)
			}(&wg)

			var responses stevedore.Responses
			for response := range responseCh {
				responses = append(responses, response)
			}

			expectedResponse := stevedore.Response{
				File:            "postgres.yaml",
				ChartName:       "stable/postgresql",
				ReleaseName:     "postgres",
				UpstallResponse: helm.UpstallResponse{},
				Err:             fmt.Errorf("error when installing postgres due to something went wrong"),
			}
			assert.Len(t, responses, 1)
			assert.Equal(t, responses[0], expectedResponse)
			wgForAssert.Wait()
		})
	})
}
