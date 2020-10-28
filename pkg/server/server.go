package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func runContainer(image string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(ioutil.Discard, reader)
	reader.Close()

	pipelineVol := mount.Mount{
		Target: "/brave",
		Source: "/Users/tahir/workspace/rnd-projects/task-pipeline/pkg/server/",
		Type:   mount.TypeBind,
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        image,
		Cmd:          []string{"/bin/sh", "-c", "go Run /brave/pipeline.go"},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts:     []mount.Mount{pipelineVol},
		AutoRemove: true,
	}, nil, "")
	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	rdr := bufio.NewReader(out)
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		log.Print(string(line))
	}

	return nil
}

func handleRunTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := runContainer("golang:latest")
		if err != nil {
			log.Println(err)
		}
	} else if r.Method == "GET" {
		fmt.Fprintf(w, "Hello")
	}
}

func handleGophy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "Gophy")
	} else if r.Method == "GET" {

	}
}

type Task struct {
}

func (t *Task) Run(msg string, reply *string) error {
	//*reply = "kjgdfgkjdgkfdghkd"
	runContainer("golang:latest")
	*reply = "done"
	return nil
}

func (t *Task) Get(msg string, reply *string) error {
	*reply = "Hey Gopher!!!"
	return nil
}

func main() {
	task := new(Task)
	err := rpc.Register(task)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}

	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":8080")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}



	mux := http.NewServeMux()
	//
	mux.HandleFunc("/api/tasks/", handleRunTask)
	//mux.HandleFunc("/api/gophy/", handleGophy)

	//http.ListenAndServe(":8080", mux)

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}

	//mux := http.NewServeMux()
	//
	//mux.HandleFunc("/api/tasks/", handleRunTask)
	//mux.HandleFunc("/api/gophy/", handleGophy)
	//
	//http.ListenAndServe(":8080", mux)
}
