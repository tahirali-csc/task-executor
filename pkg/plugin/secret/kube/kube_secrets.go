package kube

import (
	"context"
	"fmt"
	"github.com/task-executor/pkg/plugin/secret"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesSecret struct {
	client          kubernetes.Interface
	namespacePrefix string
}

func (k8 KubernetesSecret) Get(secretPath string) (*secret.Secret, error) {
	parts := strings.Split(secretPath, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unable to split kubernetes secret path into [namespace]/[secret]: %s", secretPath)
	}

	var namespace = parts[0]
	var secretName = parts[1]

	secret, found, err := k8.findSecret(namespace, secretName)
	if err != nil {
		return nil, err
	}

	if found {
		return k8.getValueFromSecret(secret)
	}

	return nil, err
}

func (k8 KubernetesSecret) getValueFromSecret(k8Secret *v1.Secret) (*secret.Secret, error) {
	sec := secret.Secret{
		Metadata: make(map[string]interface{}),
	}
	sec.Name = k8Secret.Name
	sec.Metadata["type"] = string(k8Secret.Type)

	data := map[string]interface{}{}
	for k, v := range k8Secret.Data {
		data[k] = v
	}
	sec.Data = data

	return &sec, nil
}

func (k8 KubernetesSecret) findSecret(namespace, name string) (*v1.Secret, bool, error) {
	var secret *v1.Secret
	var err error

	secret, err = k8.client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	} else {
		return secret, true, err
	}
}
