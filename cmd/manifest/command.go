package manifest

import (
	"fmt"
	"os"
	"strings"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/cmd/kubeconfig"
	"github.com/gojek/stevedore/cmd/repo"
	"github.com/gojek/stevedore/pkg/config"
	"github.com/spf13/pflag"

	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/cmd/plugin"
	"github.com/gojek/stevedore/cmd/store"
	"github.com/manifoldco/promptui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Command can be used to create a cobra Command and appropriate flags to it
type Command struct {
	overridesPath      string
	helmRepoName       string
	name               string
	dryRun             bool
	fs                 afero.Fs
	cfgFile            *string
	kubeconfig         string
	kubeconfigRequired bool
	useHelm            bool
	askConfirmation    bool
	envsPath           string
	artifactsPath      string
	confirm            bool
	helmTimeout        int64
}

const (
	applyCommand  = "apply"
	planCommand   = "plan"
	renderCommand = "render"
)

// NewApplyCmd creates a apply command
func NewApplyCmd(fs afero.Fs, cfgFile *string, kubeconfigRequired bool) *Command {
	return &Command{
		name:               applyCommand,
		dryRun:             false,
		fs:                 fs,
		cfgFile:            cfgFile,
		useHelm:            true,
		askConfirmation:    true,
		confirm:            false,
		kubeconfigRequired: kubeconfigRequired,
	}
}

// NewPlanCmd creates a plan command
func NewPlanCmd(fs afero.Fs, cfgFile *string, kubeconfigRequired bool) *Command {
	return &Command{name: planCommand,
		dryRun:             true,
		fs:                 fs,
		cfgFile:            cfgFile,
		useHelm:            true,
		kubeconfigRequired: kubeconfigRequired,
	}
}

// NewRenderCmd creates a render command
func NewRenderCmd(fs afero.Fs, cfgFile *string, kubeconfigRequired bool) *Command {
	return &Command{name: renderCommand,
		dryRun:             true,
		fs:                 fs,
		cfgFile:            cfgFile,
		kubeconfigRequired: kubeconfigRequired,
	}
}

func promptConfirmation() (string, error) {
	templates := &promptui.PromptTemplates{
		Confirm: "{{ . | yellow }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | green }}",
	}

	prompt := promptui.Prompt{
		Label:     "Confirm to apply: [y/N] ",
		IsConfirm: true,
		Templates: templates,
	}

	return prompt.Run()
}

// CobraCommand builds a cobra command for the action
func (actionCmd *Command) CobraCommand() (*cobra.Command, error) {
	shortDesc := fmt.Sprintf("%s stevedore yaml(s)", strings.Title(actionCmd.name))
	longDesc := fmt.Sprintf("Validate and %s stevedore yaml(s)", actionCmd.name)
	loader, err := plugin.GetPluginLoader()
	if err != nil {
		return nil, err
	}
	plugins, err := loader.GetAllEligiblePlugins()
	if err != nil {
		return nil, err
	}
	cmd := cobra.Command{
		Use:           actionCmd.name,
		Short:         shortDesc,
		Long:          longDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(actionCmd.overridesPath); actionCmd.overridesPath != "" && os.IsNotExist(err) {
				return fmt.Errorf("invalid file path. Provide a valid path to stevedore manifests using --overrides-path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			localStore := store.Local{}
			contextProvider := provider.NewContextProvider(actionCmd.fs, *actionCmd.cfgFile, localStore)
			overridesProvider := provider.NewOverrideProvider(actionCmd.fs, actionCmd.overridesPath)
			envProvider := provider.NewEnvProvider(actionCmd.fs, actionCmd.envsPath)
			reporter := DefaultReporter{}

			ctx, err := contextProvider.Context()
			if err != nil {
				return err
			}
			cli.Infof(ctx.String())

			resolvedKubeconfig, err := kubeconfig.ResolveAndValidate(kubeconfig.OSHomeDirResolver, actionCmd.kubeconfig, actionCmd.fs, ctx)
			if err != nil {
				return err
			}
			actionCmd.kubeconfig = resolvedKubeconfig

			if actionCmd.askConfirmation && !actionCmd.confirm {
				_, err := promptConfirmation()

				if err != nil {
					fmt.Printf("Command Cancelled %v\n", err)
					return err
				}
			}

			configProviders, err := plugins.ConfigProviders()
			if err != nil {
				return err
			}
			manifestProvider, err := plugins.ManifestProvider()
			if err != nil {
				return err
			}

			flags := cmd.Flags()
			updatedConfigProviders := make(config.Providers, 0, len(configProviders))
			for _, configProvider := range configProviders {
				configProvider.Context = make(map[string]string)
				flags.VisitAll(func(flag *pflag.Flag) {
					if strings.HasPrefix(flag.Name, configProvider.Name) {
						key := strings.TrimPrefix(flag.Name, fmt.Sprintf("%s-", configProvider.Name))
						configProvider.Context[key] = flag.Value.String()
					}
				})
				updatedConfigProviders = append(updatedConfigProviders, configProvider)
			}
			manifestProvider.Context = make(map[string]string)
			flags.VisitAll(func(flag *pflag.Flag) {
				if strings.HasPrefix(flag.Name, manifestProvider.Name) {
					key := strings.TrimPrefix(flag.Name, fmt.Sprintf("%s-", manifestProvider.Name))
					manifestProvider.Context[key] = flag.Value.String()
				}
			})
			ignoreProvider, err := provider.NewIgnoreProvider(actionCmd.fs, manifestProvider.Context["path"], localStore)
			if err != nil {
				return err
			}

			info, err := NewManifests(
				localStore,
				contextProvider,
				manifestProvider,
				overridesProvider,
				ignoreProvider,
				envProvider,
				reporter,
				updatedConfigProviders,
			)
			if err != nil {
				return err
			}

			generateArtifact := len(actionCmd.artifactsPath) != 0
			artifact := NewArtifact(actionCmd.fs, generateArtifact, actionCmd.artifactsPath)

			action := NewAction(*actionCmd, *info)
			processedInfo, err := action.Do()
			if err != nil {
				return err
			}
			return artifact.Save(processedInfo)
		},
	}
	cmd.PersistentFlags().StringVarP(&actionCmd.artifactsPath, "artifacts-path", "a", "", "Stevedore artifact(s) path (folder) to save the output as artifact")
	cmd.PersistentFlags().StringVarP(&actionCmd.envsPath, "envs-path", "e", "", "Stevedore env(s) path (can be yaml file or folder)")
	cmd.PersistentFlags().StringVarP(&actionCmd.overridesPath, "overrides-path", "o", "", "Stevedore overrides path (can be yaml file or folder)")

	if actionCmd.askConfirmation {
		cmd.PersistentFlags().BoolVar(&actionCmd.confirm, "yes", actionCmd.confirm, "Confirm to apply")
	}

	if actionCmd.useHelm {
		repo.AddRepoFlags(&cmd, &actionCmd.helmRepoName)
		cmd.PersistentFlags().Int64VarP(&actionCmd.helmTimeout, "helm-timeout", "t", 600, "Timeout in seconds(default 10 minutes)")
	}

	if actionCmd.kubeconfigRequired {
		defaultFile, err := kubeconfig.DefaultFile(kubeconfig.OSHomeDirResolver)
		if err != nil {
			return nil, err
		}
		cmd.PersistentFlags().StringVar(&actionCmd.kubeconfig, "kubeconfig", "", fmt.Sprintf("path to kubeconfig file (default: %s)", defaultFile))
	}

	err = plugins.PopulateFlags(&cmd)
	if err != nil {
		return nil, err
	}

	return &cmd, nil
}
