package kube

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/core"
	"github.com/task-executor/pkg/plugin/secret/kube"
	"os"
	"path"
	"testing"
)

func TestJob(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		t.FailNow()
		return
	}

	kubeConfig := path.Join(homeDir, ".kube", "config")

	km := &kube.KubernetesManager{
		ConfigPath: kubeConfig,
	}
	secretsFactory, _ := km.NewSecretsFactory()
	s := secretsFactory.NewSecrets()

	var cmd []string
	sec, err := s.Get("default/secret-basic-auth")

	cmd = append(cmd, "./repo-cloner")
	cmd = append(cmd, "--repo=https://github.com/tahirali-csc/hello-app.git")
	cmd = append(cmd, "--clone-dir=/data")
	if sec.Metadata["type"] == "kubernetes.io/basic-auth" {
		cmd = append(cmd, "--secret-type=basic-auth")
	} else if sec.Metadata["type"] == "kubernetes.io/ssh-auth" {
		cmd = append(cmd, "--secret-type=ssh-auth")
	}

	scmCloneContainer := core.InitContainer{
		Name:            "scm-clone",
		Image:           "repo-cloner:latest",
		ImagePullPolicy: "Never",
		Command:         cmd,
		Secrets: []core.SecretSource{
			{
				Name: "USER",
				From: core.SecretFromRef{
					Name: "secret-basic-auth",
					Key:  "username",
				},
			},
			{
				Name: "PASSWORD",
				From: core.SecretFromRef{
					Name: "secret-basic-auth",
					Key:  "password",
				},
			},
		},
		Volume: []core.InitVolume{
			{
				Name:      "www-data",
				MountPath: "/data",
			},
		},
	}

	kubeSch, err := NewKubeScheduler(&Config{
		ConfigURL:      "",
		ConfigPath:     kubeConfig,
		Namespace:      "default",
		ServiceAccount: "default",
	})

	kubeSch.Schedule(context.Background(), &core.Stage{
		Image:           "golang:1.14",
		ImagePullPolicy: "Never",
		LimitMemory:     0,
		LimitCompute:    0,
		RequestMemory:   0,
		RequestCompute:  0,
		Command:         []string{"/bin/sh", "-c", "ls -al /data && cd /data/ci && go run ci.go"},
		Volume: []core.InitVolume{
			{
				Name:      "www-data",
				MountPath: "/data",
			},
		},
	}, []core.InitContainer{scmCloneContainer})
}
