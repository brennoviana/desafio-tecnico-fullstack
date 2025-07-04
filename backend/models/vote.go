package models

type Vote struct {
	ID      int    `json:"id"`
	TopicID int    `json:"topic_id"`
	UserCPF string `json:"user_cpf"`
	Choice  string `json:"choice"`
}
