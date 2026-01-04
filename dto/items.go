package dto

import "time"

type ItemsRequest struct {
	Id         int       `json:"id"`
	CategoryId int       `json:"category_id" validate:"required"`
	RackId     int       `json:"rack_id" validate:"required"`
	Name       string    `json:"name" validate:"required,min=3"`
	Sku        string    `json:"sku" validate:"required"`
	Stock      int       `json:"stock" validate:"gte=0"`
	MinStock   int       `json:"min_stock" validate:"gte=0"`
	Price      float64   `json:"price" validate:"required,gte=0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ItemsResponse struct {
	Id         int       `json:"id"`
	CategoryId int       `json:"category_id"`
	RackId     int       `json:"rack_id"`
	Name       string    `json:"name"`
	Sku        string    `json:"sku"`
	Stock      int       `json:"stock"`
	MinStock   int       `json:"min_stock"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}