package init

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/cmd/helm/installer"
	"time"
)

// TillerInstaller install/upgrade tiller
type TillerInstaller interface {
	Install(opts *installer.Options) error
	Upgrade(opts *installer.Options) error
	Wait(options *installer.Options) error
}

type defaultTillerInstaller struct {
	kubeClient            kubernetes.Interface
	TillerWaitTimeout     int
	TillerWaitPollTimeout int
}

// Install will install the tiller in the given namespace
func (d defaultTillerInstaller) Install(opts *installer.Options) error {
	return installer.Install(d.kubeClient, opts)
}

// Upgrade will upgrade the tiller in the given namespace
func (d defaultTillerInstaller) Upgrade(opts *installer.Options) error {
	return installer.Upgrade(d.kubeClient, opts)
}

// Wait will wait the tiller deployment to have
// availableReplicas greater than zero in the given namespace
// Timeout for wait is 60 seconds with the poll period of 10 seconds
func (d defaultTillerInstaller) Wait(options *installer.Options) error {
	return checkAndWait(d.kubeClient, *options, d.TillerWaitTimeout, d.TillerWaitTimeout, d.TillerWaitPollTimeout)
}

func checkAndWait(kubeClient kubernetes.Interface, options installer.Options, timeout, waitingTime, pollTimeout int) error {
	if waitingTime > 0 {
		deployment, err := kubeClient.AppsV1().Deployments(options.Namespace).Get("tiller-deploy", v1.GetOptions{})
		if err != nil {
			return err
		}
		if deployment.Status.AvailableReplicas > 0 {
			return nil
		}
		time.Sleep(time.Duration(pollTimeout) * time.Second)
		return checkAndWait(kubeClient, options, timeout, waitingTime-pollTimeout, pollTimeout)
	}
	return fmt.Errorf("tiller are not ready in %d seconds in %s", timeout, options.Namespace)
}

// NewTillerInstaller creates a defaultTillerInstaller
func NewTillerInstaller(kubeClient kubernetes.Interface, tillerWaitTimeout int, tillerWaitPollTimeout int) TillerInstaller {
	return defaultTillerInstaller{
		kubeClient:            kubeClient,
		TillerWaitTimeout:     tillerWaitTimeout,
		TillerWaitPollTimeout: tillerWaitPollTimeout,
	}
}
