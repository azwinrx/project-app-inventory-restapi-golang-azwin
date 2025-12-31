package model

type SaleItems struct {
	Id       int     `json:"id"`
	SaleId   int     `json:"sale_id"`
	ItemId   int     `json:"item_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}