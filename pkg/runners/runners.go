package runners

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type RunOptions struct {
	Image string
}

func Run(commands []string, runOpts RunOptions) error {
	body, err := json.Marshal(runOpts)
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Post("http://localhost/api/tasks", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	return nil
}
