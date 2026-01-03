package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type Service struct {
	ItemsService ItemsService
	CategoriesService CategoriesService
	RacksService RacksService
	WarehousesService WarehousesService
	UsersService UsersService
	SalesService SalesService
	ReportsService ReportsService
}

func NewService(Repo repository.Repository) Service {
	return Service{
		ItemsService: NewItemsService(Repo.ItemsRepo),
		CategoriesService: NewCategoriesService(Repo.CategoriesRepo),
		RacksService: NewRacksService(Repo.RacksRepo),
		WarehousesService: NewWarehousesService(Repo.WarehousesRepo),
		UsersService: NewUsersService(Repo.UsersRepo),
		SalesService: NewSalesService(Repo.SalesRepo),
		ReportsService: NewReportsService(Repo.ReportsRepo),
	}
}