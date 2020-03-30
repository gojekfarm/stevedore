package manifest

import (
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/pkg/stevedore"
)

func newResponseHandler(dryRun bool) responseHandler {
	if dryRun {
		return planResponseHandler{}
	}
	return applyResponseHandler{}
}

func iterator(group stevedore.GroupedResponses, renderer func(response stevedore.Response) error) error {
	sortedFileNames := group.SortedFileNames()
	for _, fileName := range sortedFileNames {
		responses := group[fileName]
		cli.Infof("File: %s", fileName)
		for _, response := range responses {
			if err := renderer(response); err != nil {
				return err
			}
		}
	}
	return nil
}

type responseHandler interface {
	render(group stevedore.GroupedResponses) error
}

type planResponseHandler struct{}

func (p planResponseHandler) render(group stevedore.GroupedResponses) error {
	return iterator(group, p.renderer)
}

func (p planResponseHandler) renderer(response stevedore.Response) error {
	cli.Infof("ReleaseName: %s", response.ReleaseName)
	if response.Err != nil {
		cli.Errorf("Error: %v", response.Err)
	} else if response.HasDiff {
		cli.Info(response.Diff)
	} else {
		cli.Info("No diff present\n")
	}
	return nil
}

type applyResponseHandler struct{}

func (a applyResponseHandler) render(group stevedore.GroupedResponses) error {
	return iterator(group, a.renderer)
}

func (a applyResponseHandler) renderer(response stevedore.Response) error {
	cli.Infof("ReleaseName: %s", response.ReleaseName)
	if response.Err != nil {
		cli.Errorf("Error: %v", response.Err)
		return response.Err
	}
	cli.Info("Applied Successfully\n")
	return nil
}
