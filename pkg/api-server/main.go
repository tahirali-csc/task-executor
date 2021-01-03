package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/task-executor/pkg/api-server/events"
	"github.com/task-executor/pkg/api-server/services"
	staticdata "github.com/task-executor/pkg/api-server/static-data"

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

func setHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//anyone can make a CORS request (not recommended in production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//only allow GET, POST, and OPTIONS
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		//Since I was building a REST API that returned JSON, I set the content type to JSON here.
		w.Header().Set("Content-Type", "application/json")
		//Allow requests to have the following headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")
		//if it's just an OPTIONS request, nothing other than the headers in the response is needed.
		//This is essential because you don't need to handle the OPTIONS requests in your handlers now
		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	//In GoLand/Intellij set work directory pointing to api-server
	configFile := flag.String("config file", "config.yaml", "")
	flag.Parse()

	//Set default logging level
	utils.InitLogs(log.DebugLevel)

	//Load application config
	config, err := config.Load(*configFile)
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

	go func() {
		staticdata.EventStream = events.NewDBEvents()
		staticdata.EventStream.NewTopic("build")
		staticdata.EventStream.NewTopic("step")

		staticdata.BuildStreamer = services.NewBuildStreamer(staticdata.EventStream)
		go staticdata.BuildStreamer.Start()
		staticdata.EventStream.Start()
	}()

	//router := http.NewServeMux()
	router := mux.NewRouter()
	//TODO::
	scmClient, _ := github.New()
	router.HandleFunc("/api/builds", controllers.HandleBuild)
	router.HandleFunc("/api/builds/{id}/status/{status}", controllers.HandleBuildStatus)
	router.HandleFunc("/api/steps", controllers.HandleStep)
	router.HandleFunc("/api/steps/{id}/status", controllers.HandleStepStatus)
	router.HandleFunc("/api/steps/{id}/status/{status}", controllers.HandleStepStatus)
	router.HandleFunc("/api/steps/{id}/logs", controllers.HandleStepLogsUpload)
	router.HandleFunc("/api/logs", controllers.HandleLogStream)
	// router.HandleFunc("/api/testdeep", controllers.TestDeep)

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
		Handler: setHeaders(router),
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
