package storage

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
)

func AddUser(db *sql.DB, u models.User) error {
	_, err := db.Exec("INSERT INTO users (name, cpf, password) VALUES ($1, $2, $3)", u.Name, u.CPF, u.Password)
	return err
}

func GetUserByCPF(db *sql.DB, cpf string) *models.User {
	var user models.User
	err := db.QueryRow("SELECT name, cpf, password FROM users WHERE cpf = $1", cpf).Scan(&user.Name, &user.CPF, &user.Password)
	if err != nil {
		return nil
	}
	return &user
}
