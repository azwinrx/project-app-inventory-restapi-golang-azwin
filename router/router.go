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

	return r
}