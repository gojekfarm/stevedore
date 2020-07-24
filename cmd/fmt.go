package cmd

import (
	"fmt"
	"github.com/gojek/stevedore/client/yaml"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

var (
	fmtEnvsPath      string
	fmtOverridesPath string
	fmtManifestPath  string
)

type formatErrors []string

func (fmtErrors formatErrors) Error() string {
	builder := strings.Builder{}
	builder.WriteString("Unable to format, reason(s): \n")
	for _, fmtError := range fmtErrors {
		builder.WriteString(fmtError)
	}
	return builder.String()
}

type fileContent struct {
	fileName string
	content  fmt.Formatter
}

type fileContents []fileContent

type readFileContent func(reader io.Reader) (fmt.Formatter, error)

type contentFetcher func(fs afero.Fs) (fileContents, error)

var fmtCmd = &cobra.Command{
	Use:           "fmt",
	Short:         "Format stevedore yaml(s)",
	Long:          `Format stevedore yaml(s) and save it back`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var all fileContents
		fetchers := getFetchers()
		for _, fetch := range fetchers {
			contents, err := fetch(fs)
			if err != nil {
				return err
			}
			all = append(all, contents...)
		}
		return format(fs, all...)
	},
}

func getFetchers() []contentFetcher {
	fetchers := make([]contentFetcher, 0)
	if fmtManifestPath != "" {
		fetchers = append(fetchers, getManifests)
	}

	if fmtOverridesPath != "" {
		fetchers = append(fetchers, getOverrides)
	}

	if fmtEnvsPath != "" {
		fetchers = append(fetchers, getEnvs)
	}
	return fetchers
}

func getContents(fs afero.Fs, path string, readFile readFileContent) (fileContents, error) {
	files, err := yaml.NewYamlFiles(fs, path)
	if err != nil {
		return nil, fmt.Errorf("unable to format file %s, reason %v", path, err)
	}

	errors := formatErrors{}
	result := fileContents{}
	for _, file := range files {
		content, err := readFile(file.Reader())
		if err != nil {
			errors = append(errors, fmt.Sprintf("unable to read content %s, reason: %v", file.Name, err.Error()))
			continue
		}
		result = append(result, fileContent{fileName: file.Name, content: content})
	}
	if len(errors) != 0 {
		return nil, errors
	}
	return result, nil
}

func getOverrides(fs afero.Fs) (fileContents, error) {
	return getContents(fs, fmtOverridesPath, func(reader io.Reader) (fmt.Formatter, error) {
		return stevedore.NewOverrides(reader)
	})
}

func getEnvs(fs afero.Fs) (fileContents, error) {
	return getContents(fs, fmtEnvsPath, func(reader io.Reader) (fmt.Formatter, error) {
		return stevedore.NewEnv(reader)
	})
}

func getManifests(fs afero.Fs) (fileContents, error) {
	return getContents(fs, fmtManifestPath, func(reader io.Reader) (fmt.Formatter, error) {
		return stevedore.NewManifest(reader)
	})
}

func format(fs afero.Fs, fileContents ...fileContent) error {
	errors := formatErrors{}
	for _, fileContent := range fileContents {
		formatted := fmt.Sprintf("%y", fileContent.content)
		err := afero.WriteFile(fs, fileContent.fileName, []byte(formatted), 0644)

		if err != nil {
			errors = append(errors, fmt.Sprintf("unable to write override %s, reason: %v", fileContent.fileName, err.Error()))
			continue
		}
	}
	if len(errors) != 0 {
		return errors
	}
	return nil
}

func init() {
	rootCmd.ResetFlags()
	fmtCmd.PersistentFlags().StringVarP(&fmtManifestPath, "manifests-path", "f", "", "Stevedore manifest(s) path (can be yaml file or folder)")
	fmtCmd.PersistentFlags().StringVarP(&fmtEnvsPath, "envs-path", "e", "", "Stevedore env(s) path (can be yaml file or folder)")
	fmtCmd.PersistentFlags().StringVarP(&fmtOverridesPath, "overrides-path", "o", "", "Stevedore overrides path (can be yaml file or folder)")
	rootCmd.AddCommand(fmtCmd)
}
