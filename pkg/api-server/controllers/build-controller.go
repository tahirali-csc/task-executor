package controllers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/services"
	"io/ioutil"
	"net/http"
)

func HandleBuild(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to read body", 500)
			return
		}

		build := &api.Build{}
		err = json.Unmarshal(data, build)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to serialize", 500)
			return
		}

		buildService := services.NewBuildService()
		res, err := buildService.Create(build)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to save", 500)
			return
		}

		data, err = json.Marshal(res)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to save", 500)
			return
		}

		_, err = w.Write(data)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to convert", 500)
			return
		}
	}
}
