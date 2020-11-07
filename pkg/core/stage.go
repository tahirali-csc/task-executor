package core

//type Container struct {
//	Image           string
//	ImagePullPolicy string
//	Command         []string
//	Args            []string
//}

type Stage struct {
	//ID        int64             `json:"id"`
	//BuildID   int64             `json:"build_id"`build_id
	Image           string
	ImagePullPolicy string
	LimitMemory     int
	LimitCompute    int
	RequestMemory   int
	RequestCompute  int
	Command         []string
	Args            []string
	//InitContainers  []string
	Volume []InitVolume
}

type InitContainer struct {
	Name            string
	Image           string
	ImagePullPolicy string
	Command         []string
	Secrets         []SecretSource
	Volume          []InitVolume
}

type SecretFromRef struct {
	Name string
	Key  string
}
type SecretSource struct {
	Name string
	From SecretFromRef
}
type InitVolume struct {
	Name      string
	MountPath string
}
