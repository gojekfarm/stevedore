package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Namespace represents a k8s Namespace
type Namespace struct {
	*corev1.Namespace
}

// Get gets a namespace
func (n Namespace) Get(kubeClient kubernetes.Interface) error {
	namespaces := kubeClient.CoreV1().Namespaces()
	_, err := namespaces.Get(n.Name, metav1.GetOptions{})
	return err
}

// Create creates a namespace
func (n Namespace) Create(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.CoreV1().Namespaces().Create(n.Namespace)
	return err
}

// NewNamespace creates a new Namespace
func NewNamespace(name string) Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return Namespace{namespace}
}
