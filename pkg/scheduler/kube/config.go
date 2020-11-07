package kube

type Config struct {
	ConfigURL      string
	ConfigPath     string
	Namespace      string
	ServiceAccount string
	//CloneContainer CloneContainer
}

//type CloneContainer struct {
//	Image           string
//	ImagePullPolicy v1.PullPolicy
//	Repo            string
//
//}
