package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type Service struct {
	ItemsService ItemsService
}

func NewService(Repo repository.Repository) Service {
	return Service{
		ItemsService: NewItemsService(Repo.ItemsRepo),
	}
}