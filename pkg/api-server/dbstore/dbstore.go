package dbstore

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/task-executor/pkg/api-server/config"
)

var DataSource *sql.DB

func Init(config *config.AppConfig) error {
	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
		config.Database.Name, config.Database.User, config.Database.Password, config.Database.Host)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	DataSource = db
	return err
}

func Close() error {
	return DataSource.Close()
}
