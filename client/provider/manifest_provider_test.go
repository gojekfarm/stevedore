package provider_test

import (
	"github.com/gojek/stevedore/client/provider"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDefaultManifestProviderManifests(t *testing.T) {
	t.Run("recursive:false", func(t *testing.T) {
		t.Run("should return all the manifests in the dir", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			manifestFileName := "/mock/services/service.yaml"
			manifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: services
spec:
  - release:
      name: example-postgres
      chartSpec:
        name: x-stevedore-dependencies
        dependencies:
          - name: redis
            alias: cache-redis
            version: 1.9.4
            repository: 'http://helm-charts/'
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always
    mounts:
      store:
        - name: dobby
          tags: ["server"]
          useAs: server.config
`
			anotherManifestFileName := "/mock/services/another-service.yaml"
			anotherManifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: services
spec:
  - release:
      name: example-redis
      chart: stable/redis
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always`

			_ = afero.WriteFile(memFs, manifestFileName, []byte(manifestString), 0644)
			_ = afero.WriteFile(memFs, anotherManifestFileName, []byte(anotherManifestString), 0644)

			manifestProvider := provider.NewManifestProvider(memFs)
			expected := stevedore.ManifestFiles{
				stevedore.ManifestFile{
					File: "/mock/services/another-service.yaml",
					Manifest: stevedore.Manifest{
						Kind:     stevedore.KindStevedoreManifest,
						Version:  stevedore.ManifestCurrentVersion,
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.ReleaseSpecification{
								Release: stevedore.Release{
									Name:      "example-redis",
									Chart:     "stable/redis",
									Namespace: "default",
									Values: stevedore.Values{
										"user":  "${user}",
										"image": map[interface{}]interface{}{"pullPolicy": "Always"},
									},
								},
							},
						},
					},
				},
				stevedore.ManifestFile{
					File: "/mock/services/service.yaml",
					Manifest: stevedore.Manifest{
						Kind:     stevedore.KindStevedoreManifest,
						Version:  stevedore.ManifestCurrentVersion,
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.ReleaseSpecification{
								Release: stevedore.Release{
									Name: "example-postgres",
									ChartSpec: stevedore.ChartSpec{
										Name: "x-stevedore-dependencies",
										Dependencies: stevedore.Dependencies{
											stevedore.Dependency{
												Name:       "redis",
												Alias:      "cache-redis",
												Version:    "1.9.4",
												Repository: "http://helm-charts/"},
										},
									}, Namespace: "default",
									Values: stevedore.Values{
										"user":  "${user}",
										"image": map[interface{}]interface{}{"pullPolicy": "Always"},
									},
								},
								Mounts: stevedore.Configs{
									"store": []interface{}{
										map[interface{}]interface{}{
											"name": "dobby", "tags": []interface{}{"server"}, "useAs": "server.config",
										},
									},
								},
							},
						},
					},
				},
			}

			context := map[string]string{
				provider.ManifestPathFile:      "/mock/services",
				provider.ManifestRecursiveFlag: "false",
				provider.EnvironmentTypeKey:    "dev",
			}
			manifests, err := manifestProvider.Manifests(context)

			assert.NoError(t, err)
			assert.Equal(t, expected, manifests)
		})

		t.Run("should return the specified manifest file from the dir", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			manifestFileName := "/mock/services/service.yaml"
			manifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: services
spec:
  - release:
      name: example-postgres
      chart: stable/postgresql
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always`

			anotherManifestFileName := "/mock/services/another-service.yaml"
			anotherManifestString := `
deployTo:
  - contextName: services
spec:
  - release:
      name: example-redis
      chart: stable/redis
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always`

			_ = afero.WriteFile(memFs, manifestFileName, []byte(manifestString), 0644)
			_ = afero.WriteFile(memFs, anotherManifestFileName, []byte(anotherManifestString), 0644)

			manifestProvider := provider.NewManifestProvider(memFs)
			expected := stevedore.ManifestFiles{
				stevedore.ManifestFile{
					File: "/mock/services/service.yaml",
					Manifest: stevedore.Manifest{
						Kind:     stevedore.KindStevedoreManifest,
						Version:  stevedore.ManifestCurrentVersion,
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.ReleaseSpecification{
								Release: stevedore.Release{
									Name:      "example-postgres",
									Chart:     "stable/postgresql",
									Namespace: "default",
									Values: stevedore.Values{
										"user":  "${user}",
										"image": map[interface{}]interface{}{"pullPolicy": "Always"},
									},
								},
							},
						},
					},
				},
			}

			context := map[string]string{
				provider.ManifestPathFile:      "/mock/services/service.yaml",
				provider.ManifestRecursiveFlag: "false",
				provider.EnvironmentTypeKey:    "dev",
			}
			manifests, err := manifestProvider.Manifests(context)

			assert.NoError(t, err)
			assert.Equal(t, expected, manifests)
		})
	})

	t.Run("recursive:true", func(t *testing.T) {
		t.Run("should not respect recursive flag and return all the manifests in the dir", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			manifestFileName := "/mock/services/service.yaml"
			manifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: services
spec:
  - release:
      name: example-postgres
      chartSpec:
        name: x-stevedore-dependencies
        dependencies:
          - name: redis
            alias: cache-redis
            version: 1.9.4
            repository: 'http://helm-charts/'
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always`

			anotherManifestFileName := "/mock/services/another-service.yaml"
			anotherManifestString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: services
spec:
  - release:
      name: example-redis
      chart: stable/redis
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always`

			_ = afero.WriteFile(memFs, manifestFileName, []byte(manifestString), 0644)
			_ = afero.WriteFile(memFs, anotherManifestFileName, []byte(anotherManifestString), 0644)

			manifestProvider := provider.NewManifestProvider(memFs)
			expected := stevedore.ManifestFiles{
				stevedore.ManifestFile{
					File: "/mock/services/another-service.yaml",
					Manifest: stevedore.Manifest{
						Kind:     stevedore.KindStevedoreManifest,
						Version:  stevedore.ManifestCurrentVersion,
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.ReleaseSpecification{
								Release: stevedore.Release{
									Name:      "example-redis",
									Chart:     "stable/redis",
									Namespace: "default",
									Values: stevedore.Values{
										"user":  "${user}",
										"image": map[interface{}]interface{}{"pullPolicy": "Always"},
									},
								},
							},
						},
					},
				},
				stevedore.ManifestFile{
					File: "/mock/services/service.yaml",
					Manifest: stevedore.Manifest{
						Kind:     stevedore.KindStevedoreManifest,
						Version:  stevedore.ManifestCurrentVersion,
						DeployTo: stevedore.Matchers{{stevedore.ConditionContextName: "services"}},
						Spec: stevedore.ReleaseSpecifications{
							stevedore.ReleaseSpecification{
								Release: stevedore.Release{
									Name: "example-postgres",
									ChartSpec: stevedore.ChartSpec{
										Name: "x-stevedore-dependencies",
										Dependencies: stevedore.Dependencies{
											stevedore.Dependency{
												Name:       "redis",
												Alias:      "cache-redis",
												Version:    "1.9.4",
												Repository: "http://helm-charts/"},
										},
									}, Namespace: "default",
									Values: stevedore.Values{
										"user":  "${user}",
										"image": map[interface{}]interface{}{"pullPolicy": "Always"},
									},
								},
							},
						},
					},
				},
			}

			context := map[string]string{
				provider.ManifestPathFile:      "/mock/services",
				provider.ManifestRecursiveFlag: "true",
				provider.EnvironmentTypeKey:    "dev",
			}
			manifests, err := manifestProvider.Manifests(context)

			assert.NoError(t, err)
			assert.Equal(t, expected, manifests)
		})

		t.Run("should return all dependent manifest files specified in the provided manifest file from the dir", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			serviceOneFileName := "/mock/services/service-one.yaml"
			serviceOneString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: services
spec:
  - release:
      name: service-one
      chart: stable/postgresql
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always
    dependsOn:
      - service-two
`

			serviceTwoFileName := "/mock/services/service-two.yaml"
			serviceTwoString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: components
spec:
  - release:
      name: service-two
      chart: stable/redis
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always
    dependsOn:
      - service-four
`

			serviceThreeFileName := "/mock/services/service-three.yaml"
			serviceThreeString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: production
spec:
  - release:
      name: service-three
      chart: stable/rabbitmq
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always`

			serviceFourFileName := "/mock/services/service-four.yaml"
			serviceFourString := `
kind: StevedoreManifest
version: "2"
deployTo:
  - contextName: readonly
spec:
  - release:
      name: service-four
      chart: stable/mysql
      namespace: default
      values:
        user: ${user}
        image:
          pullPolicy: Always
    dependsOn:
      - service-one
`
			_ = afero.WriteFile(memFs, serviceOneFileName, []byte(serviceOneString), 0644)
			_ = afero.WriteFile(memFs, serviceTwoFileName, []byte(serviceTwoString), 0644)
			_ = afero.WriteFile(memFs, serviceThreeFileName, []byte(serviceThreeString), 0644)
			_ = afero.WriteFile(memFs, serviceFourFileName, []byte(serviceFourString), 0644)

			manifestProvider := provider.NewManifestProvider(memFs)

			expected := map[string]stevedore.ReleaseSpecification{
				"service-one": {
					Release: stevedore.Release{
						Name:      "service-one",
						Chart:     "stable/postgresql",
						Namespace: "default",
						Values: stevedore.Values{
							"user":  "${user}",
							"image": map[interface{}]interface{}{"pullPolicy": "Always"},
						},
					},
					DependsOn: []string{"service-two"},
				},
				"service-two": {
					Release: stevedore.Release{
						Name:      "service-two",
						Chart:     "stable/redis",
						Namespace: "default",
						Values: stevedore.Values{
							"user":  "${user}",
							"image": map[interface{}]interface{}{"pullPolicy": "Always"},
						},
					},
					DependsOn: []string{"service-four"},
				},
				"service-four": {
					Release: stevedore.Release{
						Name:      "service-four",
						Chart:     "stable/mysql",
						Namespace: "default",
						Values: stevedore.Values{
							"user":  "${user}",
							"image": map[interface{}]interface{}{"pullPolicy": "Always"},
						},
					},
					DependsOn: []string{"service-one"},
				},
			}

			context := map[string]string{
				provider.ManifestPathFile:      "/mock/services/service-one.yaml",
				provider.ManifestRecursiveFlag: "true",
				provider.EnvironmentTypeKey:    "dev",
			}
			manifests, err := manifestProvider.Manifests(context)

			actual := map[string]stevedore.ReleaseSpecification{}
			for _, manifest := range manifests {
				for _, releaseSpec := range manifest.Spec {
					actual[releaseSpec.Release.Name] = releaseSpec
				}
			}

			assert.NoError(t, err)
			assert.Equal(t, "/mock/services/service-one.yaml", manifests[0].File)
			assert.Equal(t, stevedore.Matchers{{stevedore.ConditionContextName: "services"}}, manifests[0].DeployTo)

			ignoreTypes := cmpopts.IgnoreTypes(stevedore.Substitute{}, stevedore.Overrides{})
			if !cmp.Equal(expected, actual, ignoreTypes) {
				assert.Fail(t, cmp.Diff(expected, actual, ignoreTypes))
			}
		})
	})

	t.Run("it should return error if the manifest files doesn't exists", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		_ = memFs.Mkdir("/mock/services", 0644)

		manifestProvider := provider.NewManifestProvider(memFs)

		context := map[string]string{
			provider.ManifestPathFile:      "/mock/services/service.yaml",
			provider.ManifestRecursiveFlag: "false",
		}
		manifestFiles, err := manifestProvider.Manifests(context)

		if assert.Error(t, err) {
			assert.Equal(t, "invalid file path. Provide a valid path to stevedore manifests using --manifests-path", err.Error())
		}
		assert.Nil(t, manifestFiles)
	})
}
