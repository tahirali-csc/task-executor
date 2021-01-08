package logs

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"io"
	"io/ioutil"

	"github.com/task-executor/pkg/core"
)

func New(db *sql.DB) core.LogStore {
	return &dbLogs{
		DataSource: db,
	}
}

type dbLogs struct {
	DataSource *sql.DB
}

// Find returns a log stream from the datastore.
func (l *dbLogs) Find(ctx context.Context, stepId int64) (io.ReadCloser, error) {
	var logs []byte
	row := l.DataSource.QueryRow(`SELECT log_data FROM logs WHERE step_id=$1`, stepId)
	err := row.Scan(&logs)

	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(
		bytes.NewBuffer(logs),
	), err
}

// Update writes copies the log stream from Reader r to the datastore.
func (l *dbLogs) Upload(ctx context.Context, stepId int64, log io.Reader) error {
	reader := bufio.NewReader(log)
	chunk := make([]byte, 100)

	for {
		read, err := reader.Read(chunk) //ReadString and ReadLine() also applicable or alternative
		if err != nil {
			break
		}

		insertStmt := `INSERT INTO logs(step_id, log_data) VALUES($1, $2)
		ON conflict(step_id) DO UPDATE SET log_data = logs.log_data || $2`

		_, err = l.DataSource.Exec(insertStmt, stepId, chunk[:read])
		if err != nil {
			return err
		}
	}

	return nil
}
