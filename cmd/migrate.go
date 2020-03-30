package cmd

import (
	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/cmd/migrate"
	"github.com/gojek/stevedore/cmd/store"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	manifestFilePath string
	overridesPath    string
	optimizeMigrate  bool
	envPath          string
)

var migrateCmd = &cobra.Command{
	Use:           "migrate",
	Short:         "Migrate stevedore manifest file(s)",
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()
		localStore := store.Local{}

		configuration, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
		if err != nil {
			return err
		}

		ignoreProvider, err := provider.NewIgnoreProvider(fs, manifestFilePath, localStore)
		if err != nil {
			return err
		}

		ignoreFiles, err := ignoreProvider.Files()
		if err != nil {
			return err
		}

		manifestStrategy := migrate.NewManifestStrategy(fs, manifestFilePath, configuration.Contexts, optimizeMigrate)
		ignoreStrategy := migrate.NewIgnoreStrategy(fs, ignoreFiles)
		overrideStrategy := migrate.NewOverrideStrategy(fs, overridesPath, optimizeMigrate)
		envStrategy := migrate.NewEnvStrategy(fs, envPath, optimizeMigrate)
		return migrate.Perform(ignoreStrategy, manifestStrategy, overrideStrategy, envStrategy)
	},
}

func init() {
	migrateCmd.Flags().StringVarP(&manifestFilePath, "manifests-path", "f", "services", "Stevedore manifest(s) path (can be yaml file or folder)")
	migrateCmd.Flags().BoolVarP(&optimizeMigrate, "optimize", "", false, "Optimize migration")
	migrateCmd.Flags().StringVarP(&overridesPath, "overrides-path", "o", "overrides", "Stevedore overrides path (can be yaml file or folder)")
	migrateCmd.Flags().StringVarP(&envPath, "env-path", "e", "envs", "Stevedore env path (can be yaml file or folder)")
	rootCmd.AddCommand(migrateCmd)
}
