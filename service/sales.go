package service

import (
	"errors"
	"project-app-inventory-restapi-golang-azwin/dto"
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type SalesService interface {
	GetSalesById(id int) (*dto.SalesResponse, error)
	GetAllSales(page, limit int) ([]dto.SalesResponse, int, error)
	CreateSales(data *dto.SalesRequest) error
	UpdateSales(id int, data *dto.SalesRequest) error
	DeleteSales(id int) error
}

type salesService struct {
	Repo repository.SalesRepository
}

func NewSalesService(repo repository.SalesRepository) SalesService {
	return &salesService{Repo: repo}
}

func (s *salesService) GetSalesById(id int) (*dto.SalesResponse, error) {
	sale, items, err := s.Repo.GetSalesById(id)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	var itemsResponse []dto.SaleItemResponse
	for _, item := range items {
		itemsResponse = append(itemsResponse, dto.SaleItemResponse{
			Id:       item.Id,
			ItemId:   item.ItemId,
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: item.Subtotal,
		})
	}

	response := &dto.SalesResponse{
		Id:          sale.Id,
		UserId:      sale.UserId,
		TotalAmount: sale.TotalAmount,
		Items:       itemsResponse,
		CreatedAt:   sale.CreatedAt,
	}

	return response, nil
}

func (s *salesService) GetAllSales(page, limit int) ([]dto.SalesResponse, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	sales, total, err := s.Repo.GetAllSales(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to DTO with items detail
	var salesResponse []dto.SalesResponse
	for _, sale := range sales {
		// Convert items to DTO
		var itemsResponse []dto.SaleItemResponse
		for _, item := range sale.Items {
			itemsResponse = append(itemsResponse, dto.SaleItemResponse{
				Id:       item.Id,
				ItemId:   item.ItemId,
				Quantity: item.Quantity,
				Price:    item.Price,
				Subtotal: item.Subtotal,
			})
		}

		salesResponse = append(salesResponse, dto.SalesResponse{
			Id:          sale.Id,
			UserId:      sale.UserId,
			TotalAmount: sale.TotalAmount,
			Items:       itemsResponse,
			CreatedAt:   sale.CreatedAt,
		})
	}

	return salesResponse, total, nil
}

func (s *salesService) CreateSales(data *dto.SalesRequest) error {
	// Validate request
	if data.UserId <= 0 {
		return errors.New("user_id is required")
	}
	if len(data.Items) == 0 {
		return errors.New("at least one item is required")
	}

	// Calculate total amount
	var totalAmount float64
	var saleItems []model.SaleItems
	for _, item := range data.Items {
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}
		if item.Price <= 0 {
			return errors.New("price must be greater than 0")
		}

		subtotal := float64(item.Quantity) * item.Price
		totalAmount += subtotal

		saleItems = append(saleItems, model.SaleItems{
			ItemId:   item.ItemId,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	// Create sale model
	sale := &model.Sales{
		UserId:      data.UserId,
		TotalAmount: totalAmount,
	}

	return s.Repo.CreateSales(sale, saleItems)
}

func (s *salesService) UpdateSales(id int, data *dto.SalesRequest) error {
	// Validate request
	if data.UserId <= 0 {
		return errors.New("user_id is required")
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range data.Items {
		subtotal := float64(item.Quantity) * item.Price
		totalAmount += subtotal
	}

	// Update sale model
	sale := &model.Sales{
		UserId:      data.UserId,
		TotalAmount: totalAmount,
	}

	return s.Repo.UpdateSales(id, sale)
}

func (s *salesService) DeleteSales(id int) error {
	return s.Repo.DeleteSales(id)
}
