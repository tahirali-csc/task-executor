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
	"strings"
	"sync"
)

type logInfo struct {
	StepId  int64
	Message string
}

func tailStep(s *api.Step, engine engine2.Engine, logsChan chan *logInfo, wg *sync.WaitGroup, notifyChan chan []byte) {
	defer wg.Done()
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

	//d, _ := json.Marshal(map[string]interface{}{
	//	"StepId": s.Id,
	//	"Status": "Started",
	//})
	//notifyChan <- d

	rd := bufio.NewReader(reader)
	for {
		dat, _, err := rd.ReadLine()
		if err != nil {
			log.Println(fmt.Sprintf("%d", s.Id) + " :: I am done!!!")
			d, _ := json.Marshal(map[string]interface{}{
				"StepId": s.Id,
				"Status": "Finished",
			})
			notifyChan <- d
			return
		}
		log.Debug(fmt.Sprintf("%d", s.Id) + string(dat))
		logsChan <- &logInfo{StepId: s.Id, Message: string(dat)}
	}

}

func HandleLogStream(w http.ResponseWriter, r *http.Request) {
	buildNumberStr := r.URL.Query()["buildNumber"][0]
	//stepNumber := r.URL.Query()["stepNumber"][0]

	buildNumber, _ := strconv.ParseInt(buildNumberStr, 10, 64)
	//TODO: Review the package redeclaration
	log.Debug("Straming Build Number==", buildNumber)

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

	logsChan := make(chan *logInfo)
	stepsChan := make(chan *api.Step)
	notifyChan := make(chan []byte)
	doneStreamingChan := make(chan bool)
	var wg sync.WaitGroup

	//Tail log steps
	go func() {
		stepsSeenSoFar := make(map[int64]struct{})
		for step := range stepsChan {
			//log.Println("Status:::", step.Status)

			//Wait if pending
			if step.Status.Name == api.PendingBuildStatus {
				continue
			}

			_, ok := stepsSeenSoFar[step.Id]
			if !ok {
				stepsSeenSoFar[step.Id] = struct{}{}

				log.Debug("Starting tail log....", step.Id, stepsSeenSoFar)
				wg.Add(1)

				d, _ := json.Marshal(map[string]interface{}{
					"StepId": step.Id,
					"Status": "Started",
				})
				notifyChan <- d

				go tailStep(step, engine, logsChan, &wg, notifyChan)
			}
		}
	}()

	go func() {
		//Read all steps initially
		steps, err := stepService.GetSteps(buildNumber)
		if err != nil {
			log.Error(err)
			return
		}

		for _, s := range steps {
			stepsChan <- s
		}

		//Subscribe to streamer events
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
					log.Debug("Got build event:::", build)
					if buildNumber == build.Id && build.Status.Name == api.FinishedBuildStatus {
						//Allow steps to finish streaming
						wg.Wait()
						//TODO:??
						ctx.Done()
						cancel()
						doneStreamingChan <- true
						//close(logsChan)
						return
					}
				}
			}
		}()

	}()

	enc := json.NewEncoder(w)
L:
	for {
		select {
		case l := <-logsChan:
			_ = formatSSEMessage(w, "", l, enc, f)
		case notify := <-notifyChan:
			log.Debug("Sending Notify", string(notify))
			//TODO: Review
			//dat, _ := json.Marshal(notify)
			nn := formatSSE("notify", string(notify))
			w.Write(nn)
			f.Flush()
		case <-doneStreamingChan:
			break L
		}
	}

	////Tail log steps
	//go func() {
	//	stepsSeenSoFar := make(map[int64]struct{})
	//	for step := range stepsChan {
	//		//log.Println("Status:::", step.Status)
	//
	//		//Wait if pending
	//		if step.Status.Name == api.PendingBuildStatus {
	//			continue
	//		}
	//
	//		_, ok := stepsSeenSoFar[step.Id]
	//		if !ok {
	//			stepsSeenSoFar[step.Id] = struct{}{}
	//
	//			log.Debug("Starting tail log....", step.Id, stepsSeenSoFar)
	//			wg.Add(1)
	//
	//			d, _ := json.Marshal(map[string]interface{}{
	//				"StepId": step.Id,
	//				"Status": "Started",
	//			})
	//			notifyChan <- d
	//
	//			go tailStep(step, engine, logsChan, &wg, notifyChan)
	//		}
	//	}
	//}()

	////Read all steps initially
	//steps, err := stepService.GetSteps(buildNumber)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//for _, s := range steps {
	//	stepsChan <- s
	//}

	////Subscribe to streamer events
	//ctx, cancel := context.WithCancel(context.Background())
	//streamer := staticdata.BuildStreamer.Subscribe(ctx, buildNumber)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	////Read live events from DB + Polling results
	//go func() {
	//	for {
	//		select {
	//		case step := <-streamer.StepChannel:
	//			stepsChan <- step
	//
	//		case build := <-streamer.BuildChannel:
	//			log.Debug("Got build event:::", build)
	//			if buildNumber == build.Id && build.Status.Name == api.FinishedBuildStatus {
	//				//Allow steps to finish streaming
	//				wg.Wait()
	//				//TODO:??
	//				ctx.Done()
	//				cancel()
	//				doneStreamingChan <- true
	//				//close(logsChan)
	//				return
	//			}
	//		}
	//	}
	//}()

	//Waiting for logs
	//	enc := json.NewEncoder(w)
	//L:
	//	for {
	//		select {
	//		case l := <-logsChan:
	//			_ = formatSSEMessage(w, "", l, enc, f)
	//		case notify := <-notifyChan:
	//			log.Debug("Sending Notify", string(notify))
	//			//TODO: Review
	//			//dat, _ := json.Marshal(notify)
	//			nn := formatSSE( "notify", string(notify))
	//			w.Write(nn)
	//			f.Flush()
	//		case <-doneStreamingChan:
	//			break L
	//		}
	//	}

	//for l := range logsChan {
	//	io.WriteString(w, "data: ")
	//	enc.Encode(string(l))
	//	io.WriteString(w, "\n\n")
	//	f.Flush()
	//}

	/*w.Write(formatSSE("close", ""))
	f.Flush()*/
	log.Println("CLosing Stream")
	d, _ := json.Marshal(map[string]interface{}{
		"Status": "Done",
	})
	nn := formatSSE("close", string(d))
	w.Write(nn)
	f.Flush()
	//_ = formatSSEMessage(w, "close", []byte(""), enc, f)
}

func formatSSEMessage(w http.ResponseWriter, event string, data interface{}, enc *json.Encoder, flusher http.Flusher) error {
	if len(event) > 0 {
		if _, err := io.WriteString(w, "event: "+event+"\n"); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "data: "); err != nil {
		return err
	}
	//if err := enc.Encode(string(data)); err != nil {
	//	return err
	//}
	//if err := enc.Encode(data); err != nil {
	//	return err
	//}
	d, _ := json.Marshal(data)
	w.Write(d)

	if _, err := io.WriteString(w, "\n\n"); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	//eventPayload = eventPayload + "data: " + data + "\n\n"
	//eventPayload = eventPayload + "data: " + data + "\n"
	return []byte(eventPayload + "\n")
}

func TestDeep(w http.ResponseWriter, r *http.Request) {
	res, _ := buildService.DeepFetch([]int64{137, 138})
	d, _ := json.Marshal(res)
	w.Write(d)
}
