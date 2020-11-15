package main

import (
	"context"
	"encoding/json"
	"fmt"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/config"
	"github.com/task-executor/pkg/api-server/controllers"
	"github.com/task-executor/pkg/api-server/dbstore"
	"github.com/task-executor/pkg/pipeline"
	"github.com/task-executor/pkg/scm/driver/github"
	steprunner "github.com/task-executor/pkg/step-runner"
	"github.com/task-executor/pkg/utils"

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

func registerShutdown(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		//Close DB
		<-c
		err := dbstore.Close()
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Shutting database connection")

		//Shutdown webserver
		err = server.Shutdown(context.Background())
		time.Sleep(time.Second * 2)
		if err != nil {
			log.Error(err)
			return
		}
	}()
}


func main() {
	utils.InitLogs(log.DebugLevel)

	config, err := config.Load()
	if err != nil {
		log.Error(err)
		return
	}

	err = dbstore.Init(config)
	if err != nil {
		log.Error(err)
		return
	}

	staticdata.Init()

	mux := http.NewServeMux()

	//TODO::
	scmClient, _ := github.New()
	mux.HandleFunc("/api/pipeline/", handlePipeline)
	mux.HandleFunc("/api/builds", controllers.HandleBuild)
	mux.HandleFunc("/api/tasks", controllers.HandleStep)
	mux.HandleFunc("/api/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			hook, err := scmClient.Webhooks.Parse(r, w)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(hook)
		}
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: mux,
	}

	registerShutdown(&server)

	err = server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
