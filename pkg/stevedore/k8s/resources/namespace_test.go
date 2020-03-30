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

func TestNewNamespace(t *testing.T) {
	actual := NewNamespace("namespace-sample")

	namespaceSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-sample",
		},
	}
	expected := Namespace{namespaceSpec}
	assert.Equal(t, expected, actual)
}

func TestNamespace_Get(t *testing.T) {
	t.Run("should get a namespace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		namespaceInterface := mock_v1.NewMockNamespaceInterface(ctrl)
		coreV1Interface.EXPECT().Namespaces().Return(namespaceInterface)
		namespaceInterface.EXPECT().Get("namespace-sample", metav1.GetOptions{}).Return(&corev1.Namespace{}, nil)
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		ns := NewNamespace("namespace-sample")
		err := ns.Get(kubeClient)

		assert.NoError(t, err)
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		namespaceInterface := mock_v1.NewMockNamespaceInterface(ctrl)
		coreV1Interface.EXPECT().Namespaces().Return(namespaceInterface)
		namespaceInterface.EXPECT().Get("namespace-sample", metav1.GetOptions{}).Return(&corev1.Namespace{}, fmt.Errorf("failed to get namespace"))
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		ns := NewNamespace("namespace-sample")
		err := ns.Get(kubeClient)

		assert.EqualError(t, err, "failed to get namespace")
	})
}

func TestNamespace_Create(t *testing.T) {
	t.Run("should create a namespace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		namespaceInterface := mock_v1.NewMockNamespaceInterface(ctrl)
		coreV1Interface.EXPECT().Namespaces().Return(namespaceInterface)
		namespaceSpec := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "namespace-sample",
			},
		}
		namespaceInterface.EXPECT().Create(namespaceSpec).Return(&corev1.Namespace{}, nil)
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		ns := NewNamespace("namespace-sample")
		err := ns.Create(kubeClient)

		assert.NoError(t, err)
	})

	t.Run("should return error on failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		coreV1Interface := mock_v1.NewMockCoreV1Interface(ctrl)
		namespaceInterface := mock_v1.NewMockNamespaceInterface(ctrl)
		coreV1Interface.EXPECT().Namespaces().Return(namespaceInterface)
		namespaceSpec := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "namespace-sample",
			},
		}
		namespaceInterface.EXPECT().Create(namespaceSpec).Return(&corev1.Namespace{}, fmt.Errorf("failed to create namespace"))
		kubeClient.EXPECT().CoreV1().Return(coreV1Interface)

		ns := NewNamespace("namespace-sample")
		err := ns.Create(kubeClient)

		assert.EqualError(t, err, "failed to create namespace")
	})
}
