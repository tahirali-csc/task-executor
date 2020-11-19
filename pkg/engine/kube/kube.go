package kube

import (
	"context"
	"github.com/task-executor/pkg/engine"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubeEngine struct {
	client *kubernetes.Clientset
}

// NewFile returns a new Kubernetes engine from a
// Kubernetes configuration file (~/.kube/config).
func NewFile(url, path, node string) (engine.Engine, error) {
	config, err := clientcmd.BuildConfigFromFlags(url, path)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &kubeEngine{client: client}, nil
}

func (e *kubeEngine) Setup(ctx context.Context, spec *engine.Spec) error {
	ns := toNamespace(spec)

	// create the project namespace. all pods and
	// containers are created within the namespace, and
	// are removed when the pipeline execution completes.
	_, err := e.client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// pv := toPersistentVolume(e.node, spec.Metadata.Namespace, spec.Metadata.Namespace, filepath.Join("/tmp", spec.Metadata.Namespace))
	// _, err = e.client.CoreV1().PersistentVolumes().Create(pv)
	// if err != nil {
	// 	return err
	// }

	// pvc := toPersistentVolumeClaim(spec.Metadata.Namespace, spec.Metadata.Namespace)
	// _, err = e.client.CoreV1().PersistentVolumeClaims(spec.Metadata.Namespace).Create(pvc)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (e *kubeEngine) Start(ctx context.Context, spec *engine.Spec) error {
	pod := toPod(spec)
	//if len(step.Docker.Ports) != 0 {
	//	service := toService(spec, step)
	//	_, err := e.client.CoreV1().Services(spec.Metadata.Namespace).Create(service)
	//	if err != nil {
	//		return err
	//	}
	//}

	//if e.node != "" {
	//	pod.Spec.Affinity = &v1.Affinity{
	//		NodeAffinity: &v1.NodeAffinity{
	//			RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
	//				NodeSelectorTerms: []v1.NodeSelectorTerm{{
	//					MatchExpressions: []v1.NodeSelectorRequirement{{
	//						Key:      "kubernetes.io/hostname",
	//						Operator: v1.NodeSelectorOpIn,
	//						Values:   []string{e.node},
	//					}},
	//				}},
	//			},
	//		},
	//	}
	//}

	_, err := e.client.CoreV1().Pods(spec.Metadata.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	return err
}

func (e *kubeEngine) Tail(ctx context.Context, spec *engine.Spec) (io.ReadCloser, error) {
	ns := spec.Metadata.Namespace
	podName := spec.Metadata.UID

	return e.client.CoreV1().Pods(ns).GetLogs(podName, &v1.PodLogOptions{Follow: true}).
		Stream(ctx)

	//up := make(chan bool)
	//
	//var podUpdated = func(old interface{}, new interface{}) {
	//	pod := new.(*v1.Pod)
	//	if pod.Name == podName {
	//		switch pod.Status.Phase {
	//		case v1.PodRunning, v1.PodSucceeded, v1.PodFailed:
	//			up <- true
	//		}
	//	}
	//}
	//si := informers.NewSharedInformerFactory(e.client, 5*time.Minute)
	//si.Core().V1().Pods().Informer().AddEventHandler(
	//	cache.ResourceEventHandlerFuncs{
	//		UpdateFunc: podUpdated,
	//	},
	//)
	//si.Start(wait.NeverStop)
	//
	//select {
	//case <-up:
	//case <-ctx.Done():
	//}
	//
	//opts := &v1.PodLogOptions{
	//	Follow: true,
	//}
	//
	//return e.client.CoreV1().RESTClient().Get().
	//	Namespace(ns).
	//	Name(podName).
	//	Resource("pods").
	//	SubResource("log").
	//	VersionedParams(opts, scheme.ParameterCodec).
	//	Stream(ctx)
}

// helper function returns a kubernetes pod for the
// given step and specification.
func toPod(spec *engine.Spec) *v1.Pod {
	//var volumes []v1.Volume
	//volumes = append(volumes, toVolumes(spec, step)...)
	//volumes = append(volumes, toConfigVolumes(spec, step)...)
	//
	//var mounts []v1.VolumeMount
	//mounts = append(mounts, toVolumeMounts(spec, step)...)
	//mounts = append(mounts, toConfigMounts(spec, step)...)
	//
	//var pullSecrets []v1.LocalObjectReference
	//if len(spec.Docker.Auths) > 0 {
	//	pullSecrets = []v1.LocalObjectReference{{
	//		Name: "docker-auth-config", // TODO move name to a const
	//	}}
	//}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Metadata.UID,
			Namespace: spec.Metadata.Namespace,
			//Labels:    spec.Metadata.Labels,
		},
		Spec: v1.PodSpec{
			//AutomountServiceAccountToken: boolptr(false),
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{{
				Name:  spec.Metadata.UID,
				Image: spec.Image,
				//ImagePullPolicy: toPullPolicy(step.Docker.PullPolicy),
				Command: spec.Command,
				//Args:            spec.Docker.Args,
				//WorkingDir:      step.WorkingDir,
				//SecurityContext: &v1.SecurityContext{
				//	Privileged: &step.Docker.Privileged,
				//},
				//Env:          toEnv(spec, step),
				//VolumeMounts: mounts,
				//Ports:        toPorts(step),
				//Resources:    toResources(step),
			}},
			//ImagePullSecrets: pullSecrets,
			//Volumes:          volumes,
		},
	}
}

// helper function converts the engine pull policy
// to the kubernetes pull policy constant.
//func toPullPolicy(from engine.PullPolicy) v1.PullPolicy {
//	switch from {
//	case engine.PullAlways:
//		return v1.PullAlways
//	case engine.PullNever:
//		return v1.PullNever
//	case engine.PullIfNotExists:
//		return v1.PullIfNotPresent
//	default:
//		return v1.PullIfNotPresent
//	}
//}
