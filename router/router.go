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
		r.Get("/", handler.ItemsHandler.GetAllItems)
		r.Post("/", handler.ItemsHandler.CreateItem)
	})
	
	return r
}