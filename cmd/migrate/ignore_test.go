package migrate

import (
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIgnoreStrategy_Do(t *testing.T) {
	t.Run("should migrate to newer v1Stevedore ignore file", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		ignoreStrategy := NewIgnoreStrategy(memFs, []string{"ignore/a.yaml", "ignore/b.yaml"})

		aFileContent := `
- matches:
    contextName: one
  components:
    - name: x-service
      reason: ignoring as the k8s resources already exist
- matches:
    contextName: two
  components:
    - name: y-service
      reason: ignoring in second environment
`
		bFileContent := `
- matches:
    contextName: three
  components:
    - name: y-service
      reason: ignoring in third environment
`
		_ = afero.WriteFile(memFs, "ignore/a.yaml", []byte(aFileContent), 0666)
		_ = afero.WriteFile(memFs, "ignore/b.yaml", []byte(bFileContent), 0666)
		expectedA := stevedore.Ignores{
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionContextName: "one",
				},
				Releases: stevedore.IgnoredReleases{
					{
						Name:   "x-service",
						Reason: "ignoring as the k8s resources already exist",
					},
				},
			},
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionContextName: "two",
				},
				Releases: stevedore.IgnoredReleases{
					{
						Name:   "y-service",
						Reason: "ignoring in second environment",
					},
				},
			},
		}
		expectedB := stevedore.Ignores{
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionContextName: "three",
				},
				Releases: stevedore.IgnoredReleases{
					{
						Name:   "y-service",
						Reason: "ignoring in third environment",
					},
				},
			},
		}

		err := ignoreStrategy.Do()
		actualA := stevedore.Ignores{}
		actualB := stevedore.Ignores{}

		_ = read(memFs, "ignore/a.yaml", &actualA)
		_ = read(memFs, "ignore/b.yaml", &actualB)

		assert.NoError(t, err)
		if !cmp.Equal(expectedA, actualA) {
			assert.Fail(t, cmp.Diff(expectedA, actualA))
		}

		if !cmp.Equal(expectedB, actualB) {
			assert.Fail(t, cmp.Diff(expectedB, actualB))
		}
	})
}
