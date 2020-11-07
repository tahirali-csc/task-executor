package kube

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/core"
	"github.com/task-executor/pkg/scheduler"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeScheduler struct {
	config *Config
	client *kubernetes.Clientset
}

func (k KubeScheduler) Schedule(ctx context.Context, stage *core.Stage, initContainers []core.InitContainer) error {
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "wwwwww",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:  corev1.RestartPolicyNever,
					InitContainers: k.newInitContainers(initContainers),
					Containers: []corev1.Container{
						{
							Name:  "pipeline",
							Image: stage.Image,
							//TODO: Will replace
							ImagePullPolicy: "Never",
							Command:         stage.Command,
							Args:            stage.Args,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "www-data",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: k.newMountVolumes(stage.Volume),
				},
			},
		},
	}

	job, err := k.client.BatchV1().Jobs(k.config.Namespace).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		//log.WithError(err).Errorln("kubernetes: cannot create job")
	} else {
		//log.Debugf("kubernetes: successfully created job")
	}

	return err
}

func NewKubeScheduler(conf *Config) (scheduler.Scheduler, error) {
	config, err := clientcmd.BuildConfigFromFlags(conf.ConfigURL, conf.ConfigPath)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KubeScheduler{client: client, config: conf}, nil
}

var _ scheduler.Scheduler = (*KubeScheduler)(nil)

func (k KubeScheduler) newInitContainers(containers []core.InitContainer) []corev1.Container {
	var initCont []corev1.Container

	for _, cont := range containers {
		c := corev1.Container{
			Name:  cont.Name,
			Image: cont.Image,
			//TODO
			ImagePullPolicy: "Never",
			Command:         cont.Command,
		}

		for _, s := range cont.Secrets {
			c.Env = append(c.Env, corev1.EnvVar{
				Name: s.Name,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: s.From.Name,
						},
						Key: s.From.Key,
					},
				},
			})
		}

		for _, v := range cont.Volume {
			c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
				Name:      v.Name,
				MountPath: v.MountPath,
			})
		}

		initCont = append(initCont, c)
	}
	return initCont
}

func (k KubeScheduler) newMountVolumes(volume []core.InitVolume) []corev1.Volume {
	var vols []corev1.Volume
	for _, v := range volume {
		vols = append(vols, corev1.Volume{
			Name:         v.Name,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})
	}
	return vols
}
