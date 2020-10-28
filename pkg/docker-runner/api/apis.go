package api

type ContainerConfig struct {
	Image   string
	Command []string
	Volumes []VolumeMount
	Env     []string
}

type VolumeMount struct {
	Source string
	Target string
}
