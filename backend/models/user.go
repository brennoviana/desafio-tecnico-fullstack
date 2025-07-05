package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}
