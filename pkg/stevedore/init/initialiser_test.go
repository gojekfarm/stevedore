package init

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/pkg/helm"
	"github.com/gojek/stevedore/pkg/internal/mocks"
	k8sresources "github.com/gojek/stevedore/pkg/stevedore/k8s/resources"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/cmd/helm/installer"
)

func TestInit(t *testing.T) {
	helm.SetHelmVersion("v2.16.3")
	helm.SetBuildMetadata("")
	t.Run("should give empty response for empty requests", func(t *testing.T) {
		initialiser := DefaultInitialiser{}

		responses, err := initialiser.Init([]Request{})

		assert.Nil(t, err)
		assert.Empty(t, responses)
	})

	t.Run("should initialize successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		request := NewInitRequest("default", "tiller", false)
		anotherRequest := NewInitRequest("app", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("app")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "app")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "app", "tiller")).Return(nil)

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    anotherRequest.namespace,
			ServiceAccount:               anotherRequest.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    anotherRequest.namespace,
			ServiceAccount:               anotherRequest.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{request, anotherRequest}

		responses, err := initialiser.Init(requests)

		assert.Nil(t, err)
		assert.Len(t, responses, 2)
		expectedResponses := Responses{
			Response{Message: "Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.", Namespace: "default"},
			Response{Message: "Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.", Namespace: "app"},
		}
		assert.Equal(t, expectedResponses, responses)
	})

	t.Run("should initialize successfully for privileged namespace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		request := NewInitRequest("kube-system", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("kube-system")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "kube-system")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminClusterRoleBinding("tiller-cluster-admin", "kube-system", "tiller")).Return(nil)

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{request}

		responses, err := initialiser.Init(requests)

		assert.Nil(t, err)
		assert.Len(t, responses, 1)
		expectedResponses := Responses{
			Response{Message: "Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.", Namespace: "kube-system"},
		}
		assert.Equal(t, expectedResponses, responses)
	})

	t.Run("should return error when failure is not due to tiller being installed already", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		request := NewInitRequest("default", "tiller", false)
		anotherRequest := NewInitRequest("app", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(fmt.Errorf("some error"))

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{request, anotherRequest}

		responses, err := initialiser.Init(requests)

		assert.EqualError(t, err, "error installing: some error")
		assert.Empty(t, responses)
	})

	t.Run("should upgrade if failure is due to tiller being present already", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		request := NewInitRequest("default", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		alreadyPresentError := &apierrors.StatusError{ErrStatus: metav1.Status{Reason: metav1.StatusReasonAlreadyExists}}
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(alreadyPresentError)
		tillerInstaller.EXPECT().Upgrade(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{request}

		responses, err := initialiser.Init(requests)

		assert.Nil(t, err)
		assert.Len(t, responses, 1)
		expectedResponse := Response{"default", "Tiller (the Helm server-side component) has been upgraded to the current version."}
		assert.Equal(t, expectedResponse, responses[0])
	})

	t.Run("should force upgrade when init request has force upgrade enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		request := NewInitRequest("default", "tiller", true)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		alreadyPresentError := &apierrors.StatusError{ErrStatus: metav1.Status{Reason: metav1.StatusReasonAlreadyExists}}
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
			ForceUpgrade:                 true,
		}).Return(alreadyPresentError)
		tillerInstaller.EXPECT().Upgrade(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
			ForceUpgrade:                 true,
		}).Return(nil)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
			ForceUpgrade:                 true,
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{request}

		responses, err := initialiser.Init(requests)

		assert.Nil(t, err)
		assert.Len(t, responses, 1)
		expectedResponse := Response{"default", "Tiller (the Helm server-side component) has been upgraded to the current version."}
		assert.Equal(t, expectedResponse, responses[0])
	})

	t.Run("should return error if tiller upgrade fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		request := NewInitRequest("default", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		alreadyPresentError := &apierrors.StatusError{ErrStatus: metav1.Status{Reason: metav1.StatusReasonAlreadyExists}}
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(alreadyPresentError)
		tillerInstaller.EXPECT().Upgrade(&installer.Options{
			Namespace:                    request.namespace,
			ServiceAccount:               request.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(fmt.Errorf("some error"))

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{request}

		responses, err := initialiser.Init(requests)

		assert.Error(t, err, "error when upgrading: some error")
		assert.Empty(t, responses)
	})

	t.Run("should return error even if one namespace creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		requestOne := NewInitRequest("default", "tiller", false)
		requestTwo := NewInitRequest("app", "tiller", false)
		requestThree := NewInitRequest("core", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("app")).Return(fmt.Errorf("namespace creation error"))

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    requestOne.namespace,
			ServiceAccount:               requestOne.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    requestOne.namespace,
			ServiceAccount:               requestOne.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{requestOne, requestTwo, requestThree}

		responses, err := initialiser.Init(requests)

		assert.EqualError(t, err, "failed to create namespace app: namespace creation error")
		assert.Empty(t, responses)
	})

	t.Run("should return error even if one serviceaccount creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		requestOne := NewInitRequest("default", "tiller", false)
		requestTwo := NewInitRequest("app", "tiller", false)
		requestThree := NewInitRequest("core", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("app")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "app")).Return(fmt.Errorf("serviceaccount creation error"))

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    requestOne.namespace,
			ServiceAccount:               requestOne.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    requestOne.namespace,
			ServiceAccount:               requestOne.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{requestOne, requestTwo, requestThree}

		responses, err := initialiser.Init(requests)

		assert.EqualError(t, err, "failed to create serviceaccount tiller in namespace app: serviceaccount creation error")
		assert.Empty(t, responses)
	})

	t.Run("should return error even if one rolebinding creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		requestOne := NewInitRequest("default", "tiller", false)
		requestTwo := NewInitRequest("app", "tiller", false)
		requestThree := NewInitRequest("core", "tiller", false)

		resourceCreator := mocks.NewMockResourceCreator(ctrl)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "default")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "default", "tiller")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewNamespace("app")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewServiceAccount("tiller", "app")).Return(nil)
		resourceCreator.EXPECT().CreateIfNotExists(k8sresources.NewClusterAdminRoleBinding("tiller-cluster-admin", "app", "tiller")).Return(fmt.Errorf("rolebinding creation error"))

		tillerInstaller := mocks.NewMockTillerInstaller(ctrl)
		tillerInstaller.EXPECT().Install(&installer.Options{
			Namespace:                    requestOne.namespace,
			ServiceAccount:               requestOne.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)
		tillerInstaller.EXPECT().Wait(&installer.Options{
			Namespace:                    requestOne.namespace,
			ServiceAccount:               requestOne.serviceAccount,
			AutoMountServiceAccountToken: true,
			ImageSpec:                    "gcr.io/kubernetes-helm/tiller:v2.16.3",
		}).Return(nil)

		initialiser := NewDefaultInitialiser(tillerInstaller, resourceCreator)
		requests := []Request{requestOne, requestTwo, requestThree}

		responses, err := initialiser.Init(requests)

		assert.EqualError(t, err, "failed to create rolebinding tiller-cluster-admin in namespace app: rolebinding creation error")
		assert.Empty(t, responses)
	})
}
