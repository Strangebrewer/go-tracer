package example

import "time"

type Example struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateExampleRequest struct {
	Name string `json:"name"`
}

type UpdateExampleRequest struct {
	Name string `json:"name"`
}
