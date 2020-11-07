package kube

import (
	"github.com/task-executor/pkg/plugin/secret"
	"k8s.io/client-go/kubernetes"
)

type kubernetesFactory struct {
	client          kubernetes.Interface
	namespacePrefix string
}

func NewKubernetesSecretFactory(client kubernetes.Interface, namespacePrefix string) *kubernetesFactory {
	return &kubernetesFactory{
		client:          client,
		namespacePrefix: namespacePrefix,
	}
}

func (factory *kubernetesFactory) NewSecrets() secret.Secrets {
	return &KubernetesSecret{
		client:          factory.client,
		namespacePrefix: factory.namespacePrefix,
	}
}
