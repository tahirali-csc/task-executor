package api

type PipelineConfig struct {
}

type RunConfig struct {
	Image   string
	Command []string
	//TODO: Will review. Object/Property
	BuildId int64
}
