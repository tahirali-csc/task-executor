package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/services"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
)

var stepService = services.NewStepService()

//var stepRunner = runner.NewRunner()

func HandleStep(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createStep(r, w)
	} else if r.Method == http.MethodGet {
		findStep(r, w)
	}
}

//func HandleStepStatus(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		findStatus(w, r)
//	}
//	//} else if r.Method == http.MethodPost {
//	//	updateStatus(w, r)
//	//}
//}

func HandleStepStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		updateStatus(w, r)
	} else if r.Method == http.MethodGet {
		findStatus(w, r)
	}
}

func updateStatus(w http.ResponseWriter, r *http.Request) {
	stepIdVar := mux.Vars(r)["id"]
	status := mux.Vars(r)["status"]
	statusId := staticdata.BuildStatusList[status].Id
	stepId, _ := strconv.ParseInt(stepIdVar, 10, 64)
	err := stepService.UpdateStatus(stepId, statusId)
	if err != nil {
		log.Println(err)
		log.Println("Can not update status")
		return
	}
}

func findStatus(w http.ResponseWriter, r *http.Request) {
	stepIdVar := mux.Vars(r)["id"]
	stepId, _ := strconv.ParseInt(stepIdVar, 10, 64)
	status, err := buildService.GetStatus(stepId)
	if err != nil {
		log.Println(err)
		log.Println("Can not get status")
		return
	}

	dat, err := json.Marshal(status)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to parse data", 500)
		return
	}

	_, err = w.Write(dat)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to convert", 500)
		return
	}
	return
}

type stepExec struct {
	Name     string
	Image    string
	Cmd      []string
	Args     []string
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

	now := time.Now()
	res, err := stepService.Create(&api.Step{
		Build: api.Build{
			Id: step.BuildId,
		},
		Name: step.Name,
		Status: api.BuildStatus{
			Id: buildStatus.Id,
		},
		StartTs: &now,
	})
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to parse data", 500)
		return
	}

	//staticdata.StepChannel <- res
	//staticdata.EventBroker.Publish(step.BuildId, res)

	dat, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to parse data", 500)
		return
	}
	w.Write(dat)

	//rs := &core.StepRun{
	//	Name:     step.Name,
	//	BuildId:  step.BuildId,
	//	Args:     step.Args,
	//	Cmd:      step.Cmd,
	//	Image:    step.Image,
	//	Memory:   step.Memory,
	//	CpuLimit: step.CpuLimit,
	//	Step:     res,
	//}

	//go stepRunner.Run(rs)
}

func findStep(r *http.Request, w http.ResponseWriter) {
	builds, err := stepService.Filter(r.URL.Query())
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to retrieve data", 500)
		return
	}

	data, err := json.Marshal(builds)
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

func HandleStepLogsUpload(res http.ResponseWriter, req *http.Request) {
	stepIdVar := mux.Vars(req)["id"]
	stepId, _ := strconv.ParseInt(stepIdVar, 10, 64)

	if req.Method == http.MethodPost {
		staticdata.LogStore.Upload(context.Background(), stepId, req.Body)
	} else if req.Method == http.MethodGet {
		logs, err := staticdata.LogStore.Find(context.Background(), stepId)
		if err != nil {
			log.Error(err)
			return
		}
		dat, err := ioutil.ReadAll(logs)
		if err != nil {
			log.Error(err)
			return
		}
		res.Write(dat)
	}
}
