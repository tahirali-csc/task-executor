package controllers

import (
	"encoding/json"
	"github.com/task-executor/pkg/api-server/services"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
)

var buildService = services.NewBuildService()

func HandleBuild(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		createBuild(r, w)
	} else if r.Method == http.MethodGet {
		findBuild(r, w)
	}
}

func findBuild(r *http.Request, w http.ResponseWriter) {
	builds, err := buildService.Filter(r.URL.Query())
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

func createBuild(r *http.Request, w http.ResponseWriter) {
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

	res, err := buildService.Create(build)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to save", 500)
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
