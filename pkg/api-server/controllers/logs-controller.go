package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/task-executor/pkg/api"
	engine2 "github.com/task-executor/pkg/engine"
	"github.com/task-executor/pkg/engine/kube"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func HandleLogStream(w http.ResponseWriter, r *http.Request) {

	buildNumber := r.URL.Query()["buildNumber"][0]
	//stepNumber := r.URL.Query()["stepNumber"][0]

	bid, _ := strconv.ParseInt(buildNumber, 10, 64)
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
	stepsDoneChan := make(chan bool)
	var wg sync.WaitGroup

	go tailLogs(stepsChan, stepsDoneChan, logsChan, engine, &wg)
	go streamSteps(bid, stepsChan, stepsDoneChan)

	enc := json.NewEncoder(w)
	for l := range logsChan {
		io.WriteString(w, "data: ")
		enc.Encode(string(l))
		io.WriteString(w, "\n\n")
		f.Flush()
	}

	log.Print("Done streaming...")
}

func streamSteps(buildId int64, stepsChan chan *api.Step, stepsDoneChan chan bool) {

	steps, err := stepService.GetSteps(buildId)
	if err != nil {
		log.Println(err)
		return
	}

	for _, s := range steps {
		log.Println("Sending:::", s.Id)
		stepsChan <- s
	}

	stepsDoneChan <- true

	//close(stepsChan)
}

func tailLogs(stepsChan chan *api.Step, stepsDoneChan chan bool, logsChan chan []byte,
	engine engine2.Engine, wg *sync.WaitGroup) {

	//for s := range stepsChan {
L:
	for {
		select {
		case s := <-stepsChan:
			wg.Add(1)
			log.Println("Received::", s.Id)
			go tailStep(s, engine, logsChan, wg)

		case <-stepsDoneChan:
			break L
		}
	}

	wg.Wait()
	close(logsChan)

	//}
}

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
