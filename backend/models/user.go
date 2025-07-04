package models

type User struct {
	Name     string
	CPF      string
	Password string
}

type Topic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Session struct {
	ID      int   `json:"id"`
	TopicID int   `json:"topic_id"`
	OpenAt  int64 `json:"open_at"`
	CloseAt int64 `json:"close_at"`
}

type Vote struct {
	ID      int    `json:"id"`
	TopicID int    `json:"topic_id"`
	UserCPF string `json:"user_cpf"`
	Choice  string `json:"choice"`
}
