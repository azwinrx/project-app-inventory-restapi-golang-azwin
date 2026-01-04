package dto

import "time"

type SalesRequest struct {
    UserId int                    `json:"user_id" validate:"required"`
    Items  []SaleItemRequest      `json:"items" validate:"required,dive"`  // Gabung detail items
}

type SaleItemRequest struct {
    ItemId   int     `json:"item_id" validate:"required"`
    Quantity int     `json:"quantity" validate:"required,gte=1"`
    Price    float64 `json:"price" validate:"required,gte=0"`
}

type SalesResponse struct {
    Id          int               `json:"id"`
    UserId      int               `json:"user_id"`
    TotalAmount float64           `json:"total_amount"`
    Items       []SaleItemResponse `json:"items"`  // Include detail
    CreatedAt   time.Time         `json:"created_at"`
}

type SaleItemResponse struct {
    Id       int     `json:"id"`
    ItemId   int     `json:"item_id"`
    Quantity int     `json:"quantity"`
    Price    float64 `json:"price"`
    Subtotal float64 `json:"subtotal"`
}