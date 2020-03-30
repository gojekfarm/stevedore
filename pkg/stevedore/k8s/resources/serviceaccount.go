package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceAccount represents a k8s ServiceAccount
type ServiceAccount struct {
	*corev1.ServiceAccount
}

// Get gets a service account
func (sa ServiceAccount) Get(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.Name, metav1.GetOptions{})
	return err
}

// Create creates a service account
func (sa ServiceAccount) Create(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.CoreV1().ServiceAccounts(sa.Namespace).Create(sa.ServiceAccount)
	return err
}

// NewServiceAccount creates a new ServiceAccount
func NewServiceAccount(name, namespaceName string) ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
	}
	return ServiceAccount{serviceAccount}
}
