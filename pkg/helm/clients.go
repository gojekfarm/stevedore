package helm

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/kube"
)

// Clients represents a collection of helm clients mapped to their respective tiller namespaces
type Clients map[string]Client

// Close closes all the helm clients
func (clients Clients) Close() {
	for _, client := range clients {
		client.Close()
	}
}

// NewHelmClients creates helm clients for given tiller namespaces
func NewHelmClients(tillerNamespaces []string, kubecontext, kubeconfig string) (Clients, error) {
	clientConfig, err := kube.GetConfig(kubecontext, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error when creating kubernetes client config due to %v", err)
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes config due to %v", err)
	}
	return NewHelmClientsUsingKubeClient(tillerNamespaces, client, clientConfig)
}

// NewHelmClientsUsingKubeClient creates helm clients for given tiller namespaces
func NewHelmClientsUsingKubeClient(tillerNamespaces []string, kubeClient kubernetes.Interface, kubeClientConfig *rest.Config) (Clients, error) {
	helmClients := make(Clients)
	for _, namespace := range tillerNamespaces {
		if _, present := helmClients[namespace]; !present {
			client, err := NewHelmClient(namespace, kubeClient, kubeClientConfig)
			if err != nil {
				helmClients.Close()
				return nil, err
			}
			helmClients[namespace] = client
		}
	}

	return helmClients, nil
}
