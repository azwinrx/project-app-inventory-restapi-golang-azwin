package model

import "time"

type Sales struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}
