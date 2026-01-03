package dto

import "time"

type SalesRequest struct {
    UserId int                    `json:"user_id"`
    Items  []SaleItemRequest      `json:"items"`  // Gabung detail items
}

type SaleItemRequest struct {
    ItemId   int     `json:"item_id"`
    Quantity int     `json:"quantity"`
    Price    float64 `json:"price"`
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