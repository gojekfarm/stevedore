package init

import (
	"fmt"

	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/version"
)

// Request to represent the request struct for stevedore Init
type Request struct {
	namespace          string
	serviceAccount     string
	forceUpgradeTiller bool
}

func (request Request) getOpts() installer.Options {
	return installer.Options{
		Namespace:                    request.namespace,
		ServiceAccount:               request.serviceAccount,
		AutoMountServiceAccountToken: true,
		ImageSpec:                    fmt.Sprintf("%s:%s", "gcr.io/kubernetes-helm/tiller", version.Version),
		ForceUpgrade:                 request.forceUpgradeTiller,
	}
}

// NewInitRequest to create new Request
func NewInitRequest(namespace, serviceAccount string, forceUpgradeTiller bool) Request {
	return Request{namespace: namespace, serviceAccount: serviceAccount, forceUpgradeTiller: forceUpgradeTiller}
}
