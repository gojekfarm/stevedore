package stevedore

import (
	"fmt"
	"github.com/gojek/stevedore/internal/cli/command"
	"github.com/gojek/stevedore/pkg/stevedore"
)

var binaryPath string

func binary() string {
	if binaryPath == "" {
		return "stevedore"
	}
	return binaryPath
}

// Stevedore represents stevedore cli command
type Stevedore struct {
	command    string
	args       []string
	config     string
	kubeConfig string
	context    string
}

// Environ return the environment variable needed
//to be used when running command
func (stevedore Stevedore) Environ() []string {
	return []string{fmt.Sprintf("STEVEDORE_CONTEXT=%s", stevedore.context)}
}

// Print the underlying command
func (stevedore Stevedore) Print() string {
	return command.Print(stevedore)
}

// Command returns the stevedore command
/// with all the necessary parameters
func (stevedore Stevedore) Command() (string, []string) {
	args := []string{
		stevedore.command,
		fmt.Sprintf("--config=%s", stevedore.config),
		fmt.Sprintf("--kubeconfig=%s", stevedore.kubeConfig),
	}
	return binary(), append(args, stevedore.args...)
}

func saveManifest(manifest stevedore.Manifest) (string, error) {
	manifestFile, err := command.SaveManifest(manifest)
	if err != nil {
		return "", fmt.Errorf("unable to save manifest, reason %v", err)
	}
	return manifestFile, nil
}

func saveOverrides(overrides stevedore.Overrides) (string, error) {
	manifestFile, err := command.SaveOverrides(overrides)
	if err != nil {
		return "", fmt.Errorf("unable to save manifest, reason %v", err)
	}
	return manifestFile, nil
}

func saveEnvs(envs stevedore.Env) (string, error) {
	manifestFile, err := command.SaveEnvs(envs)
	if err != nil {
		return "", fmt.Errorf("unable to save manifest, reason %v", err)
	}
	return manifestFile, nil
}

func save(
	manifest stevedore.Manifest,
	overrides stevedore.Overrides,
	env stevedore.Env,
) (
	manifestPath string,
	overridesPath string,
	envPath string,
	err error,
) {
	manifestPath, err = saveManifest(manifest)
	if err != nil {
		return "", "", "", err
	}

	overridesPath, err = saveOverrides(overrides)
	if err != nil {
		return "", "", "", err
	}

	envPath, err = saveEnvs(env)
	if err != nil {
		return "", "", "", err
	}

	return manifestPath, overridesPath, envPath, err
}

func newStevedore(
	cmd string,
	config string,
	context string,
	kubeConfig string,
	args []string,
) (command.Command, error) {
	stevedoreCmd := &Stevedore{
		command:    cmd,
		args:       args,
		config:     config,
		context:    context,
		kubeConfig: kubeConfig,
	}

	stevedoreCommand := command.NewCommand(stevedoreCmd)
	return stevedoreCommand, nil
}

// Apply invoke stevedore apply
func Apply(
	config string,
	context string,
	kubeConfig string,
	manifest stevedore.Manifest,
	overrides stevedore.Overrides,
	envs stevedore.Env,
	helmRepoName string,
	confirm bool,
) (command.Command, error) {
	manifestPath, overridesPath, envPath, err := save(manifest, overrides, envs)
	if err != nil {
		return command.Command{}, err
	}

	args := []string{fmt.Sprintf("--manifests-path=%s", manifestPath),
		fmt.Sprintf("--overrides-path=%s", overridesPath),
		fmt.Sprintf("--envs-path=%s", envPath),
		fmt.Sprintf("--helm-repo-name=%s", helmRepoName),
	}
	if confirm {
		args = append(args, "--yes")
	}

	return newStevedore("apply", config, context, kubeConfig, args)
}

// Init invoke stevedore init
func Init(
	config string,
	context string,
	kubeConfig string,
	namespacesFilePath string,
	timeout string,
) (command.Command, error) {
	if timeout == "" {
		timeout = "60"
	}

	args := []string{
		fmt.Sprintf("--namespaces-file=%s", namespacesFilePath),
		fmt.Sprintf("--timeout=%s", timeout),
	}

	return newStevedore("init", config, context, kubeConfig, args)
}
