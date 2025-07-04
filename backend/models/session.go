package models

type Session struct {
	ID      int   `json:"id"`
	TopicID int   `json:"topic_id"`
	OpenAt  int64 `json:"open_at"`
	CloseAt int64 `json:"close_at"`
}
