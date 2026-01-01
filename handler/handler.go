package handler

import (
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
)

type Handler struct {
	ItemsHandler ItemsHandler
}

func NewHandler(service service.Service, config utils.Configuration) Handler {
	return Handler{
		ItemsHandler: NewItemsHandler(service.ItemsService, config),
	}
}
