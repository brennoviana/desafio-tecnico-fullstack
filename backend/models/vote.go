package models

type Vote struct {
	ID      int    `json:"id"`
	TopicID int    `json:"topic_id"`
	UserID  int    `json:"user_id"`
	Choice  string `json:"choice"`
}
