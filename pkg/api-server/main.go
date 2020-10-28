package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/config"
	"github.com/task-executor/pkg/api-server/controllers"
	"github.com/task-executor/pkg/pipeline"
	steprunner "github.com/task-executor/pkg/step-runner"
	"github.com/task-executor/pkg/utils"
	"io/ioutil"

	"net/http"
)

func handlePipeline(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	http.Error(w, "Unable to parse body", 500)
		//	return
		//}
		//
		//runConfig := &api.RunConfig{}
		//err = json.Unmarshal(body, runConfig)
		//if err != nil {
		//	http.Error(w, "Unable to parse body", 500)
		//	return
		//}

		pipeline.Run()
	}
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to parse body", 500)
			return
		}

		runConfig := &api.RunConfig{}
		err = json.Unmarshal(body, runConfig)
		if err != nil {
			http.Error(w, "Unable to parse body", 500)
			return
		}

		log.Println("The request object:::", runConfig)

		steprunner.Run(runConfig)
	}
}

func main() {
	utils.InitLogs(log.DebugLevel)

	appConfig, err := config.Load()
	if err != nil {
		log.Error(err)
		return
	}

	mux := http.NewServeMux()

	//TODO::
	mux.HandleFunc("/api/pipeline/", handlePipeline)
	mux.HandleFunc("/api/builds/", controllers.HandleBuild)
	mux.HandleFunc("/api/tasks/", handleTask)

	err = http.ListenAndServe(fmt.Sprintf(":%s", appConfig.Server.Port), mux)
	if err != nil {
		log.Error(err)
	}
}
