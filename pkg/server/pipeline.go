package main

import (
	cp "github.com/tahirali-csc/client-api"
	"log"
)

//func Run(image string, cmd []string) error {
//
//	var cmdToRun []string
//	for _, c := range cmd {
//		cmdToRun = append(cmdToRun, c)
//	}
//
//	input := make(map[string]interface{})
//	input["image"] = image
//	input["command"] = cmdToRun
//
//	body, err := json.Marshal(input)
//	if err != nil {
//		return err
//	}
//
//	client := http.Client{}
//	res, err := client.Post("http://host.docker.internal:8081/api/tasks/",
//		"application/json", bytes.NewReader(body))
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	dat, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	log.Println(string(dat))
//	return nil
//}

func pipeline() {
	log.Println("I am running a pipeline")
	cp.Run("alpine:latest", []string{"ls -la"})

	log.Println("Running next")
	cp.Run("alpine:latest", []string{"date"})

	log.Println("Running 2")
	cp.Run("alpine:latest", []string{"ls"})
}

func main() {
	pipeline()
}


