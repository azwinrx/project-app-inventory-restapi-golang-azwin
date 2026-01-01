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
}

func NewHandler(service service.Service, config utils.Configuration) Handler {
	return Handler{
		ItemsHandler: NewItemsHandler(service.ItemsService, config),
		CategoriesHandler: NewCategoriesHandler(service.CategoriesService, config),
		RacksHandler: NewRacksHandler(service.RacksService, config),
		WarehousesHandler: NewWarehousesHandler(service.WarehousesService, config),
	}
}
