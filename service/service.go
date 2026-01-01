package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type Service struct {
	ItemsService ItemsService
	CategoriesService CategoriesService
}

func NewService(Repo repository.Repository) Service {
	return Service{
		ItemsService: NewItemsService(Repo.ItemsRepo),
		CategoriesService: NewCategoriesService(Repo.CategoriesRepo),
	}
}