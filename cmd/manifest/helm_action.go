package manifest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/pkg/stevedore"
)

// HelmAction to apply manifest
type HelmAction struct {
	info         Info
	kubeconfig   string
	dryRun       bool
	parallel     bool
	filter       bool
	helmRepoName string
	helmTimeout  int64
}

type actionErrors []error

// NewHelmAction returns HelmAction with given arguments
func NewHelmAction(info Info, kubeconfig string, dryRun bool, parallel bool, filter bool, helmRepoName string, helmTimeout int64) HelmAction {
	return HelmAction{info, kubeconfig, dryRun, parallel, filter, helmRepoName, helmTimeout}
}

// Do will plan/apply manifests
func (action HelmAction) Do() (Info, error) {
	manifestFiles := action.info.ManifestFiles
	opts := stevedore.Opts{DryRun: action.dryRun, Parallel: action.parallel, Filter: action.filter}

	responses, err := stevedore.CreateResponse(context.TODO(), manifestFiles, opts, action.helmRepoName, action.helmTimeout)

	if err != nil {
		return Info{}, err
	}
	if len(responses) == 0 {
		cli.Warn("No changes in the plan")
		return Info{}, nil
	}

	group := responses.GroupByFile()
	responseHandler := newResponseHandler(action.dryRun)
	if err := responseHandler.render(group); err != nil {
		return Info{}, err
	}

	summarizer := TableSummarizer{writer: cli.OutputStream()}
	summarizer.Display(group)

	filteredInfos := action.info.FilterBy(responses)
	errors := getAllErrors(responses)

	if len(errors) == 0 {
		return filteredInfos, nil
	}

	return filteredInfos, errors
}

func getAllErrors(responses stevedore.Responses) actionErrors {
	var errors actionErrors
	for _, response := range responses {
		if response.Err != nil {
			errors = append(errors, response.Err)
		}
	}
	return errors
}

func (errors actionErrors) Error() string {
	buff := bytes.NewBufferString(fmt.Sprintf("Failed with %d Errors\n", len(errors)))
	for i, err := range errors {
		buff.WriteString(fmt.Sprintf("\n%d. %s\n", i+1, err.Error()))
	}
	return buff.String()
}
