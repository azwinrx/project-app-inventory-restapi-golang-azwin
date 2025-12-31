package router

import (
	"project-app-inventory-restapi-golang-azwin/handler"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(itemsHandler handler.ItemsHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/items", itemsHandler.GetAllItemsHandler)

	return r
}