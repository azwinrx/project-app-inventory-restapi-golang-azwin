package dto

import "time"

type CategoriesRequest struct {
		Id	   int    `json:"id"`
	Name     string `json:"name" validate:"required,min=3"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoriesResponse struct {
		Id	   int    `json:"id"`
	Name     string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}