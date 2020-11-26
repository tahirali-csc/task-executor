package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api-server/services"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
	"github.com/task-executor/pkg/api-server/trigger"
	"io/ioutil"
	"net/http"
	"strconv"
)

var buildService = services.NewBuildService()

func HandleBuild(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		createBuild(r, w)
	} else if r.Method == http.MethodGet {
		findBuild(r, w)
	}
}

func HandleBuildStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		updateBuildStatus(r, w)
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
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to read body", 500)
		return
	}

	namespace := r.URL.Query()["namespace"][0]
	repoName := r.URL.Query()["repoName"][0]
	//TODO::
	//branch := r.URL.Query()["branch"]

	repoService := services.RepoService{}
	repo, err := repoService.FindByNamespaceAndName(namespace, repoName)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to find repo", 500)
		return
	}
	log.Println(repo)

	//TODO: Find who triggered the build??

	//TODO: Branch??

	//build := &api.Build{}
	//err = json.Unmarshal(data, build)
	//if err != nil {
	//	log.Error(err)
	//	http.Error(w, "Unable to serialize", 500)
	//	return
	//}

	buildTrigger, err := trigger.NewBuildTrigger()
	build, err := buildTrigger.Trigger(repo)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to save", 500)
		return
	}

	dat, _ := json.Marshal(build)
	w.Write(dat)
}

func updateBuildStatus(r *http.Request, w http.ResponseWriter) {
	buildIdVar := mux.Vars(r)["id"]
	status := mux.Vars(r)["status"]
	statusId := staticdata.BuildStatusList[status].Id
	buildId, _ := strconv.ParseInt(buildIdVar, 10, 64)
	err := buildService.UpdateStatus(buildId, statusId)
	if err != nil {
		log.Println(err)
		log.Println("Can not update status")
		return
	}
}
