package model

import "time"

type Racks struct {
	Id       int    `json:"id"`
	WarehouseId int    `json:"warehouse_id"`
	Name    string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}