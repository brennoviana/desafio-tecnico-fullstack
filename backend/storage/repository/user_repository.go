package repository

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddUser(u models.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, cpf, password) VALUES ($1, $2, $3)", u.Name, u.CPF, u.Password)
	return err
}

func (r *UserRepository) GetUserByCPF(cpf string) *models.User {
	var user models.User
	err := r.db.QueryRow("SELECT name, cpf, password FROM users WHERE cpf = $1", cpf).Scan(&user.Name, &user.CPF, &user.Password)
	if err != nil {
		return nil
	}
	return &user
}