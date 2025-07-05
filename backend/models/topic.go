package models

type Topic struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
