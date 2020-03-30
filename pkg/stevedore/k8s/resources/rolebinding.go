package resources

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	apiGroup           = "rbac.authorization.k8s.io"
	clusterRoleKind    = "ClusterRole"
	serviceAccountKind = "ServiceAccount"
	clusterAdminRole   = "cluster-admin"
)

// RoleBinding represents a k8s RoleBinding resource
type RoleBinding struct {
	*rbacv1.RoleBinding
}

// Get gets a role binding
func (r RoleBinding) Get(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.RbacV1().RoleBindings(r.Namespace).Get(r.Name, metav1.GetOptions{})
	return err
}

// Create creates a role binding
func (r RoleBinding) Create(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.RbacV1().RoleBindings(r.Namespace).Create(r.RoleBinding)
	return err
}

// NewClusterAdminRoleBinding creates a new RoleBinding
// such that the given ServiceAccount is bound to cluster-admin ClusterRole
func NewClusterAdminRoleBinding(name, namespaceName, serviceAccountName string) RoleBinding {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: apiGroup,
			Kind:     clusterRoleKind,
			Name:     clusterAdminRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      serviceAccountKind,
				Name:      serviceAccountName,
				Namespace: namespaceName,
			},
		},
	}
	return RoleBinding{roleBinding}
}
