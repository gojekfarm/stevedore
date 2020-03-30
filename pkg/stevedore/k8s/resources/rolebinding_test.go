package resources

import (
	"fmt"
	"testing"

	mockK8s "github.com/gojek/stevedore/pkg/internal/mocks/client-go/kubernetes"
	mock_v1 "github.com/gojek/stevedore/pkg/internal/mocks/client-go/kubernetes/typed/rbac/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewClusterAdminRoleBinding(t *testing.T) {
	actual := NewClusterAdminRoleBinding("jerry-cluster-admin", "rolebinding", "jerry")

	roleBindingSpec := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jerry-cluster-admin",
			Namespace: "rolebinding",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "jerry",
				Namespace: "rolebinding",
			},
		},
	}
	expected := RoleBinding{roleBindingSpec}
	assert.Equal(t, expected, actual)
}

func TestRoleBinding_Get(t *testing.T) {
	t.Run("should get a rolebinding", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		rbacV1Interface := mock_v1.NewMockRbacV1Interface(ctrl)
		kubeClient.EXPECT().RbacV1().Return(rbacV1Interface)
		roleBindingsInterface := mock_v1.NewMockRoleBindingInterface(ctrl)
		rbacV1Interface.EXPECT().RoleBindings("rolebinding").Return(roleBindingsInterface)
		roleBindingsInterface.EXPECT().Get("jerry-cluster-admin", metav1.GetOptions{}).Return(&rbacv1.RoleBinding{}, nil)

		rb := NewClusterAdminRoleBinding("jerry-cluster-admin", "rolebinding", "jerry")
		err := rb.Get(kubeClient)

		assert.NoError(t, err)
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		rbacV1Interface := mock_v1.NewMockRbacV1Interface(ctrl)
		kubeClient.EXPECT().RbacV1().Return(rbacV1Interface)
		roleBindingsInterface := mock_v1.NewMockRoleBindingInterface(ctrl)
		rbacV1Interface.EXPECT().RoleBindings("rolebinding").Return(roleBindingsInterface)
		roleBindingsInterface.EXPECT().Get("jerry-cluster-admin", metav1.GetOptions{}).Return(&rbacv1.RoleBinding{}, fmt.Errorf("failed to get rolebinding"))

		rb := NewClusterAdminRoleBinding("jerry-cluster-admin", "rolebinding", "jerry")
		err := rb.Get(kubeClient)

		assert.EqualError(t, err, "failed to get rolebinding")
	})
}

func TestRoleBinding_Create(t *testing.T) {
	t.Run("should create a rolebinding", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		rbacV1Interface := mock_v1.NewMockRbacV1Interface(ctrl)
		kubeClient.EXPECT().RbacV1().Return(rbacV1Interface)
		roleBindingInterface := mock_v1.NewMockRoleBindingInterface(ctrl)
		rbacV1Interface.EXPECT().RoleBindings("rolebinding").Return(roleBindingInterface)
		roleBindingSpec := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jerry-cluster-admin",
				Namespace: "rolebinding",
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "cluster-admin",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "jerry",
					Namespace: "rolebinding",
				},
			},
		}
		roleBindingInterface.EXPECT().Create(roleBindingSpec).Return(&rbacv1.RoleBinding{}, nil)

		rb := NewClusterAdminRoleBinding("jerry-cluster-admin", "rolebinding", "jerry")
		err := rb.Create(kubeClient)

		assert.NoError(t, err)
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		rbacV1Interface := mock_v1.NewMockRbacV1Interface(ctrl)
		kubeClient.EXPECT().RbacV1().Return(rbacV1Interface)
		roleBindingInterface := mock_v1.NewMockRoleBindingInterface(ctrl)
		rbacV1Interface.EXPECT().RoleBindings("rolebinding").Return(roleBindingInterface)
		roleBindingSpec := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jerry-cluster-admin",
				Namespace: "rolebinding",
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "cluster-admin",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "jerry",
					Namespace: "rolebinding",
				},
			},
		}
		roleBindingInterface.EXPECT().Create(roleBindingSpec).Return(&rbacv1.RoleBinding{}, fmt.Errorf("failed to create rolebinding"))

		rb := NewClusterAdminRoleBinding("jerry-cluster-admin", "rolebinding", "jerry")
		err := rb.Create(kubeClient)

		assert.EqualError(t, err, "failed to create rolebinding")
	})
}
