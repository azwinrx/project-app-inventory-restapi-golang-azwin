package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type Service struct {
	ItemsRepo ItemsService
}

func NewService(itemsRepo repository.ItemsRepository) *Service {
	return &Service{
		ItemsRepo: NewItemsService(itemsRepo),
	}
}