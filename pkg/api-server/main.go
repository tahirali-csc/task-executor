package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api-server/config"
	"github.com/task-executor/pkg/api-server/controllers"
	"github.com/task-executor/pkg/api-server/dbstore"
	"github.com/task-executor/pkg/scm/driver/github"
	"github.com/task-executor/pkg/utils"

	"net/http"
)

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

	//router := http.NewServeMux()
	router := mux.NewRouter()
	//TODO::
	scmClient, _ := github.New()
	router.HandleFunc("/api/builds", controllers.HandleBuild)
	router.HandleFunc("/api/steps", controllers.HandleStep)
	router.HandleFunc("/api/steps/{id}/status", controllers.HandleStepStatus)
	router.HandleFunc("/api/steps/{id}/status/{status}", controllers.HandleStepStatusUpdate)
	router.HandleFunc("/api/logs", controllers.HandleLogStream)
	router.HandleFunc("/api/callback", func(w http.ResponseWriter, r *http.Request) {
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
		Handler: router,
	}

	registerShutdown(&server)

	//go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Error(err)
		}
	//}()

	////TODO: Also start runner at the same time
	//go func() {
	//	time.Sleep(time.Second * 5)
	//	runner := runner.NewRunner()
	//	runner.Run()
	//}()

	//for {
	//}

}
