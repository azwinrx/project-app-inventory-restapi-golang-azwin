package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type ItemsService interface {
	GetAllItems(page, limit int) ([]model.Items, int, error)
	Create(i *model.Items) error
	Update(id int, item *model.Items) error
	Delete(id int) error
}

type itemsService struct {
	Repo repository.ItemsRepository
}

func NewItemsService(repo repository.ItemsRepository) ItemsService {
	return &itemsService{Repo: repo}
}

func (i *itemsService) GetAllItems(page, limit int) ([]model.Items, int, error) {
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
	
	return i.Repo.GetAllItems(page, limit)
}

func (i *itemsService) Create(id *model.Items) error {
	return i.Repo.CreateItems(id)
}

func (i *itemsService) Update(id int, item *model.Items) error {
	return i.Repo.UpdateItems(id, item)
}

func (i *itemsService) Delete(id int) error {
	return i.Repo.DeleteItems(id)
}