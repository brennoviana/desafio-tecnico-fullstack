package connection

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/config"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	user := config.AppConfig.DBUser
	password := config.AppConfig.DBPassword
	dbname := config.AppConfig.DBName
	host := config.AppConfig.DBHost
	port := config.AppConfig.DBPort

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        cpf TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    )`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS topics (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		topic_id INTEGER NOT NULL REFERENCES topics(id),
		open_at BIGINT NOT NULL,
		close_at BIGINT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS votes (
		id SERIAL PRIMARY KEY,
		topic_id INTEGER NOT NULL REFERENCES topics(id),
		user_cpf TEXT NOT NULL REFERENCES users(cpf),
		choice TEXT NOT NULL,
		UNIQUE(topic_id, user_cpf)
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
