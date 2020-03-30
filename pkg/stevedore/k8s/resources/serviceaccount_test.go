package resources

import (
	"fmt"
	"testing"

	mockK8s "github.com/gojek/stevedore/pkg/internal/mocks/client-go/kubernetes"
	mock_v1 "github.com/gojek/stevedore/pkg/internal/mocks/client-go/kubernetes/typed/core/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewServiceAccount(t *testing.T) {
	actual := NewServiceAccount("jerry", "tom")

	serviceAccountSpec := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jerry",
			Namespace: "tom",
		},
	}
	expected := ServiceAccount{serviceAccountSpec}
	assert.Equal(t, expected, actual)
}

func TestServiceAccount_Get(t *testing.T) {
	t.Run("should get a serviceaccount", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		serviceAccountInterface := mock_v1.NewMockServiceAccountInterface(ctrl)
		coreV1Interface.EXPECT().ServiceAccounts("tom").Return(serviceAccountInterface)
		serviceAccountInterface.EXPECT().Get("jerry", metav1.GetOptions{}).Return(&corev1.ServiceAccount{}, nil)
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		sa := NewServiceAccount("jerry", "tom")
		err := sa.Get(kubeClient)

		assert.NoError(t, err)
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		namespaceInterface := mock_v1.NewMockServiceAccountInterface(ctrl)
		coreV1Interface.EXPECT().ServiceAccounts("tom").Return(namespaceInterface)
		namespaceInterface.EXPECT().Get("jerry", metav1.GetOptions{}).Return(&corev1.ServiceAccount{}, fmt.Errorf("failed to get serviceaccount"))
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		sa := NewServiceAccount("jerry", "tom")
		err := sa.Get(kubeClient)

		assert.EqualError(t, err, "failed to get serviceaccount")
	})
}

func TestServiceAccount_Create(t *testing.T) {
	t.Run("should create a serviceaccount", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		serviceAccountInterface := mock_v1.NewMockServiceAccountInterface(ctrl)
		coreV1Interface.EXPECT().ServiceAccounts("tom").Return(serviceAccountInterface)
		serviceAccountSpec := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jerry",
				Namespace: "tom",
			},
		}
		serviceAccountInterface.EXPECT().Create(serviceAccountSpec).Return(&corev1.ServiceAccount{}, nil)
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		sa := NewServiceAccount("jerry", "tom")
		err := sa.Create(kubeClient)

		assert.NoError(t, err)
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		namespaceInterface := mock_v1.NewMockServiceAccountInterface(ctrl)
		coreV1Interface.EXPECT().ServiceAccounts("tom").Return(namespaceInterface)
		serviceAccountSpec := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jerry",
				Namespace: "tom",
			},
		}
		namespaceInterface.EXPECT().Create(serviceAccountSpec).Return(&corev1.ServiceAccount{}, fmt.Errorf("failed to create serviceaccount"))
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		sa := NewServiceAccount("jerry", "tom")
		err := sa.Create(kubeClient)

		assert.EqualError(t, err, "failed to create serviceaccount")
	})
}
