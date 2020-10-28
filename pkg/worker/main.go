package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/utils"
	"io/ioutil"
	"net/http"
)

func handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to read body", 500)
			return
		}

		runConfig := &api.RunConfig{}
		err = json.Unmarshal(body, runConfig)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to parse body", 500)
			return
		}

		client := http.Client{}
		client.Post("http://localhost:8080/api/tasks/", "application/json", bytes.NewReader(body))
	}
}

func main() {
	utils.InitLogs(log.DebugLevel)

	log.Println("Starting Worker")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks/", handleTask)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Error(err)
	}
}
