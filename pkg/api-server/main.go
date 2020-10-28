package main

import (
	"encoding/json"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/pipeline"
	steprunner "github.com/task-executor/pkg/step-runner"
	"io/ioutil"
	"log"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/api/pipeline/", handlePipeline)
	mux.HandleFunc("/api/tasks/", handleTask)

	http.ListenAndServe(":8080", mux)
}
