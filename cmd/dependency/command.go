package dependency

import (
	"fmt"
	"strings"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/cmd/plugin"
	"github.com/gojek/stevedore/cmd/repo"
	"github.com/gojek/stevedore/cmd/store"
	"github.com/gojek/stevedore/pkg/manifest"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Action type represents manifest related actions
type Action interface {
	Do(impl manifest.ProviderImpl) error
}

// NewAction to create action based on command name
func NewAction(cmd Command) Action {
	switch cmd.name {
	case buildCommand:
		return NewBuildAction(afero.NewOsFs(), cmd.helmRepoName, cmd.artifactsPath)
	default:
		return ShowAction{}
	}
}

// Command can be used to create a cobra Command and appropriate flags to it
type Command struct {
	name          string
	artifactsPath string
	helmRepoName  string
	shortDesc     string
	longDesc      string
}

const (
	buildCommand = "build"
)

// NewBuildCommand creates a apply command
func NewBuildCommand() Command {
	return Command{
		name:      "build",
		shortDesc: "Build dependencies and push to chart museum",
		longDesc:  "Validate, build and push dependencies from stevedore manifest(s)",
	}
}

// NewShowCommand creates a apply command
func NewShowCommand() Command {
	return Command{
		name:      "show",
		shortDesc: "Show dependencies",
		longDesc:  "Show dependencies needed by stevedore manifest(s)",
	}
}

// CobraCommand builds a cobra command for the action
func (command Command) CobraCommand(fs afero.Fs, cfgFile *string, localStore store.Local) (*cobra.Command, error) {
	loader, err := plugin.GetPluginLoader()
	if err != nil {
		return nil, err
	}
	manifestPlugin, err := loader.GetManifestPlugin()
	if err != nil {
		return nil, err
	}
	cmd := cobra.Command{
		Use:           command.name,
		Short:         command.shortDesc,
		Long:          command.longDesc,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestProvider, err := manifestPlugin.ManifestProvider()
			if err != nil {
				return err
			}
			manifestProvider.Context = make(map[string]string)
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if strings.HasPrefix(flag.Name, manifestProvider.Name) {
					key := strings.TrimPrefix(flag.Name, fmt.Sprintf("%s-", manifestProvider.Name))
					manifestProvider.Context[key] = flag.Value.String()
				}
			})
			contextMap := map[string]string{provider.EnvironmentTypeKey: "dev"}
			manifestProvider.MergeToContext(contextMap)
			action := NewAction(command)
			return action.Do(manifestProvider)
		},
	}

	repo.AddRepoFlags(&cmd, &command.helmRepoName)
	cmd.PersistentFlags().StringVarP(&command.artifactsPath, "artifacts-path", "a", "", "Stevedore artifact(s) path (folder) to save the output as artifact")

	err = manifestPlugin.PopulateFlags(&cmd)
	if err != nil {
		return nil, err
	}

	return &cmd, nil
}
