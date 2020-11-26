package trigger

import (
	"context"
	"fmt"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/services"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
	"github.com/task-executor/pkg/core"
	"github.com/task-executor/pkg/plugin/secret"
	kubesecret "github.com/task-executor/pkg/plugin/secret/kube"
	"github.com/task-executor/pkg/scheduler"
	"github.com/task-executor/pkg/scheduler/kube"
	"os"
	"path"
)

type BuildTrigger struct {
	Scheduler    scheduler.Scheduler
	BuildService services.BuildService
}

func NewBuildTrigger() (*BuildTrigger, error) {
	buildTrigger := &BuildTrigger{}

	sch, err := provideKubernetesScheduler()
	if err != nil {
		return nil, err
	}

	buildTrigger.Scheduler = sch
	return buildTrigger, nil
}

func (trigger *BuildTrigger) Trigger(repo *api.Repo) (*api.Build, error) {
	//Start in Pending State
	pendingStatus := staticdata.BuildStatusList[api.PendingBuildStatus]

	pendingBuild := api.Build{
		//TODO::
		RepoBranch: api.RepoBranch{
			Id:   22,
			Repo: *repo,
			Name: "master",
		},
		Status: pendingStatus,
	}

	res, err := trigger.BuildService.Create(&pendingBuild)
	if err != nil {
		return nil, err
	}

	// Fetch Secret
	secret, err := injectSecret(repo)
	if err != nil {
		return nil, err
	}

	secretObj, err := secret.Get(repo.SecretName)
	if err != nil {
		return nil, err
	}

	var cmd []string
	cmd = append(cmd, "./repo-cloner")
	cmd = append(cmd, "--clone-dir=/data")

	if repo.AuthType.Name == api.BasicAuthType {
		cmd = append(cmd, "--repo="+repo.HttpUrl)
		cmd = append(cmd, "--secret-type=basic-auth")
	} else if repo.AuthType.Name == api.SSHAuthType {
		cmd = append(cmd, "--repo="+repo.SSHUrl)
		cmd = append(cmd, "--secret-type=ssh-auth")
	}

	scmCloneContainer := core.InitContainer{
		Name: "scm-clone",
		//TODO: review the image tag & pull policy
		Image:           "repo-cloner:latest",
		ImagePullPolicy: "Never",
		Command:         cmd,
		Secrets: []core.SecretSource{
			{
				Name: "USER",
				From: core.SecretFromRef{
					Name: secretObj.Name,
					Key:  "username",
				},
			},
			{
				Name: "PASSWORD",
				From: core.SecretFromRef{
					Name: secretObj.Name,
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

	//id := uuid.New()

	trigger.Scheduler.Schedule(context.Background(), &core.Stage{
		Name:            fmt.Sprintf("te-build-%d", res.Id),
		Image:           "golang:1.14",
		ImagePullPolicy: "Never",
		LimitMemory:     0,
		LimitCompute:    0,
		RequestMemory:   0,
		RequestCompute:  0,
		Command:         []string{"/bin/sh", "-c", "cd /data/ci && go get && go run ci.go"},
		Volume: []core.InitVolume{
			{
				Name:      "www-data",
				MountPath: "/data",
			},
		},
		BuildId: res.Id,
	}, []core.InitContainer{scmCloneContainer})

	return res, nil
}

//func scheduleBuild() {
//	homeDir, err := os.UserHomeDir()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	kubeConfig := path.Join(homeDir, ".kube", "config")
//
//	kubeSch, err := kube.NewKubeScheduler(&kube.Config{
//		ConfigURL:      "",
//		ConfigPath:     kubeConfig,
//		Namespace:      "default",
//		ServiceAccount: "default",
//		//TODO::
//	})
//
//	km := &kubesecret.KubernetesManager{
//		ConfigPath: kubeConfig,
//	}
//	secretsFactory, _ := km.NewSecretsFactory()
//	s := secretsFactory.NewSecrets()
//
//	var cmd []string
//	sec, err := s.Get("default/secret-basic-auth")
//
//	cmd = append(cmd, "./repo-cloner")
//	cmd = append(cmd, "--repo=https://github.com/tahirali-csc/hello-app.git")
//	cmd = append(cmd, "--clone-dir=/data")
//	if sec.Metadata["type"] == "kubernetes.io/basic-auth" {
//		cmd = append(cmd, "--secret-type=basic-auth")
//	} else if sec.Metadata["type"] == "kubernetes.io/ssh-auth" {
//		cmd = append(cmd, "--secret-type=ssh-auth")
//	}
//
//	scmCloneContainer := core.InitContainer{
//		Name:            "scm-clone",
//		Image:           "repo-cloner:latest",
//		ImagePullPolicy: "Never",
//		Command:         cmd,
//		Secrets: []core.SecretSource{
//			{
//				Name: "USER",
//				From: core.SecretFromRef{
//					Name: "secret-basic-auth",
//					Key:  "username",
//				},
//			},
//			{
//				Name: "PASSWORD",
//				From: core.SecretFromRef{
//					Name: "secret-basic-auth",
//					Key:  "password",
//				},
//			},
//		},
//		Volume: []core.InitVolume{
//			{
//				Name:      "www-data",
//				MountPath: "/data",
//			},
//		},
//	}
//
//	kubeSch.Schedule(context.Background(), &core.Stage{
//		Image:           "golang:1.14",
//		ImagePullPolicy: "Never",
//		LimitMemory:     0,
//		LimitCompute:    0,
//		RequestMemory:   0,
//		RequestCompute:  0,
//		Command:         []string{"/bin/sh", "-c", "ls -al /data && cd /data/ci && go run ci.go"},
//		Volume: []core.InitVolume{
//			{
//				Name:      "www-data",
//				MountPath: "/data",
//			},
//		},
//	}, []core.InitContainer{scmCloneContainer})
//}

func provideKubernetesScheduler() (scheduler.Scheduler, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeConfig := path.Join(homeDir, ".kube", "config")

	kubeSch, err := kube.NewKubeScheduler(&kube.Config{
		ConfigURL:  "",
		ConfigPath: kubeConfig,
		//TODO: Will review
		Namespace:      "default",
		ServiceAccount: "default",
		//TODO: Add more properties
		//TODO: Externalize
		HostURL: "http://192.168.64.1:8080",
	})

	return kubeSch, err
}
func injectSecret(repo *api.Repo) (secret.Secrets, error) {
	if repo.SecretType.Name == api.KubernetesSecretType {
		//TODO: review duplicate!!
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeConfig := path.Join(homeDir, ".kube", "config")

		kubeSecManager := &kubesecret.KubernetesManager{
			ConfigPath: kubeConfig,
		}
		secretsFactory, err := kubeSecManager.NewSecretsFactory()
		if err != nil {
			return nil, err
		}

		kubeSecret := secretsFactory.NewSecrets()
		return kubeSecret, nil
	}
	return nil, nil
}
