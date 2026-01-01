package dto

import "time"

type WarehousesRequest struct {
	Id        int       `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Location  string    `json:"location" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WarehousesResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
