package build

import (
	log "github.com/sirupsen/logrus"
	apiclient "github.com/task-executor/pkg/api-client"
	"time"
)

type Controller struct {
}

func NewBuildController() *Controller {
	return &Controller{}
}

func (bc *Controller) Start() {
	buildFilter := apiclient.NewBuildFilterBuilder()
	buildFilter.WithStatus(3)
	buildFilter.ByCreatedTs(false)

	buildClient := apiclient.NewBuildClient(apiclient.Config{
		Host: "localhost:8080",
	})

	for {
		builds, err := buildClient.GetBuilds(buildFilter)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, b := range builds {
			log.Println(b.CreatedTs)
		}

		time.Sleep(time.Second * 2)
		log.Println("")
	}

}
