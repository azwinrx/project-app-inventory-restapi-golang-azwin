package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type ItemsService interface {
	GetAllItems(page, limit int) ([]model.Items, int, error)
}

type itemsService struct {
	Repo repository.ItemsRepository
}

func NewItemsService(repo repository.ItemsRepository) ItemsService {
	return &itemsService{Repo: repo}
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

