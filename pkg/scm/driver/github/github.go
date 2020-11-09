package github

import "github.com/task-executor/pkg/scm"

func New() (*scm.Client, error) {
	client := &scm.Client{}
	client.Webhooks = &webhookService{}
	return client, nil
}
