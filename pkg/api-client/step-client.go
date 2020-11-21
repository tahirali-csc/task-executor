package apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/task-executor/pkg/api"
	"io/ioutil"
	"net/http"
	"strings"
)

type StepClient struct {
	Config
}

type StepFilterBuilder struct {
	query []string
	sort  []string
}

func NewStepFilterBuilder() *StepFilterBuilder {
	return &StepFilterBuilder{}
}

func (filter *StepFilterBuilder) WithId(id int64) *StepFilterBuilder {
	filter.query = append(filter.query, fmt.Sprintf("id=%d", id))
	return filter
}

func (filter *StepFilterBuilder) WithStatus(status int) *StepFilterBuilder {
	filter.query = append(filter.query, fmt.Sprintf("status=%d", status))
	return filter
}

func (filter *StepFilterBuilder) ByCreatedTs(asc bool) *StepFilterBuilder {
	sign := ""
	if !asc {
		sign = "-"
	}

	filter.sort = append(filter.sort, sign+"createdTs")
	return filter
}

func NewStepClient(config Config) StepClient {
	sc := StepClient{}
	sc.Config = config
	return sc
}

func (bc StepClient) GetSteps(filter *StepFilterBuilder) ([]api.Step, error) {

	query := strings.Join(filter.query, "&")
	if len(filter.sort) > 0 {
		query += "&sortBy=" + strings.Join(filter.sort, ",")
	}

	client := http.Client{}
	res, err := client.Get(fmt.Sprintf("http://%s/api/steps?%s", bc.Host, query))

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var steps []api.Step
	err = json.Unmarshal(data, &steps)
	if err != nil {
		return nil, err
	}

	return steps, nil
}