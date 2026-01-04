package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type ItemsService interface {
	GetItemsById(id int) (*model.Items, error) 
	GetAllItems(page, limit int) ([]model.Items, int, error)
	GetLowStockItems(threshold int) ([]model.Items, error)
	CreateItems(data *model.Items) error
	UpdateItems(id int, data *model.Items) error
	DeleteItems(id int) error
}

type itemsService struct {
	Repo repository.ItemsRepository
}

func NewItemsService(repo repository.ItemsRepository) ItemsService {
	return &itemsService{Repo: repo}
}

func (s *itemsService) GetItemsById(id int) (*model.Items, error) {
	return s.Repo.GetItemsById(id)
}

func (s *itemsService) GetAllItems(page, limit int) ([]model.Items, int, error) {
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
	
	return s.Repo.GetAllItems(page, limit)
}

func (s *itemsService) GetLowStockItems(threshold int) ([]model.Items, error) {
	// Default threshold to 5 if not provided or invalid
	if threshold < 1 {
		threshold = 5
	}
	
	return s.Repo.GetLowStockItems(threshold)
}

func (s *itemsService) CreateItems(data *model.Items) error {
	return s.Repo.CreateItems(data)
}

func (s *itemsService) UpdateItems(id int, data *model.Items) error {
	return s.Repo.UpdateItems(id, data)
}

func (s *itemsService) DeleteItems(id int) error {
	return s.Repo.DeleteItems(id)
}