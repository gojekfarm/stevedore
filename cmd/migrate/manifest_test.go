package migrate

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManifestStrategyDo(t *testing.T) {
	t.Run("should convert older stevedore manifest file into newer one", func(t *testing.T) {
		manifest := `
deployTo:
  - one
  - two
applications:
  - component:
      name: example-postgres
      # chart: official/postgresql
      chartSpec:
        name: example-postgres-dependencies
        dependencies:
          - name: postgresql
            alias: db
            repository: https://kubernetes-charts.storage.googleapis.com
            version: 3.14.2
          - name: redis
            repository: https://kubernetes-charts.storage.googleapis.com
            version: 6.4.2
      namespace: default
      privileged: true
      values:
        postgresqlPassword: myPostgresPassword
        image:
          pullPolicy: IfNotPresent
    configs:
      consul:
      - name: app
`
		memFs := afero.NewMemMapFs()
		manifestFilePath := "old/sample.yaml"
		_ = memFs.Mkdir("old", 0666)
		_ = afero.WriteFile(memFs, manifestFilePath, []byte(manifest), 0666)
		expected := stevedore.Manifest{
			Kind:    stevedore.KindStevedoreManifest,
			Version: stevedore.ManifestCurrentVersion,
			DeployTo: stevedore.Matchers{
				{stevedore.ConditionContextName: "one"},
				{stevedore.ConditionContextName: "two"},
			},
			Spec: stevedore.ReleaseSpecifications{
				{
					Release: stevedore.Release{
						Name:      "example-postgres",
						Namespace: "default",
						ChartSpec: stevedore.ChartSpec{
							Name: "example-postgres-dependencies",
							Dependencies: stevedore.Dependencies{
								{
									Name:       "postgresql",
									Alias:      "db",
									Version:    "3.14.2",
									Repository: "https://kubernetes-charts.storage.googleapis.com",
								},
								{
									Name:       "redis",
									Version:    "6.4.2",
									Repository: "https://kubernetes-charts.storage.googleapis.com",
								},
							},
						},
						Values: stevedore.Values{
							"postgresqlPassword": "myPostgresPassword",
							"image": map[interface{}]interface{}{
								"pullPolicy": "IfNotPresent",
							},
						},
						Privileged: true,
					},
					Configs: stevedore.Configs{"consul": []interface{}{map[interface{}]interface{}{"name": "app"}}},
				},
			},
		}

		migrateStrategy := NewManifestStrategy(memFs, manifestFilePath, stevedore.Contexts{}, true)

		err := migrateStrategy.Do()
		assert.NoError(t, err)

		actual := stevedore.Manifest{}
		err = read(memFs, manifestFilePath, &actual)
		assert.NoError(t, err)

		ignoreTypes := cmpopts.IgnoreUnexported(stevedore.Release{}, stevedore.ReleaseSpecification{})
		if !cmp.Equal(expected, actual, ignoreTypes) {
			assert.Fail(t, cmp.Diff(expected, actual, ignoreTypes))
		}
	})
}
