package resources

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ClusterRoleBinding represents a k8s ClusterRoleBinding resource
type ClusterRoleBinding struct {
	*rbacv1.ClusterRoleBinding
}

// Get gets a cluster role binding
func (r ClusterRoleBinding) Get(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.RbacV1().ClusterRoleBindings().Get(r.Name, metav1.GetOptions{})
	return err
}

// Create creates a cluster role binding
func (r ClusterRoleBinding) Create(kubeClient kubernetes.Interface) error {
	_, err := kubeClient.RbacV1().ClusterRoleBindings().Create(r.ClusterRoleBinding)
	return err
}

// NewClusterAdminClusterRoleBinding creates a new ClusterRoleBinding
// such that the given ServiceAccount is bound to cluster-admin ClusterRole
func NewClusterAdminClusterRoleBinding(name, namespaceName, serviceAccountName string) ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
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
	return ClusterRoleBinding{clusterRoleBinding}
}
