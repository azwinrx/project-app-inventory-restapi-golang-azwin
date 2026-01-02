package dto

import "time"

type RacksRequest struct {
	Id          int       `json:"id"`
	WarehouseId int       `json:"warehouse_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RacksResponse struct {
	Id          int       `json:"id"`
	WarehouseId int       `json:"warehouse_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
