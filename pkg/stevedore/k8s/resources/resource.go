package resources

import (
	"fmt"
	"github.com/gojek/stevedore/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

// Resource represents a k8s resource
type Resource interface {
	Get(kubeClient kubernetes.Interface) error
	Create(kubeClient kubernetes.Interface) error
}

// ResourceCreator exposes methods to create k8s resources in various ways
type ResourceCreator interface {
	CreateIfNotExists(resource Resource) error
}

type defaultResourceCreator struct {
	kubeClient kubernetes.Interface
}

// NewResourceCreator creates a defaultResourceCreator
func NewResourceCreator(kubeClient kubernetes.Interface) ResourceCreator {
	return defaultResourceCreator{kubeClient: kubeClient}
}

// CreateIfNotExists creates a k8s resource if it does not exists already
func (d defaultResourceCreator) CreateIfNotExists(resource Resource) error {
	err := resource.Get(d.kubeClient)
	if err != nil {
		errStatus, ok := err.(apierrors.APIStatus)
		if !ok || errStatus.Status().Code != http.StatusNotFound {
			errMsg := fmt.Sprintf("unable to get resource: %v", err)
			log.Error(errMsg)
			return fmt.Errorf(errMsg)
		}

		err = resource.Create(d.kubeClient)
		if err != nil {
			errMsg := fmt.Sprintf("unable to create resource: %v", err)
			log.Error(errMsg)
			return fmt.Errorf(errMsg)
		}
		log.Info("created resource successfully")
		return nil
	}

	log.Info("resource already exists")
	return nil
}
