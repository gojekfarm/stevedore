package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gojek/stevedore/client/yaml"
	"github.com/gojek/stevedore/cmd/cli"
	kubeConfigPkg "github.com/gojek/stevedore/cmd/kubeconfig"
	"github.com/gojek/stevedore/cmd/plugin"
	"github.com/gojek/stevedore/pkg/stevedore"
	stevedoreInit "github.com/gojek/stevedore/pkg/stevedore/init"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/pkg/kube"
)

type namespaces struct {
	Names []string `yaml:"namespaces"`
}

const tillerServiceAccount = "tiller"
const tillerPollTimeout = 10

var kubeConfig string
var timeout int
var namespacesFile string

// CobraCommand builds a cobra command for the action
func CobraCommand() (*cobra.Command, error) {
	loader, err := plugin.GetPluginLoader()
	if err != nil {
		return nil, err
	}
	plugins, err := loader.GetManifestPlugin()
	if err != nil {
		return nil, err
	}
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Install/Upgrade tiller",
		Long: `Install/Upgrade tiller in the namespace inferred through namespaces file or manifests files.
For every inferred namespace this command will:
1. create a ServiceAccount called tiller
2. binds tiller ServiceAccount to the "cluster-admin" ClusterRole as RoleBinding
3. install tiller server and configures it to use the tiller ServiceAccount created in 1`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			manifestsPath, err := cmd.Flags().GetString("manifests-path")
			if err != nil {
				return err
			}

			if namespacesFile == "" && manifestsPath == "" {
				return fmt.Errorf("either --namespaces-file or --manifests-path is required for inferring namespaces")
			}

			if namespacesFile != "" && manifestsPath != "" {
				return fmt.Errorf("only one of --namespaces-file or --manifests-path should be enough for inferring namespaces")
			}

			if namespacesFile != "" {
				if _, err := os.Stat(namespacesFile); err != nil {
					return fmt.Errorf("invalid file path. Provide a valid path to stevedore namespace file using --namespaces-file")
				}
				// mark --manifests-path flag as not required
				// to avoid flag value required error from manifests provider plugin flag
				err := cmd.PersistentFlags().SetAnnotation("manifests-path", cobra.BashCompOneRequiredFlag, []string{"false"})

				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			stevedoreConfig, err := stevedore.NewConfigurationFromFile(fs, cfgFile, localStore)
			if err != nil {
				return err
			}

			context, err := stevedoreConfig.CurrentContext()
			if err != nil {
				return err
			}

			resolvedKubeConfig, err := kubeConfigPkg.ResolveAndValidate(kubeConfigPkg.OSHomeDirResolver, kubeConfig, fs, context)
			if err != nil {
				return err
			}
			kubeConfig = resolvedKubeConfig

			client, err := getKubeClient(kubeConfig, context.KubernetesContext)
			if err != nil {
				return err
			}

			initialiser := stevedoreInit.CreateDefaultInitialiser(client, timeout, tillerPollTimeout)

			var initRequests []stevedoreInit.Request

			if namespacesFile != "" {
				yamlFile, err := yaml.NewYamlFile(fs, namespacesFile)
				if err != nil {
					return err
				}
				initRequests, err = getRequests(yamlFile)
				if err != nil {
					return err
				}
			} else {
				manifestProvider, err := plugins.ManifestProvider()
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
				stevedoreContext, err := context.Map()
				if err != nil {
					return err
				}
				manifestProvider.MergeToContext(stevedoreContext)
				manifestFiles, err := manifestProvider.Provider.Manifests(manifestProvider.Context)
				if err != nil {
					return err
				}
				for _, namespace := range manifestFiles.AllNamespaces() {
					initRequest := stevedoreInit.NewInitRequest(namespace, tillerServiceAccount, forceUpgrade)
					initRequests = append(initRequests, initRequest)
				}
			}

			initResponses, err := initialiser.Init(initRequests)
			if err != nil {
				return err
			}
			cli.Info(initResponses)
			return nil
		},
	}
	err = plugins.PopulateFlags(initCmd)
	if err != nil {
		return nil, err
	}
	return initCmd, nil
}

func getRequests(file yaml.File) ([]stevedoreInit.Request, error) {
	namespaces := &namespaces{}
	err := stevedore.ValidateAndGenerate(file.Reader(), namespaces)
	if err != nil {
		return nil, err
	}
	return createInitRequests(*namespaces), nil
}

func createInitRequests(n namespaces) []stevedoreInit.Request {
	namespaceNames := n.Names
	requests := make([]stevedoreInit.Request, 0, len(namespaceNames))
	for _, namespace := range namespaceNames {
		requests = append(requests, stevedoreInit.NewInitRequest(namespace, tillerServiceAccount, forceUpgrade))
	}
	return requests
}

func getKubeClient(kubeConfig, kubeContext string) (kubernetes.Interface, error) {
	clientConfig, err := kube.GetConfig(kubeContext, kubeConfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error when creating kubernetes client config due to %v", err)
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes config due to %v", err)
	}
	return client, nil
}

var forceUpgrade bool

func init() {
	initCmd, err := CobraCommand()
	if err != nil {
		cli.DieIf(err, closePlugins)
	}
	rootCmd.AddCommand(initCmd)

	defaultFile, err := kubeConfigPkg.DefaultFile(kubeConfigPkg.OSHomeDirResolver)
	if err != nil {
		cli.DieIf(err, closePlugins)
	}
	initCmd.Flags().IntVar(&timeout, "timeout", 60, "time in seconds to wait for tiller to become operational")
	initCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", fmt.Sprintf("path to kubeConfig file (default: %s)", defaultFile))
	initCmd.PersistentFlags().BoolVar(&forceUpgrade, "force-upgrade", false, "force upgrade Tiller to the helm version supported by stevedore")
	initCmd.PersistentFlags().StringVarP(&namespacesFile, "namespaces-file", "n", "", "Stevedore namespaces yaml file")
}
