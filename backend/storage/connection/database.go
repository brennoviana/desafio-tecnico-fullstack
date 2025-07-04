package connection

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/config"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	user := config.AppConfig.Database.User
	password := config.AppConfig.Database.Password
	dbname := config.AppConfig.Database.Name
	host := config.AppConfig.Database.Host
	port := config.AppConfig.Database.Port

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
