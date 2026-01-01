package dto

import "time"

type ItemsRequest struct {
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