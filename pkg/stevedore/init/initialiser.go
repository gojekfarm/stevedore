package init

import (
	"fmt"
	"github.com/gojek/stevedore/log"
	"github.com/gojek/stevedore/pkg/stevedore"
	k8sresources "github.com/gojek/stevedore/pkg/stevedore/k8s/resources"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/cmd/helm/installer"
)

// Initialiser exposes methods for preparing namespaces for stevedore
type Initialiser interface {
	Init(request []Request) (Responses, error)
}

// DefaultInitialiser prepares namespaces for stevedore
type DefaultInitialiser struct {
	resourceCreator k8sresources.ResourceCreator
	tillerInstaller TillerInstaller
}

func (d DefaultInitialiser) upstallTiller(opts installer.Options, upgrade bool) (Response, error) {
	if err := d.tillerInstaller.Install(&opts); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return Response{}, fmt.Errorf("error installing: %s", err)
		}
		if upgrade {
			if err := d.tillerInstaller.Upgrade(&opts); err != nil {
				return Response{}, fmt.Errorf("error when upgrading: %s", err)
			}
			return Response{opts.Namespace, "Tiller (the Helm server-side component) has been upgraded to the current version."}, nil
		}
		return Response{opts.Namespace, "Warning: Tiller is already installed in the cluster.\n" +
			"(Use --client-only to suppress this message, or --upgrade to upgrade Tiller to the current version.)"}, nil
	}
	return Response{opts.Namespace, "Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster."}, nil
}

// Init prepares namespaces for stevedore
func (d DefaultInitialiser) Init(requests []Request) (Responses, error) {
	var initResponses Responses
	for _, request := range requests {
		namespaceName := request.namespace
		log.Info(fmt.Sprintf("creating namespace %s", namespaceName))
		namespace := k8sresources.NewNamespace(namespaceName)
		if err := d.resourceCreator.CreateIfNotExists(namespace); err != nil {
			return nil, fmt.Errorf("failed to create namespace %s: %v", namespaceName, err)
		}

		serviceAccountName := request.serviceAccount
		log.Info(fmt.Sprintf("creating serviceaccount %s in namespace %s", serviceAccountName, namespaceName))
		serviceAccount := k8sresources.NewServiceAccount(serviceAccountName, namespaceName)
		if err := d.resourceCreator.CreateIfNotExists(serviceAccount); err != nil {
			return nil, fmt.Errorf("failed to create serviceaccount %s in namespace %s: %v", serviceAccountName, namespaceName, err)
		}

		roleBindingName := fmt.Sprintf("%s-cluster-admin", serviceAccountName)
		if namespaceName == stevedore.PrivilegedNamespace {
			log.Info(fmt.Sprintf("creating cluster rolebinding %s", roleBindingName))
			clusterRoleBinding := k8sresources.NewClusterAdminClusterRoleBinding(roleBindingName, namespaceName, serviceAccountName)
			if err := d.resourceCreator.CreateIfNotExists(clusterRoleBinding); err != nil {
				return nil, fmt.Errorf("failed to create cluster rolebinding %s: %v", roleBindingName, err)
			}
		} else {
			log.Info(fmt.Sprintf("creating rolebinding %s in namespace %s", roleBindingName, namespaceName))
			roleBinding := k8sresources.NewClusterAdminRoleBinding(roleBindingName, namespaceName, serviceAccountName)
			if err := d.resourceCreator.CreateIfNotExists(roleBinding); err != nil {
				return nil, fmt.Errorf("failed to create rolebinding %s in namespace %s: %v", roleBindingName, namespaceName, err)
			}
		}

		initResponse, err := d.upstallTiller(request.getOpts(), true)
		if err != nil {
			return nil, err
		}
		err = d.waitForTiller(request.getOpts())
		if err != nil {
			return nil, err
		}
		initResponses = append(initResponses, initResponse)
	}
	return initResponses, nil
}

func (d DefaultInitialiser) waitForTiller(opts installer.Options) error {
	return d.tillerInstaller.Wait(&opts)
}

// NewDefaultInitialiser creates DefaultInitialiser
func NewDefaultInitialiser(tillerInstaller TillerInstaller, resourceCreator k8sresources.ResourceCreator) DefaultInitialiser {
	return DefaultInitialiser{
		tillerInstaller: tillerInstaller,
		resourceCreator: resourceCreator,
	}
}

// CreateDefaultInitialiser creates DefaultInitialiser
func CreateDefaultInitialiser(client kubernetes.Interface, tillerWaitTimeout, tillerPollTimeout int) DefaultInitialiser {
	tillerInstaller := NewTillerInstaller(client, tillerWaitTimeout, tillerPollTimeout)
	resourceCreator := k8sresources.NewResourceCreator(client)
	return DefaultInitialiser{
		tillerInstaller: tillerInstaller,
		resourceCreator: resourceCreator,
	}
}
