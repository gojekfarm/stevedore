package stevedore

import (
	"sort"

	"github.com/gojek/stevedore/pkg/helm"
)

// Response represents a stevedore release response
type Response struct {
	File                  string
	ReleaseName           string
	ChartName             string
	ChartVersion          string
	CurrentReleaseVersion int32
	helm.UpstallResponse
	Err error
}

// Responses is a collection of stevedore release response
type Responses []Response

// SortByReleaseName sorts the responses by manifest name
func (responses Responses) SortByReleaseName() {
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].ReleaseName < responses[j].ReleaseName
	})
}

// GroupedResponses is group collection of stevedore release responses
type GroupedResponses map[string]Responses

// GroupByFile groups stevedore release responses by file name
func (responses Responses) GroupByFile() GroupedResponses {
	accumulator := make(GroupedResponses)
	for _, response := range responses {
		file := response.File
		if responses, ok := accumulator[file]; ok {
			responses = append(responses, response)
			accumulator[file] = responses
		} else {
			accumulator[file] = Responses{response}
		}
	}
	return accumulator
}

// GetReleaseNames returns release names
func (responses Responses) GetReleaseNames() []string {
	var releaseNames []string
	for _, response := range responses {
		releaseNames = append(releaseNames, response.ReleaseName)
	}
	return releaseNames
}

// Find returns release names
func (responses Responses) Find(releaseName string) Response {
	for _, response := range responses {
		if releaseName == response.ReleaseName {
			return response
		}
	}
	return Response{}
}

// SortedFileNames return list of file names in sorted order
func (g GroupedResponses) SortedFileNames() []string {
	var keys []string
	for k := range g {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
