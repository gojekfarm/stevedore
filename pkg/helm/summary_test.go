package helm_test

import (
	"github.com/gojek/stevedore/pkg/helm"
	"reflect"
	"testing"
)

func TestResourcesGroupByKind(t *testing.T) {
	t.Run("should group resources by kind", func(t *testing.T) {
		resources := helm.Resources{
			{Name: "app-db-exporter", Kind: "ConfigMap"},
			{Name: "app-db", Kind: "ConfigMap"},
			{Name: "app-server", Kind: "Deployment"},
			{Name: "app-cm", Kind: "ConfigMap"},
			{Name: "app-db-keeper", Kind: "StatefulSet"},
			{Name: "app-db-exporter", Kind: "Deployment"},
		}

		expected := map[string]helm.Resources{
			"ConfigMap": {
				{Name: "app-db-exporter", Kind: "ConfigMap"},
				{Name: "app-db", Kind: "ConfigMap"},
				{Name: "app-cm", Kind: "ConfigMap"},
			},
			"Deployment": {
				{Name: "app-server", Kind: "Deployment"},
				{Name: "app-db-exporter", Kind: "Deployment"},
			},
			"StatefulSet": {
				{Name: "app-db-keeper", Kind: "StatefulSet"},
			},
		}

		actual := resources.GroupByKind()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", actual, expected)
		}
	})
}
