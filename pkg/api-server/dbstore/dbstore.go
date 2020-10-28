package dbstore

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/task-executor/pkg/api-server/config"
)

func GetDb() (*sql.DB, error) {
	//TODO : Will review
	config := config.Get()
	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
		config.Database.Name, config.Database.User, config.Database.Password, config.Database.Host)
	return sql.Open("postgres", connString)
}

func Release(db *sql.DB) error {
	return db.Close()
}
