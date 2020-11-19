package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	engine2 "github.com/task-executor/pkg/engine"
	"github.com/task-executor/pkg/engine/kube"
	"io"
	"net/http"
	"strconv"
)

func HandleLogStream(w http.ResponseWriter, r *http.Request) {

	//buildNumber := r.URL.Query()["buildNumber"][0]
	stepNumber := r.URL.Query()["stepNumber"][0]

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

	sid, _ := strconv.ParseInt(stepNumber, 10, 64)

	spec := &engine2.Spec{
		Metadata: engine2.Metadata{
			UID:       fmt.Sprintf("te-step-%d", sid),
			Namespace: "default",
		},
	}

	reader, err := engine.Tail(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return
	}

	enc := json.NewEncoder(w)
	rd := bufio.NewReader(reader)

	for {
		da, _, err := rd.ReadLine()

		if err != nil {
			break
		}
		log.Println(string(da))

		io.WriteString(w, "data: ")
		enc.Encode(string(da))
		io.WriteString(w, "\n\n")
		f.Flush()

	}

	//io.WriteString(w, "event: error\ndata: eof\n\n")
	//f.Flush()

}
