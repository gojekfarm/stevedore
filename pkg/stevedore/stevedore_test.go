package stevedore_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/gojek/stevedore/pkg/helm"
	mockDependencyBuilder "github.com/gojek/stevedore/pkg/internal/mocks/chart"
	mocks "github.com/gojek/stevedore/pkg/internal/mocks/helm"
	mockUpstaller "github.com/gojek/stevedore/pkg/internal/mocks/upstaller"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStevedoreDo(t *testing.T) {
	var timeout int64 = 10

	t.Run("should do nothing for empty releaseSpecifications", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockClient(ctrl)
		upstaller := mockUpstaller.NewMockUpstaller(ctrl)

		opts := stevedore.Opts{DryRun: true, Parallel: true}
		s := stevedore.Stevedore{
			Client:    client,
			Opts:      opts,
			Upstaller: upstaller,
		}
		responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{}, timeout)

		assert.Empty(t, responses)
		assert.Nil(t, err)
	})

	t.Run("should create response for all releaseSpecifications when filter is disabled", func(t *testing.T) {
		releaseSpecificationX := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)
		releaseSpecificationY := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("y-stevedore", "default", "chart/y-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "y-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)

		releaseSpecifications := stevedore.ReleaseSpecifications{
			releaseSpecificationX,
			releaseSpecificationY,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockClient(ctrl)
		upstaller := mockUpstaller.NewMockUpstaller(ctrl)
		dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationX).Times(2).Return(releaseSpecificationX, true, nil)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationY).Times(2).Return(releaseSpecificationY, true, nil)
		dependencyBuilder.EXPECT().UpdateRepo().Times(4).Return(nil)

		opts := stevedore.Opts{DryRun: true, Parallel: true, Filter: false}

		upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, fileName, responseCh, proceedCh, wg, _ interface{}, t int64) {
			proceedCh.(chan<- bool) <- true
			defer wg.(*sync.WaitGroup).Done()

			response := stevedore.Response{
				File:            fileName.(string),
				ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			}
			responseCh.(chan<- stevedore.Response) <- response
		}).Times(4)

		expectedResponses := stevedore.Responses{
			{
				File:            "postgres.yaml",
				ReleaseName:     releaseSpecificationX.Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "postgres.yaml",
				ReleaseName:     releaseSpecificationY.Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "redis.yaml",
				ReleaseName:     releaseSpecificationX.Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "redis.yaml",
				ReleaseName:     releaseSpecificationY.Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
		}

		s := stevedore.Stevedore{
			Client:            client,
			Opts:              opts,
			Upstaller:         upstaller,
			DependencyBuilder: dependencyBuilder,
		}
		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
			Spec:     releaseSpecifications,
		}
		postgresRequest := stevedore.ManifestFile{File: "postgres.yaml", Manifest: manifest}
		redisRequest := stevedore.ManifestFile{File: "redis.yaml", Manifest: manifest}
		responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{postgresRequest, redisRequest}, timeout)

		assert.Len(t, responses, 4)
		assert.ElementsMatch(t, expectedResponses, responses)
		assert.Nil(t, err)
	})

	t.Run("should create response only for releaseSpecifications having diff when filter is enabled", func(t *testing.T) {
		releaseSpecificationX := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)
		releaseSpecificationY := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("y-stevedore", "default", "chart/y-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "y-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)

		releaseSpecifications := stevedore.ReleaseSpecifications{
			releaseSpecificationX,
			releaseSpecificationY,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockClient(ctrl)
		upstaller := mockUpstaller.NewMockUpstaller(ctrl)
		dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationX).Times(2).Return(releaseSpecificationX, true, nil)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationY).Times(2).Return(releaseSpecificationY, true, nil)
		dependencyBuilder.EXPECT().UpdateRepo().Times(4).Return(nil)

		opts := stevedore.Opts{DryRun: true, Parallel: true, Filter: true}

		upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, fileName, responseCh, proceedCh, wg, _ interface{}, t int64) {
			proceedCh.(chan<- bool) <- true
			defer wg.(*sync.WaitGroup).Done()

			componentName := releaseSpecification.(stevedore.ReleaseSpecification).Release.Name
			response := stevedore.Response{
				File:            fileName.(string),
				ReleaseName:     componentName,
				UpstallResponse: helm.UpstallResponse{HasDiff: false},
				Err:             nil,
			}
			if componentName == "x-stevedore" {
				response.UpstallResponse = helm.UpstallResponse{HasDiff: true}
			}
			responseCh.(chan<- stevedore.Response) <- response
		}).Times(4)

		expectedResponses := stevedore.Responses{
			{
				File:            "postgres.yaml",
				ReleaseName:     releaseSpecificationX.Release.Name,
				UpstallResponse: helm.UpstallResponse{HasDiff: true},
				Err:             nil,
			},
			{
				File:            "redis.yaml",
				ReleaseName:     releaseSpecificationX.Release.Name,
				UpstallResponse: helm.UpstallResponse{HasDiff: true},
				Err:             nil,
			},
		}

		s := stevedore.Stevedore{
			Client:            client,
			Opts:              opts,
			Upstaller:         upstaller,
			DependencyBuilder: dependencyBuilder,
		}
		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
			Spec:     releaseSpecifications,
		}
		postgresRequest := stevedore.ManifestFile{File: "postgres.yaml", Manifest: manifest}
		redisRequest := stevedore.ManifestFile{File: "redis.yaml", Manifest: manifest}
		responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{postgresRequest, redisRequest}, timeout)

		assert.Len(t, responses, 2)
		assert.ElementsMatch(t, expectedResponses, responses)
		assert.Nil(t, err)
	})

	t.Run("when there is an error", func(t *testing.T) {
		t.Run("should abort for non parallel and non dry run", func(t *testing.T) {
			releaseSpecificationOne := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationTwo := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationThree := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)

			releaseSpecifications := stevedore.ReleaseSpecifications{
				releaseSpecificationOne,
				releaseSpecificationTwo,
				releaseSpecificationThree,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			upstaller := mockUpstaller.NewMockUpstaller(ctrl)
			dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationOne).Return(releaseSpecificationOne, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationTwo).Return(releaseSpecificationTwo, true, nil)
			dependencyBuilder.EXPECT().UpdateRepo().Times(2).Return(nil)

			file := "postgres.yaml"
			opts := stevedore.Opts{DryRun: false, Parallel: false}
			counter := 0
			upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), file, gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, _, responseCh, proceedCh, wg, _ interface{}, t int64) {
				counter++
				defer wg.(*sync.WaitGroup).Done()

				var err error
				if counter == 2 {
					err = fmt.Errorf("error")
				}
				proceedCh.(chan<- bool) <- counter != 2
				response := stevedore.Response{
					File:            file,
					ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             err,
				}
				responseCh.(chan<- stevedore.Response) <- response
			}).Times(2)

			expectedResponses := stevedore.Responses{
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationOne.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationTwo.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             fmt.Errorf("error"),
				},
			}

			s := stevedore.Stevedore{
				Client:            client,
				Opts:              opts,
				Upstaller:         upstaller,
				DependencyBuilder: dependencyBuilder,
			}
			manifest := stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
				Spec:     releaseSpecifications,
			}
			manifestFile := stevedore.ManifestFile{File: "postgres.yaml", Manifest: manifest}
			responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{manifestFile}, timeout)

			assert.Len(t, responses, 2)
			assert.Equal(t, expectedResponses, responses)
			assert.Error(t, err)
		})

		t.Run("should continue for non parallel and dry run", func(t *testing.T) {
			releaseSpecificationOne := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationTwo := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationThree := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)

			releaseSpecifications := stevedore.ReleaseSpecifications{
				releaseSpecificationOne,
				releaseSpecificationTwo,
				releaseSpecificationThree,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			upstaller := mockUpstaller.NewMockUpstaller(ctrl)
			dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationOne).Return(releaseSpecificationOne, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationTwo).Return(releaseSpecificationTwo, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationThree).Return(releaseSpecificationThree, true, nil)
			dependencyBuilder.EXPECT().UpdateRepo().Times(3).Return(nil)

			file := "postgres.yaml"
			opts := stevedore.Opts{DryRun: true, Parallel: false}
			counter := 0
			upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), file, gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, _, responseCh, proceedCh, wg, _ interface{}, t int64) {
				counter++
				defer wg.(*sync.WaitGroup).Done()

				var err error
				if counter == 2 {
					err = fmt.Errorf("error")
				}
				proceedCh.(chan<- bool) <- true
				response := stevedore.Response{
					File:            file,
					ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             err,
				}
				responseCh.(chan<- stevedore.Response) <- response
			}).Times(3)

			expectedResponses := stevedore.Responses{
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationOne.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationTwo.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             fmt.Errorf("error"),
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationThree.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
			}

			s := stevedore.Stevedore{
				Client:            client,
				Opts:              opts,
				Upstaller:         upstaller,
				DependencyBuilder: dependencyBuilder,
			}
			manifest := stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
				Spec:     releaseSpecifications,
			}
			manifestFile := stevedore.ManifestFile{File: "postgres.yaml", Manifest: manifest}
			responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{manifestFile}, timeout)

			assert.Len(t, responses, 3)
			assert.Equal(t, expectedResponses, responses)
			assert.Nil(t, err)
		})

		t.Run("should continue for parallel and non dry run", func(t *testing.T) {
			releaseSpecificationOne := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationTwo := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationThree := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)

			releaseSpecifications := stevedore.ReleaseSpecifications{
				releaseSpecificationOne,
				releaseSpecificationTwo,
				releaseSpecificationThree,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			upstaller := mockUpstaller.NewMockUpstaller(ctrl)
			dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationOne).Return(releaseSpecificationOne, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationTwo).Return(releaseSpecificationTwo, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationThree).Return(releaseSpecificationThree, true, nil)
			dependencyBuilder.EXPECT().UpdateRepo().Times(3).Return(nil)

			file := "postgres.yaml"
			opts := stevedore.Opts{DryRun: false, Parallel: true}
			counter := 0
			upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), file, gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, _, responseCh, proceedCh, wg, _ interface{}, t int64) {
				counter++
				defer wg.(*sync.WaitGroup).Done()

				var err error
				if counter == 2 {
					err = fmt.Errorf("error")
				}
				proceedCh.(chan<- bool) <- true
				response := stevedore.Response{
					File:            file,
					ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             err,
				}
				responseCh.(chan<- stevedore.Response) <- response
			}).Times(3)

			expectedResponses := stevedore.Responses{
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationOne.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationTwo.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             fmt.Errorf("error"),
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationThree.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
			}

			s := stevedore.Stevedore{
				Client:            client,
				Opts:              opts,
				Upstaller:         upstaller,
				DependencyBuilder: dependencyBuilder,
			}
			manifest := stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
				Spec:     releaseSpecifications,
			}
			manifestFile := stevedore.ManifestFile{File: "postgres.yaml", Manifest: manifest}
			responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{manifestFile}, timeout)

			assert.Len(t, responses, 3)
			assert.Equal(t, expectedResponses, responses)
			assert.Nil(t, err)
		})

		t.Run("should continue for parallel and dry run", func(t *testing.T) {
			releaseSpecificationOne := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationTwo := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)
			releaseSpecificationThree := stevedore.NewReleaseSpecification(
				stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
				stevedore.Configs{
					"store": []map[string]interface{}{
						{"name": "x-stevedore", "tags": []string{"server"}},
						{"name": "store-ns"},
					},
				}, nil)

			releaseSpecifications := stevedore.ReleaseSpecifications{
				releaseSpecificationOne,
				releaseSpecificationTwo,
				releaseSpecificationThree,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockClient(ctrl)
			upstaller := mockUpstaller.NewMockUpstaller(ctrl)
			dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationOne).Return(releaseSpecificationOne, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationTwo).Return(releaseSpecificationTwo, true, nil)
			dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationThree).Return(releaseSpecificationThree, true, nil)
			dependencyBuilder.EXPECT().UpdateRepo().Times(3).Return(nil)

			file := "postgres.yaml"
			opts := stevedore.Opts{DryRun: true, Parallel: true}
			counter := 0
			upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), file, gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, _, responseCh, proceedCh, wg, _ interface{}, t int64) {
				counter++
				defer wg.(*sync.WaitGroup).Done()

				var err error
				if counter == 2 {
					err = fmt.Errorf("error")
				}
				proceedCh.(chan<- bool) <- true
				response := stevedore.Response{
					File:            file,
					ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             err,
				}
				responseCh.(chan<- stevedore.Response) <- response
			}).Times(3)

			expectedResponses := stevedore.Responses{
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationOne.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationTwo.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             fmt.Errorf("error"),
				},
				{
					File:            "postgres.yaml",
					ReleaseName:     releaseSpecificationThree.Release.Name,
					UpstallResponse: helm.UpstallResponse{},
					Err:             nil,
				},
			}

			s := stevedore.Stevedore{
				Client:            client,
				Opts:              opts,
				Upstaller:         upstaller,
				DependencyBuilder: dependencyBuilder,
			}
			manifest := stevedore.Manifest{
				DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
				Spec:     releaseSpecifications,
			}
			manifestFile := stevedore.ManifestFile{File: "postgres.yaml", Manifest: manifest}
			responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{manifestFile}, timeout)

			assert.Len(t, responses, 3)
			assert.Equal(t, expectedResponses, responses)
			assert.Nil(t, err)
		})
	})

	t.Run("should return error when chart build fails", func(t *testing.T) {

		releaseSpecificationOne := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("x-stevedore", "default", "", "", stevedore.ChartSpec{Name: "x-stevedore-dependencies", Dependencies: stevedore.Dependencies{{Name: "postgres-cluster", Repository: "http://some-chart-museum", Version: "5.0.6"}}}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)
		releaseSpecificationTwo := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("y-stevedore", "default", "", "", stevedore.ChartSpec{Name: "y-stevedore-dependencies", Dependencies: stevedore.Dependencies{{Name: "postgres-cluster", Repository: "http://some-chart-museum", Version: "5.0.6"}}}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "y-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)
		releaseSpecifications := stevedore.ReleaseSpecifications{
			releaseSpecificationOne,
			releaseSpecificationTwo,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockClient(ctrl)
		upstaller := mockUpstaller.NewMockUpstaller(ctrl)
		dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationOne).Return(releaseSpecificationOne, true, nil)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationTwo).Return(releaseSpecificationTwo, false, fmt.Errorf("build chart error occurred"))
		dependencyBuilder.EXPECT().UpdateRepo().Return(nil)

		file := "releaseSpecifications.yaml"

		opts := stevedore.Opts{DryRun: true, Parallel: true}
		upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), file, gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, _, responseCh, proceedCh, wg, _ interface{}, t int64) {
			defer wg.(*sync.WaitGroup).Done()

			proceedCh.(chan<- bool) <- true
			response := stevedore.Response{
				File:            file,
				ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			}
			responseCh.(chan<- stevedore.Response) <- response

		})

		expectedResponses := stevedore.Responses{
			{
				File:            file,
				ReleaseName:     releaseSpecificationOne.Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:                  file,
				ReleaseName:           releaseSpecificationTwo.Release.Name,
				ChartName:             releaseSpecificationTwo.Release.Chart,
				ChartVersion:          releaseSpecificationTwo.Release.ChartVersion,
				CurrentReleaseVersion: releaseSpecificationTwo.Release.CurrentReleaseVersion,
				Err:                   fmt.Errorf("failed to Build chart with error: build chart error occurred"),
			},
		}

		s := stevedore.Stevedore{
			Client:            client,
			Opts:              opts,
			Upstaller:         upstaller,
			DependencyBuilder: dependencyBuilder,
		}
		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
			Spec:     releaseSpecifications,
		}
		manifestFile := stevedore.ManifestFile{File: file, Manifest: manifest}
		responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{manifestFile}, timeout)

		assert.Len(t, responses, 2)
		assert.ElementsMatch(t, expectedResponses, responses)
		assert.Nil(t, err)
	})

	t.Run("should return error when repo update fails", func(t *testing.T) {

		releaseSpecificationOne := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("x-stevedore", "default", "chart/x-stevedore-dependencies", "", stevedore.ChartSpec{}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "x-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)
		releaseSpecificationTwo := stevedore.NewReleaseSpecification(
			stevedore.NewRelease("y-stevedore", "default", "", "", stevedore.ChartSpec{Name: "y-stevedore-dependencies", Dependencies: stevedore.Dependencies{{Name: "postgres-cluster", Repository: "http://some-chart-museum", Version: "5.0.6"}}}, 0, stevedore.Values{}, stevedore.Substitute{}, stevedore.Overrides{}),
			stevedore.Configs{
				"store": []map[string]interface{}{
					{"name": "y-stevedore", "tags": []string{"server"}},
					{"name": "store-ns"},
				},
			}, nil)
		releaseSpecifications := stevedore.ReleaseSpecifications{
			releaseSpecificationOne,
			releaseSpecificationTwo,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockClient(ctrl)
		upstaller := mockUpstaller.NewMockUpstaller(ctrl)
		dependencyBuilder := mockDependencyBuilder.NewMockDependencyBuilder(ctrl)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationOne).Return(releaseSpecificationOne, false, nil)
		dependencyBuilder.EXPECT().BuildChart(context.TODO(), releaseSpecificationTwo).Return(releaseSpecificationTwo, true, nil)
		dependencyBuilder.EXPECT().UpdateRepo().Return(fmt.Errorf("repo update failed"))

		file := "releaseSpecifications.yaml"
		opts := stevedore.Opts{DryRun: true, Parallel: true}
		upstaller.EXPECT().Upstall(context.TODO(), client, gomock.Any(), file, gomock.Any(), gomock.Any(), gomock.Any(), opts, timeout).Do(func(_, _, releaseSpecification, _, responseCh, proceedCh, wg, _ interface{}, t int64) {
			defer wg.(*sync.WaitGroup).Done()

			proceedCh.(chan<- bool) <- true
			response := stevedore.Response{
				File:            file,
				ReleaseName:     releaseSpecification.(stevedore.ReleaseSpecification).Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			}
			responseCh.(chan<- stevedore.Response) <- response

		})

		expectedResponses := stevedore.Responses{
			{
				File:            file,
				ReleaseName:     releaseSpecificationOne.Release.Name,
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:                  file,
				ReleaseName:           releaseSpecificationTwo.Release.Name,
				ChartName:             releaseSpecificationTwo.Release.Chart,
				ChartVersion:          releaseSpecificationTwo.Release.ChartVersion,
				CurrentReleaseVersion: releaseSpecificationTwo.Release.CurrentReleaseVersion,
				Err:                   fmt.Errorf("error updating helm repo with error: repo update failed"),
			},
		}

		s := stevedore.Stevedore{
			Client:            client,
			Opts:              opts,
			Upstaller:         upstaller,
			DependencyBuilder: dependencyBuilder,
		}
		manifest := stevedore.Manifest{
			DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
			Spec:     releaseSpecifications,
		}
		manifestFile := stevedore.ManifestFile{File: file, Manifest: manifest}
		responses, err := s.Do(context.TODO(), stevedore.ManifestFiles{manifestFile}, timeout)

		assert.Len(t, responses, 2)
		assert.ElementsMatch(t, expectedResponses, responses)
		assert.Nil(t, err)
	})

}
