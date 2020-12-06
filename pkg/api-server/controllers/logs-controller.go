package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	staticdata "github.com/task-executor/pkg/api-server/static-data"
	engine2 "github.com/task-executor/pkg/engine"
	"github.com/task-executor/pkg/engine/kube"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func tailStep(s *api.Step, engine engine2.Engine, logsChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Spec:::", s)

	spec := &engine2.Spec{
		Metadata: engine2.Metadata{
			UID:       fmt.Sprintf("te-step-%d", s.Id),
			Namespace: "default",
		},
	}

	reader, err := engine.Tail(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return
	}

	rd := bufio.NewReader(reader)
	for {
		dat, _, err := rd.ReadLine()
		if err != nil {
			log.Println(fmt.Sprintf("%d", s.Id) + " :: I am done!!!")
			return
		}
		log.Println(fmt.Sprintf("%d", s.Id) + string(dat))
		logsChan <- dat
	}

}

func HandleLogStream(w http.ResponseWriter, r *http.Request) {
	buildNumberStr := r.URL.Query()["buildNumber"][0]
	//stepNumber := r.URL.Query()["stepNumber"][0]

	buildNumber, _ := strconv.ParseInt(buildNumberStr, 10, 64)
	//TODO: Review the package redeclaration

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")

	f, ok := w.(http.Flusher)
	if !ok {
		return
	}

	io.WriteString(w, ": ping\n\n")
	f.Flush()

	engine, err := kube.NewFile("", "/Users/tahir/.kube/config", "")
	if err != nil {
		log.Println(err)
		return
	}

	logsChan := make(chan []byte)
	stepsChan := make(chan *api.Step)
	stepsSeenSoFar := make(map[int64]struct{})
	var wg sync.WaitGroup

	//Tail log steps
	go func() {
		for step := range stepsChan {
			_, ok := stepsSeenSoFar[step.Id]
			if !ok {
				wg.Add(1)
				go tailStep(step, engine, logsChan, &wg)
			}
		}
	}()

	//Read all steps initially
	steps, err := stepService.GetSteps(buildNumber)
	if err != nil {
		log.Println(err)
		return
	}

	for _, s := range steps {
		stepsChan <- s
		stepsSeenSoFar[s.Id] = struct{}{}
	}

	ctx, cancel := context.WithCancel(context.Background())
	streamer := staticdata.BuildStreamer.Subscribe(ctx, buildNumber)
	if err != nil {
		log.Println(err)
		return
	}

	//Read live events from DB + Polling results
	go func() {
		for {
			select {
			case step := <-streamer.StepChannel:
				stepsChan <- step

			case build := <-streamer.BuildChannel:
				if buildNumber == build.Id && build.Status.Id == 3 {
					//Allow steps to finish streaming
					wg.Wait()
					//TODO:??
					ctx.Done()
					cancel()
					close(logsChan)
					return
				}
			}
		}
	}()

	//Waiting for logs
	enc := json.NewEncoder(w)
	for l := range logsChan {
		io.WriteString(w, "data: ")
		enc.Encode(string(l))
		io.WriteString(w, "\n\n")
		f.Flush()
	}
}

func TestDeep(w http.ResponseWriter, r *http.Request) {
	res, _ := buildService.DeepFetch([]int64{137, 138})
	d, _ := json.Marshal(res)
	w.Write(d)
}
