package kube

import (
	"github.com/task-executor/pkg/plugin/secret"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesManager struct {
	InClusterConfig bool   `long:"in-cluster" description:"Enables the in-cluster client."`
	ConfigPath      string `long:"config-path" description:"Path to Kubernetes config when running ATC outside Kubernetes."`
	NamespacePrefix string `long:"namespace-prefix" default:"concourse-" description:"Prefix to use for Kubernetes namespaces under which secrets will be looked up."`
}

func (manager KubernetesManager) buildConfig() (*rest.Config, error) {
	if manager.InClusterConfig {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags("", manager.ConfigPath)
}

func (manager KubernetesManager) NewSecretsFactory() (secret.SecretsFactory, error) {
	config, err := manager.buildConfig()
	if err != nil {
		return nil, err
	}

	config.QPS = 100
	config.Burst = 100

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return NewKubernetesSecretFactory(clientset, manager.NamespacePrefix), nil
}
