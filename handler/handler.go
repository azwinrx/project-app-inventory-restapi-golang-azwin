package handler

import (
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
)

type Handler struct {
	ItemsHandler ItemsHandler
	CategoriesHandler CategoriesHandler
	RacksHandler RacksHandler
	WarehousesHandler WarehousesHandler
	UsersHandler UsersHandler
	SalesHandler SalesHandler
	ReportsHandler ReportsHandler
}

func NewHandler(service service.Service, config utils.Configuration) Handler {
	return Handler{
		ItemsHandler: NewItemsHandler(service.ItemsService, config),
		CategoriesHandler: NewCategoriesHandler(service.CategoriesService, config),
		RacksHandler: NewRacksHandler(service.RacksService, config),
		WarehousesHandler: NewWarehousesHandler(service.WarehousesService, config),
		UsersHandler: NewUsersHandler(service.UsersService, config),
		SalesHandler: NewSalesHandler(service.SalesService, config),
		ReportsHandler: NewReportsHandler(service.ReportsService, config),
	}
}
