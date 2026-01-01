package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type Service struct {
	ItemsService ItemsService
	CategoriesService CategoriesService
	RacksService RacksService
	WarehousesService WarehousesService
}

func NewService(Repo repository.Repository) Service {
	return Service{
		ItemsService: NewItemsService(Repo.ItemsRepo),
		CategoriesService: NewCategoriesService(Repo.CategoriesRepo),
		RacksService: NewRacksService(Repo.RacksRepo),
		WarehousesService: NewWarehousesService(Repo.WarehousesRepo),
	}
}