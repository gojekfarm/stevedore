package resources_test

import (
	"fmt"
	"github.com/gojek/stevedore/pkg/internal/mocks"
	mockK8s "github.com/gojek/stevedore/pkg/internal/mocks/client-go/kubernetes"
	k8sresources "github.com/gojek/stevedore/pkg/stevedore/k8s/resources"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestCreateIfNotExists(t *testing.T) {
	t.Run("should create resource if does not exist already", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		resource := mocks.NewMockResource(ctrl)
		notFoundErr := &apierrors.StatusError{ErrStatus: metav1.Status{Code: 404}}
		resource.EXPECT().Get(kubeClient).Return(notFoundErr)
		resource.EXPECT().Create(kubeClient).Return(nil)

		resourceCreator := k8sresources.NewResourceCreator(kubeClient)
		err := resourceCreator.CreateIfNotExists(resource)

		assert.NoError(t, err)
	})

	t.Run("should return error if resource creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		resource := mocks.NewMockResource(ctrl)
		notFoundErr := &apierrors.StatusError{ErrStatus: metav1.Status{Code: 404}}
		resource.EXPECT().Get(kubeClient).Return(notFoundErr)
		resource.EXPECT().Create(kubeClient).Return(fmt.Errorf("some error"))

		resourceCreator := k8sresources.NewResourceCreator(kubeClient)
		err := resourceCreator.CreateIfNotExists(resource)

		assert.EqualError(t, err, "unable to create resource: some error")
	})

	t.Run("should not create resource if it exists already", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		resource := mocks.NewMockResource(ctrl)
		resource.EXPECT().Get(kubeClient).Return(nil)

		resourceCreator := k8sresources.NewResourceCreator(kubeClient)
		err := resourceCreator.CreateIfNotExists(resource)

		assert.NoError(t, err)
	})

	t.Run("should return error if get resource fails and does not return a APIStatus error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		resource := mocks.NewMockResource(ctrl)
		resource.EXPECT().Get(kubeClient).Return(fmt.Errorf("some error"))

		resourceCreator := k8sresources.NewResourceCreator(kubeClient)
		err := resourceCreator.CreateIfNotExists(resource)

		assert.EqualError(t, err, "unable to get resource: some error")
	})

	t.Run("should return error if get namespace fails with non-404 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		kubeClient := mockK8s.NewMockInterface(ctrl)
		resource := mocks.NewMockResource(ctrl)
		internalErr := &apierrors.StatusError{ErrStatus: metav1.Status{Code: 500, Message: "internal server error"}}
		resource.EXPECT().Get(kubeClient).Return(internalErr)

		resourceCreator := k8sresources.NewResourceCreator(kubeClient)
		err := resourceCreator.CreateIfNotExists(resource)

		assert.EqualError(t, err, "unable to get resource: internal server error")
	})
}
