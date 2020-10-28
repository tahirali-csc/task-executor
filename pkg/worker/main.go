package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
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
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
	log.SetLevel(log.DebugLevel)
	log.Debugf("fksdjfsd")

	log.Println("Starting Worker")
	mux := http.NewServeMux()

	mux.HandleFunc("/api/tasks/", handleTask)

	http.ListenAndServe(":8081", mux)
}
