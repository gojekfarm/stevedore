package helm

import (
	"fmt"

	"github.com/gojek/stevedore/internal/cli/command"
)

var binaryPath string

func binary() string {
	if binaryPath == "" {
		return "helm2.16"
	}
	return binaryPath
}

// Helm represents helm cli command
type Helm struct {
	command    string
	subCommand string
	args       []string
}

// Command returns helm command and its argument
func (helm Helm) Command() (string, []string) {
	args := []string{
		helm.command,
		helm.subCommand,
	}

	return binary(), append(args, helm.args...)
}

// Print the underlying command
func (helm Helm) Print() string {
	return command.Print(helm)
}

// Environ return the environment variable needed
//to be used when running command
func (helm Helm) Environ() []string {
	return []string{}
}

func newHelm(
	cmd,
	subCommand string,
	args []string,
) (command.Command, error) {
	helmCmd := Helm{
		command:    cmd,
		subCommand: subCommand,
		args:       args,
	}
	return command.NewCommand(helmCmd), nil
}

func globalFlags(tillerNamespace, kubeConfig, kubeContext string) []string {
	return []string{
		fmt.Sprintf("--tiller-namespace=%s", tillerNamespace),
		fmt.Sprintf("--kubeconfig=%s", kubeConfig),
		fmt.Sprintf("--kube-context=%s", kubeContext),
	}
}

// RepoList returns the helm repos
func RepoList() (command.Command, error) {
	return newHelm(
		"repo",
		"list",
		[]string{},
	)
}

// RepoAdd add the repo with name and URL
func RepoAdd(repo Repo) (command.Command, error) {
	return newHelm(
		"repo",
		"add",
		[]string{
			repo.Name,
			repo.URL,
		},
	)
}

// RepoUpdate update local helm repo cache
func RepoUpdate() (command.Command, error) {
	return newHelm(
		"repo",
		"update",
		[]string{},
	)
}

// List get the helm release for a give name tiller namespace
// releaseName is optional
func List(
	tillerNamespace,
	releaseName,
	kubeContext,
	kubeConfig string,
) (command.Command, error) {
	return newHelm(
		"ls",
		"",
		append(globalFlags(tillerNamespace, kubeConfig, kubeContext), releaseName),
	)
}

// Delete purges the helm release for the given tiller namespace
// and releaseName
func Delete(
	tillerNamespace,
	releaseName,
	kubeContext,
	kubeConfig string,
) (command.Command, error) {
	return newHelm(
		"delete",
		"",
		globalFlags(tillerNamespace, kubeConfig, kubeContext),
	)
}
