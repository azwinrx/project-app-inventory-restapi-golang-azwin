package router

import (
	"project-app-inventory-restapi-golang-azwin/handler"
	mCostume "project-app-inventory-restapi-golang-azwin/middleware"
	"project-app-inventory-restapi-golang-azwin/service"

	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler handler.Handler, service service.Service, log *zap.Logger) http.Handler {
	r := chi.NewRouter()


	mw := mCostume.NewMiddlewareCustome(service, log)
	r.Mount("/", ApiV1(handler, mw))

	return r
}

func ApiV1(handler handler.Handler, mw mCostume.MiddlewareCostume) *chi.Mux{
	r := chi.NewRouter()
	r.Use(mw.Logging)
	
	r.Route("/items", func(r chi.Router) {
		r.Get("/{id}", handler.ItemsHandler.GetItemsById)
		r.Get("/", handler.ItemsHandler.GetAllItems)
		r.Post("/", handler.ItemsHandler.CreateItems)
		r.Put("/{id}", handler.ItemsHandler.UpdateItems)
		r.Delete("/{id}", handler.ItemsHandler.DeleteItems)
	})
	
	r.Route("/categories", func(r chi.Router) {
		r.Get("/{id}", handler.CategoriesHandler.GetCategoriesById)
		r.Get("/", handler.CategoriesHandler.GetAllCategories)
		r.Post("/", handler.CategoriesHandler.CreateCategories)
		r.Put("/{id}", handler.CategoriesHandler.UpdateCategories)
		r.Delete("/{id}", handler.CategoriesHandler.DeleteCategories)
	})

	r.Route("/racks", func(r chi.Router) {
		r.Get("/{id}", handler.RacksHandler.GetRacksById)
		r.Get("/", handler.RacksHandler.GetAllRacks)
		r.Post("/", handler.RacksHandler.CreateRacks)
		r.Put("/{id}", handler.RacksHandler.UpdateRacks)
		r.Delete("/{id}", handler.RacksHandler.DeleteRacks)
	})

	r.Route("/warehouses", func(r chi.Router) {
		r.Get("/{id}", handler.WarehousesHandler.GetWarehousesById)
		r.Get("/", handler.WarehousesHandler.GetAllWarehouses)
		r.Post("/", handler.WarehousesHandler.CreateWarehouses)
		r.Put("/{id}", handler.WarehousesHandler.UpdateWarehouses)
		r.Delete("/{id}", handler.WarehousesHandler.DeleteWarehouses)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", handler.UsersHandler.GetUsersByID)
		r.Get("/", handler.UsersHandler.GetAllUsers)
		r.Get("/email/{email}", handler.UsersHandler.GetUsersByEmail)
		r.Post("/", handler.UsersHandler.CreateUsers)
	})

	return r
}