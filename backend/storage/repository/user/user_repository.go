package user

import (
	"database/sql"
	"desafio-tecnico-fullstack/backend/models"
)

type UserRepository interface {
	AddUser(u models.User) error
	GetUserByCPF(cpf string) *models.User
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) AddUser(u models.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, cpf, password) VALUES ($1, $2, $3)", u.Name, u.CPF, u.Password)
	return err
}

func (r *userRepository) GetUserByCPF(cpf string) *models.User {
	var user models.User
	err := r.db.QueryRow("SELECT id, name, cpf, password FROM users WHERE cpf = $1", cpf).Scan(&user.ID, &user.Name, &user.CPF, &user.Password)
	if err != nil {
		return nil
	}
	return &user
}
