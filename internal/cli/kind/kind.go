package kind

import (
	"fmt"

	"github.com/gojek/stevedore/internal/cli/command"
)

var binaryPath string

func binary() string {
	if binaryPath == "" {
		return "kind"
	}
	return binaryPath
}

// Kind represents necessary information to
// mange k8s cluster
type Kind struct {
	command    string
	subCommand string
	args       []string
}

// Command returns kind command and its argument
func (kind Kind) Command() (string, []string) {
	args := []string{kind.command, kind.subCommand}
	return binary(), append(args, kind.args...)
}

// Print the underlying command
func (kind Kind) Print() string {
	return command.Print(kind)
}

// Environ return the environment variable needed
//to be used when running command
func (kind Kind) Environ() []string {
	return []string{}
}

func newKind(
	cmd,
	subCommand string,
	args []string,
) (command.Command, error) {
	helmCmd := Kind{
		command:    cmd,
		subCommand: subCommand,
		args:       args,
	}
	return command.NewCommand(helmCmd), nil
}

func globalFlags(name string) []string {
	return []string{
		fmt.Sprintf("--name=%s", name),
	}
}

// Create creates a new kind cluster
func Create(
	image string,
	kubeConfig string,
	name string,
	wait int,
) (command.Command, error) {
	args := []string{
		fmt.Sprintf("--image=%s", image),
		fmt.Sprintf("--kubeconfig=%s", kubeConfig),
	}
	if wait != 0 {
		args = append(args, fmt.Sprintf("--wait=%d", wait))
	}

	return newKind(
		"create",
		"cluster",
		append(args, globalFlags(name)...),
	)
}

// GetClusters returns the existing kind clusters
func GetClusters() (command.Command, error) {
	return newKind("get",
		"clusters",
		[]string{},
	)
}

// GetKubeConfig returns the kube config for the cluster
func GetKubeConfig(clusterName string) (command.Command, error) {
	args := []string{
		fmt.Sprintf("--name=%s", clusterName),
	}
	return newKind("get",
		"kubeconfig",
		args,
	)
}

// Delete deletes the kind cluster
func Delete(
	name string,
) (command.Command, error) {
	return newKind(
		"delete",
		"cluster",
		globalFlags(name),
	)
}
