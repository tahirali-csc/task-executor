package main

import "net/http"

func PollPendingBuilds() {
	client := http.Client{}
	client.Get("http://localhost:8080/api/build/")

}

func main() {

	for {

	}
}
