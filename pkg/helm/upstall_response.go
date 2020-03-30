package helm

import (
	"strings"

	"github.com/databus23/helm-diff/manifest"
)

// UpstallResponse to encapsulate all release related details
type UpstallResponse struct {
	ExistingSpecs         map[string]*manifest.MappingResult
	NewSpecs              map[string]*manifest.MappingResult
	HasDiff               bool
	Diff                  string
	ChartVersion          string
	CurrentReleaseVersion int32
}

// Summary generates Summary for helm release
func (up UpstallResponse) Summary() Summary {
	summary := Summary{Added: Resources{}, Modified: Resources{}, Destroyed: Resources{}}
	for key, spec := range up.NewSpecs {
		resourceName := strings.TrimSpace(strings.Split(spec.Name, ",")[1])
		if _, ok := up.ExistingSpecs[key]; !ok {
			summary.Added = append(summary.Added, Resource{Name: resourceName, Kind: spec.Kind})
		} else {
			if ok := spec.Content == up.ExistingSpecs[key].Content; !ok {
				summary.Modified = append(summary.Modified, Resource{Name: resourceName, Kind: spec.Kind})
			}
		}
	}
	for key, spec := range up.ExistingSpecs {
		resourceName := strings.TrimSpace(strings.Split(spec.Name, ",")[1])
		if _, ok := up.NewSpecs[key]; !ok {
			summary.Destroyed = append(summary.Destroyed, Resource{Name: resourceName, Kind: spec.Kind})
		}
	}
	return summary
}
