package cmd

import (
	"fmt"
	"io"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFmtFormat(t *testing.T) {
	t.Run("should format manifest", func(t *testing.T) {
		manifestFilePath := "/tmp/manifest.yaml"
		fs := afero.NewMemMapFs()
		manifestString := `kind: StevedoreManifest
version: "2"
deployTo:
- contextName: cluster-1
- contextName: cluster-2
spec:
- release:
    name: test-dobby
    namespace: default
    chart: helm-repo/dobby-dependencies
    values:
      db:
        persistence:
          size: ${DB_SIZE}
        metrics:
          enabled: false
        store:
          backend:
            type: consul
            endpoint: http://consul:8500
  configs:
    consul:
      path: dobby`

		_ = afero.WriteFile(fs, manifestFilePath, []byte(manifestString), 0644)

		contents, err := getContents(fs, manifestFilePath, func(reader io.Reader) (formatter fmt.Formatter, err error) {
			return stevedore.NewManifest(reader)
		})
		assert.NoError(t, err)

		err = format(fs, contents...)
		assert.NoError(t, err)

		actual, err := afero.ReadFile(fs, manifestFilePath)
		assert.NoError(t, err)

		expected := `kind: StevedoreManifest
version: "2"
deployTo:
- contextName: cluster-1
- contextName: cluster-2
spec:
- release:
    name: test-dobby
    namespace: default
    chart: helm-repo/dobby-dependencies
    values:
      db:
        metrics:
          enabled: false
        persistence:
          size: ${DB_SIZE}
        store:
          backend:
            endpoint: http://consul:8500
            type: consul
  configs:
    consul:
      path: dobby
`

		assert.Equal(t, expected, string(actual))
	})

	t.Run("should format overrides", func(t *testing.T) {
		overridesPath := "/tmp/overrides.yaml"
		fs := afero.NewMemMapFs()
		overridesString := `
kind: StevedoreOverride
version: "2"
spec:
- matches:
    environmentType: staging
    contextType: components
  values:
    redis:
      persistence:
        size: 5Gi
`

		_ = afero.WriteFile(fs, overridesPath, []byte(overridesString), 0644)

		contents, err := getContents(fs, overridesPath, func(reader io.Reader) (formatter fmt.Formatter, err error) {
			return stevedore.NewOverrides(reader)
		})
		assert.NoError(t, err)

		err = format(fs, contents...)
		assert.NoError(t, err)

		actual, err := afero.ReadFile(fs, overridesPath)
		assert.NoError(t, err)

		expected := `kind: StevedoreOverride
version: "2"
spec:
- matches:
    contextType: components
    environmentType: staging
  values:
    redis:
      persistence:
        size: 5Gi
`
		assert.Equal(t, expected, string(actual))
	})

	t.Run("should format env", func(t *testing.T) {
		envFilePath := "/tmp/env.yaml"
		fs := afero.NewMemMapFs()
		manifestString := `
kind: StevedoreEnv
version: "2"
spec:
- matches:
    contextName: cluster-1
  env:
    RABBIT_MQ_PASSWORD: 456
    INGRESS_HOST: some.host
    AVAILABILITY: false
- matches:
    contextName: cluster-1
  env:
    DBNAME: admin
    PASSWORD: admin
    HOST: example.com
`
		_ = afero.WriteFile(fs, envFilePath, []byte(manifestString), 0644)

		contents, err := getContents(fs, envFilePath, func(reader io.Reader) (formatter fmt.Formatter, err error) {
			return stevedore.NewEnv(reader)
		})
		assert.NoError(t, err)

		err = format(fs, contents...)
		assert.NoError(t, err)

		actual, err := afero.ReadFile(fs, envFilePath)
		assert.NoError(t, err)

		expected := `kind: StevedoreEnv
version: "2"
spec:
- matches:
    contextName: cluster-1
  env:
    AVAILABILITY: false
    INGRESS_HOST: some.host
    RABBIT_MQ_PASSWORD: 456
- matches:
    contextName: cluster-1
  env:
    DBNAME: admin
    HOST: example.com
    PASSWORD: admin
`

		assert.Equal(t, expected, string(actual))
	})
}
