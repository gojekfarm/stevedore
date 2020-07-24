package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/cmd/store"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"

	"github.com/gojek/stevedore/pkg/stevedore"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage stevedore config",
}

var configViewCmd = &cobra.Command{
	Use:           "view",
	Short:         "Display complete configuration",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}
		return yaml.NewEncoder(os.Stdout).Encode(stevedoreConfig)
	},
}

var configGetContextsCmd = &cobra.Command{
	Use:           "get-contexts <optional-context-name>",
	Short:         "Describe one or many contexts",
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var givenContext string
		if len(args) == 1 && strings.TrimSpace(args[0]) != "" {
			givenContext = args[0]
		}

		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"CURRENT", "NAME", "TYPE", "ENVIRONMENT", "ENVIRONMENT TYPE", "KUBERNETES CONTEXT"})

		if givenContext != "" {
			for _, ctx := range stevedoreConfig.Contexts {
				current := ""
				if ctx.Name == stevedoreConfig.Current {
					current = "*"
				}
				if ctx.Name == givenContext {
					table.Append([]string{current, ctx.Name, ctx.Type, ctx.Environment, ctx.EnvironmentType, ctx.KubernetesContext})
					break
				}
			}
		} else {
			for _, ctx := range stevedoreConfig.Contexts {
				current := ""
				if ctx.Name == stevedoreConfig.Current {
					current = "*"
				}
				table.Append([]string{current, ctx.Name, ctx.Type, ctx.Environment, ctx.EnvironmentType, ctx.KubernetesContext})
			}
		}
		table.Render()
		return nil
	},
}

var configUseContextCmd = &cobra.Command{
	Use:           "use-context CONTEXT_NAME",
	Short:         "Sets the current-context in a stevedore config file",
	SilenceErrors: true,
	SilenceUsage:  true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("provide a single CONTEXT_NAME")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if current, ok := stevedore.OverriddenContext(localStore); ok {
			return fmt.Errorf("use-context cannot be used when context is overridden with %s env variable, current context is %s", stevedore.ContextEnv, current)
		}

		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}
		name := args[0]

		err = stevedoreConfig.Use(name)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully switched to context: %s\n", name)
		return nil
	},
}

var ctx = stevedore.Context{}
var localStore store.Local

const (
	nameFlag            = "name"
	typeFlag            = "type"
	environmentFlag     = "environment"
	environmentTypeFlag = "environment-type"
	kubeContextFlag     = "kube-context"
)

var addCtxtErrors = map[string]string{
	"Name":              fmt.Sprintf("Provide a name for stevedore context using --%s", nameFlag),
	"Type":              fmt.Sprintf("Provide a valid type (services|components|readonly) for stevedore context using --%s", typeFlag),
	"Environment":       fmt.Sprintf("Provide a environment for stevedore context using --%s", environmentFlag),
	"EnvironmentType":   fmt.Sprintf("Provide a environment type for stevedore context using --%s", environmentTypeFlag),
	"KubernetesContext": fmt.Sprintf("Provide a kubecontext for stevedore context using --%s", kubeContextFlag),
}

type contextErrors []string

func newContextErrors(source validator.ValidationErrors) contextErrors {
	errors := contextErrors{}
	for _, err := range source {
		errors = append(errors, addCtxtErrors[err.Field()])
	}
	return errors
}

func (errors contextErrors) Error() string {
	buff := bytes.NewBufferString("")
	for _, err := range errors {
		buff.WriteString(err)
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

var configAddContextCmd = &cobra.Command{
	Use:           "add-context",
	Short:         "Adds context to stevedore config file",
	SilenceErrors: true,
	SilenceUsage:  true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := ctx.IsValid()
		if ctxErrs, ok := err.(validator.ValidationErrors); ok {
			return newContextErrors(ctxErrs)
		}
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}

		err = stevedoreConfig.Add(ctx)
		if err != nil {
			return err
		}

		fmt.Println("Successfully added the below context:")

		table := tablewriter.NewWriter(os.Stdout)
		table.Append([]string{"Name", ctx.Name})
		table.Append([]string{"Type", ctx.Type})
		table.Append([]string{"Environment", ctx.Environment})
		table.Append([]string{"Environment Type", ctx.EnvironmentType})
		table.Append([]string{"Kubernetes Context", ctx.KubernetesContext})
		table.Render()

		err = stevedoreConfig.Use(ctx.Name)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully switched to context: %s\n", ctx.Name)
		return nil
	},
}

var configDeleteContextCmd = &cobra.Command{
	Use:           "delete-context CONTEXT_NAME",
	Short:         "Delete the specified context from the stevedore config",
	SilenceErrors: true,
	SilenceUsage:  true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("provide a single CONTEXT_NAME")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}
		name := args[0]

		err = stevedoreConfig.Delete(name)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully deleted context: %s\n", name)
		return nil
	},
}

var configRenameContextCmd = &cobra.Command{
	Use:           "rename-context OLD_CONTEXT_NAME NEW_CONTEXT_NAME",
	Short:         "Renames a context from the stevedore config file",
	SilenceErrors: true,
	SilenceUsage:  true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("provide both OLD_CONTEXT_NAME and NEW_CONTEXT_NAME")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}

		oldContextName := args[0]
		newContextName := args[1]

		err = stevedoreConfig.Rename(oldContextName, newContextName)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully renamed %s to: %s\n", oldContextName, newContextName)
		return nil
	},
}

var configShowContextCmd = &cobra.Command{
	Use:           "show-context",
	Short:         "Shows the current-context",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}

		context, err := stevedoreConfig.CurrentContext()
		if err != nil {
			return err
		}
		cli.PrintYaml(context)
		return nil
	},
}

func init() {
	localStore = store.Local{}
	rootCmd.AddCommand(configCmd)

	configAddContextCmd.PersistentFlags().StringVar(&ctx.Name, nameFlag, "", "Stevedore context name")
	configAddContextCmd.PersistentFlags().StringVar(&ctx.Type, typeFlag, "", "Type of kubernetes cluster of the stevedore context (eg. components|services|readonly)")
	configAddContextCmd.PersistentFlags().StringVar(&ctx.Environment, environmentFlag, "", "Environment of the stevedore context")
	configAddContextCmd.PersistentFlags().StringVar(&ctx.EnvironmentType, environmentTypeFlag, "", "Type of Environment of stevedore context (eg. staging|production)")
	configAddContextCmd.PersistentFlags().StringVar(&ctx.KubernetesContext, kubeContextFlag, "", "Kubernetes cluster of the stevedore context")

	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configGetContextsCmd)
	configCmd.AddCommand(configUseContextCmd)
	configCmd.AddCommand(configAddContextCmd)
	configCmd.AddCommand(configDeleteContextCmd)
	configCmd.AddCommand(configRenameContextCmd)
	configCmd.AddCommand(configShowContextCmd)
}
