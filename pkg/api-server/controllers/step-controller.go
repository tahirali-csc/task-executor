package controllers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/services"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
	"io/ioutil"
	"net/http"
)

var stepService = services.NewStepService()

func HandleStep(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		createStep(r, w)
	} else if r.Method == http.MethodGet {

	}
}

type stepExec struct {
	Name     string
	Image    string
	Cmd      []string
	CpuLimit int
	Memory   int
	BuildId  int64
}

func createStep(r *http.Request, w http.ResponseWriter) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to read body", 500)
		return
	}

	step := &stepExec{}
	err = json.Unmarshal(data, step)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to parse data", 500)
		return
	}
	log.Println("Step:::", step)

	buildStatus := staticdata.BuildStatusList[api.PendingBuildStatus]

	res, err := stepService.Create(&api.Step{
		Build: api.Build{
			Id: step.BuildId,
		},
		Name: step.Name,
		Status: api.BuildStatus{
			Id: buildStatus.Id,
		},
	})
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to parse data", 500)
		return
	}

	data, err = json.Marshal(res)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to convert", 500)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to convert", 500)
		return
	}
}
