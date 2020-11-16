package kube

import (
	"github.com/task-executor/pkg/engine"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// helper function returns a kubernetes namespace
// for the given specification.
func toNamespace(spec *engine.Spec) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   spec.Metadata.Namespace,
			//Labels: spec.Metadata.Labels,
		},
	}
}