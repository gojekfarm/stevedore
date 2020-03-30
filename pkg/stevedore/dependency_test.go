package stevedore

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/helm/pkg/chartutil"
	"testing"
)

func TestNewDependencies(t *testing.T) {
	t.Run("should convert chartutil.Dependencies to stevedore.Dependencies", func(t *testing.T) {
		chartDependencies := []*chartutil.Dependency{
			{
				Name:         "name",
				Version:      "0.0.1",
				Alias:        "alias",
				Repository:   "https://repo",
				Condition:    "condition",
				Enabled:      true,
				ImportValues: []interface{}{"one", "two"},
				Tags:         []string{"tag", "alias"},
			},
		}
		expected := Dependencies{
			{
				Name:         "name",
				Version:      "0.0.1",
				Alias:        "alias",
				Repository:   "https://repo",
				Condition:    "condition",
				Enabled:      true,
				ImportValues: []interface{}{"one", "two"},
				Tags:         []string{"tag", "alias"},
			},
		}

		actual := NewDependencies(chartDependencies)

		assert.Equal(t, expected, actual)
	})

	t.Run("should not fail when converting nil", func(t *testing.T) {
		expected := Dependencies{}

		actual := NewDependencies(nil)

		assert.Equal(t, expected, actual)
	})
}

func TestChartUtilDependency(t *testing.T) {
	t.Run("should convert stevedore.Dependencies to chartutil.Dependencies", func(t *testing.T) {
		dependency := Dependency{
			Name:         "name",
			Version:      "0.0.1",
			Alias:        "alias",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}

		expected := chartutil.Dependency{
			Name:         "name",
			Version:      "0.0.1",
			Alias:        "alias",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}

		actual := dependency.ChartUtilDependency()

		assert.Equal(t, expected, actual)
	})

	t.Run("should not fail when converting nil", func(t *testing.T) {
		expected := Dependencies{}

		actual := NewDependencies(nil)

		assert.Equal(t, expected, actual)
	})
}

func TestDependenciesCheckSum(t *testing.T) {
	t.Run("should compute 8 character shasum", func(t *testing.T) {
		dependencyA := Dependency{
			Name:         "dependency-a",
			Version:      "0.0.1",
			Alias:        "dependency-a",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}
		dependencyB := Dependency{
			Name:         "dependency-b",
			Version:      "0.0.1",
			Alias:        "dependency-b",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}

		dependencies := Dependencies{dependencyA, dependencyB}

		sum, err := dependencies.CheckSum()

		assert.NoError(t, err)
		assert.Equal(t, "352a45d2", sum)
	})

	t.Run("should return same shasum even if order of dependencies change", func(t *testing.T) {
		dependencyA := Dependency{
			Name:         "dependency-a",
			Version:      "0.0.1",
			Alias:        "dependency-a",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}
		dependencyB := Dependency{
			Name:         "dependency-b",
			Version:      "0.0.1",
			Alias:        "dependency-b",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}

		dependencies := Dependencies{dependencyA, dependencyB}
		anotherDependencies := Dependencies{dependencyB, dependencyA}

		sum, err := dependencies.CheckSum()
		assert.NoError(t, err)

		anotherSum, err := anotherDependencies.CheckSum()

		assert.NoError(t, err)
		assert.Equal(t, sum, anotherSum)
	})

	t.Run("should use name if alias not mentioned", func(t *testing.T) {
		dependencyA := Dependency{
			Name:         "redis",
			Version:      "0.0.1",
			Alias:        "redisA",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}
		dependencyB := Dependency{
			Name:         "redis",
			Version:      "0.0.1",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}
		dependencyC := Dependency{
			Name:         "redis",
			Version:      "0.0.1",
			Alias:        "redisC",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}
		dependencyD := Dependency{
			Name:         "postgres",
			Version:      "0.0.1",
			Repository:   "https://repo",
			Condition:    "condition",
			Enabled:      true,
			ImportValues: []interface{}{"one", "two"},
			Tags:         []string{"tag", "alias"},
		}

		combination1 := Dependencies{dependencyA, dependencyB, dependencyC, dependencyD}
		combination2 := Dependencies{dependencyA, dependencyB, dependencyC, dependencyD}
		combination3 := Dependencies{dependencyB, dependencyA, dependencyC, dependencyD}
		combination4 := Dependencies{dependencyB, dependencyC, dependencyA, dependencyD}
		combination5 := Dependencies{dependencyC, dependencyA, dependencyB, dependencyD}
		combination6 := Dependencies{dependencyC, dependencyB, dependencyA, dependencyD}

		combination7 := Dependencies{dependencyA, dependencyB, dependencyD, dependencyC}
		combination8 := Dependencies{dependencyA, dependencyB, dependencyD, dependencyC}
		combination9 := Dependencies{dependencyB, dependencyA, dependencyD, dependencyC}
		combination10 := Dependencies{dependencyB, dependencyC, dependencyD, dependencyA}
		combination11 := Dependencies{dependencyC, dependencyA, dependencyD, dependencyB}
		combination12 := Dependencies{dependencyC, dependencyB, dependencyD, dependencyA}

		combination13 := Dependencies{dependencyA, dependencyD, dependencyB, dependencyC}
		combination14 := Dependencies{dependencyA, dependencyD, dependencyB, dependencyC}
		combination15 := Dependencies{dependencyB, dependencyD, dependencyA, dependencyC}
		combination16 := Dependencies{dependencyB, dependencyD, dependencyC, dependencyA}
		combination17 := Dependencies{dependencyC, dependencyD, dependencyA, dependencyB}
		combination18 := Dependencies{dependencyC, dependencyD, dependencyB, dependencyA}

		combination19 := Dependencies{dependencyD, dependencyA, dependencyB, dependencyC}
		combination20 := Dependencies{dependencyD, dependencyA, dependencyB, dependencyC}
		combination21 := Dependencies{dependencyD, dependencyB, dependencyA, dependencyC}
		combination22 := Dependencies{dependencyD, dependencyB, dependencyC, dependencyA}
		combination23 := Dependencies{dependencyD, dependencyC, dependencyA, dependencyB}
		combination24 := Dependencies{dependencyD, dependencyC, dependencyB, dependencyA}

		combination1sum, _ := combination1.CheckSum()
		combination2sum, _ := combination2.CheckSum()
		combination3sum, _ := combination3.CheckSum()
		combination4sum, _ := combination4.CheckSum()
		combination5sum, _ := combination5.CheckSum()
		combination6sum, _ := combination6.CheckSum()

		combination7sum, _ := combination7.CheckSum()
		combination8sum, _ := combination8.CheckSum()
		combination9sum, _ := combination9.CheckSum()
		combination10sum, _ := combination10.CheckSum()
		combination11sum, _ := combination11.CheckSum()
		combination12sum, _ := combination12.CheckSum()

		combination13sum, _ := combination13.CheckSum()
		combination14sum, _ := combination14.CheckSum()
		combination15sum, _ := combination15.CheckSum()
		combination16sum, _ := combination16.CheckSum()
		combination17sum, _ := combination17.CheckSum()
		combination18sum, _ := combination18.CheckSum()

		combination19sum, _ := combination19.CheckSum()
		combination20sum, _ := combination20.CheckSum()
		combination21sum, _ := combination21.CheckSum()
		combination22sum, _ := combination22.CheckSum()
		combination23sum, _ := combination23.CheckSum()
		combination24sum, _ := combination24.CheckSum()

		assert.Equal(t, combination1sum, combination2sum)
		assert.Equal(t, combination1sum, combination3sum)
		assert.Equal(t, combination1sum, combination4sum)
		assert.Equal(t, combination1sum, combination5sum)
		assert.Equal(t, combination1sum, combination6sum)
		assert.Equal(t, combination1sum, combination7sum)
		assert.Equal(t, combination1sum, combination8sum)
		assert.Equal(t, combination1sum, combination9sum)
		assert.Equal(t, combination1sum, combination10sum)
		assert.Equal(t, combination1sum, combination11sum)
		assert.Equal(t, combination1sum, combination12sum)

		assert.Equal(t, combination1sum, combination13sum)
		assert.Equal(t, combination1sum, combination14sum)
		assert.Equal(t, combination1sum, combination15sum)
		assert.Equal(t, combination1sum, combination16sum)
		assert.Equal(t, combination1sum, combination17sum)
		assert.Equal(t, combination1sum, combination18sum)
		assert.Equal(t, combination1sum, combination19sum)
		assert.Equal(t, combination1sum, combination20sum)
		assert.Equal(t, combination1sum, combination21sum)
		assert.Equal(t, combination1sum, combination22sum)
		assert.Equal(t, combination1sum, combination23sum)
		assert.Equal(t, combination1sum, combination24sum)
	})

}

func TestDependenciesContains(t *testing.T) {
	dependencies := Dependencies{
		{Name: "dep-1", Alias: "alias-1"},
		{Name: "dep-1", Alias: "alias-2"},
		{Name: "dep-2", Alias: "alias-3"},
	}
	t.Run("should return matching dependencies and true if releaseSpecification contains given chart name as dependency", func(t *testing.T) {
		expected := Dependencies{
			{Name: "dep-1", Alias: "alias-1"},
			{Name: "dep-1", Alias: "alias-2"},
		}

		actual, ok := dependencies.Contains("dep-1")

		assert.True(t, ok)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return false if releaseSpecification does not contain given chart name as dependency", func(t *testing.T) {
		expected := Dependencies{}

		actual, ok := dependencies.Contains("dep-3")

		assert.False(t, ok)
		assert.Equal(t, expected, actual)
	})
}
