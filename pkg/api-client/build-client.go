package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/task-executor/pkg/api"
)

type BuildClient struct {
	Config
}

type BuildFilterBuilder struct {
	query []string
	sort  []string
}

func NewBuildFilterBuilder() *BuildFilterBuilder {
	return &BuildFilterBuilder{}
}

func (filter *BuildFilterBuilder) WithId(id int64) *BuildFilterBuilder {
	filter.query = append(filter.query, fmt.Sprintf("id=%d", id))
	return filter
}

func (filter *BuildFilterBuilder) WithStatus(status int) *BuildFilterBuilder {
	filter.query = append(filter.query, fmt.Sprintf("status=%d", status))
	return filter
}

func (filter *BuildFilterBuilder) ByCreatedTs(asc bool) *BuildFilterBuilder {
	sign := ""
	if !asc {
		sign = "-"
	}

	filter.sort = append(filter.sort, sign+"createdTs")
	return filter
}

func NewBuildClient(config Config) BuildClient {
	bc := BuildClient{}
	bc.Config = config
	return bc
}

func (bc BuildClient) GetBuilds(filter *BuildFilterBuilder) ([]api.Build, error) {

	query := strings.Join(filter.query, "&")
	if len(filter.sort) > 0 {
		query += "&sortBy=" + strings.Join(filter.sort, ",")
	}

	client := http.Client{}
	res, err := client.Get(fmt.Sprintf("http://%s/api/builds?%s", bc.Host, query))

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var builds []api.Build
	err = json.Unmarshal(data, &builds)
	if err != nil {
		return nil, err
	}

	return builds, nil
}
