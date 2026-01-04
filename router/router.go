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
		// get low stock items (must be before /{id} to avoid conflict)
		r.Get("/low-stock", handler.ItemsHandler.GetLowStockItems)
		// get item by id
		r.Get("/{id}", handler.ItemsHandler.GetItemsById)
		// get all items
		r.Get("/", handler.ItemsHandler.GetAllItems)
		// create item
		r.Post("/", handler.ItemsHandler.CreateItems)
		// update item
		r.Put("/{id}", handler.ItemsHandler.UpdateItems)
		// delete item
		r.Delete("/{id}", handler.ItemsHandler.DeleteItems)
	})
	
	r.Route("/categories", func(r chi.Router) {
		// get category by id
		r.Get("/{id}", handler.CategoriesHandler.GetCategoriesById)
		// get all categories
		r.Get("/", handler.CategoriesHandler.GetAllCategories)
		// create category
		r.Post("/", handler.CategoriesHandler.CreateCategories)
		// update category
		r.Put("/{id}", handler.CategoriesHandler.UpdateCategories)
		// delete category
		r.Delete("/{id}", handler.CategoriesHandler.DeleteCategories)
	})

	r.Route("/racks", func(r chi.Router) {
		// get rack by id
		r.Get("/{id}", handler.RacksHandler.GetRacksById)
		// get all racks
		r.Get("/", handler.RacksHandler.GetAllRacks)
		// create rack
		r.Post("/", handler.RacksHandler.CreateRacks)
		// update rack
		r.Put("/{id}", handler.RacksHandler.UpdateRacks)
		// delete rack
		r.Delete("/{id}", handler.RacksHandler.DeleteRacks)
	})

	r.Route("/warehouses", func(r chi.Router) {
		// get warehouse by id
		r.Get("/{id}", handler.WarehousesHandler.GetWarehousesById)
		// get all warehouses
		r.Get("/", handler.WarehousesHandler.GetAllWarehouses)
		// create warehouse
		r.Post("/", handler.WarehousesHandler.CreateWarehouses)
		// update warehouse
		r.Put("/{id}", handler.WarehousesHandler.UpdateWarehouses)
		// delete warehouse
		r.Delete("/{id}", handler.WarehousesHandler.DeleteWarehouses)
	})

	r.Route("/users", func(r chi.Router) {
		// get user by id
		r.Get("/{id}", handler.UsersHandler.GetUsersByID)
		// get all users
		r.Get("/", handler.UsersHandler.GetAllUsers)
		// get user by email
		r.Get("/email", handler.UsersHandler.GetUsersByEmail)
		// create user
		r.Post("/", handler.UsersHandler.CreateUsers)
		// update user
		r.Put("/{id}", handler.UsersHandler.UpdateUsers)
		// delete user
		r.Delete("/{id}", handler.UsersHandler.DeleteUsers)
	})

	r.Route("/sales", func(r chi.Router) {
		// get sale by id
		r.Get("/{id}", handler.SalesHandler.GetSalesById)
		// get all sales
		r.Get("/", handler.SalesHandler.GetAllSales)
		// create sale
		r.Post("/", handler.SalesHandler.CreateSales)
		// update sale
		r.Put("/{id}", handler.SalesHandler.UpdateSales)
		// delete sale
		r.Delete("/{id}", handler.SalesHandler.DeleteSales)
	})

	r.Route("/reports", func(r chi.Router) {
		// get items report - total barang
		r.Get("/items", handler.ReportsHandler.GetItemsReport)
		// get sales report - penjualan
		r.Get("/sales", handler.ReportsHandler.GetSalesReport)
		// get revenue report - pendapatan
		r.Get("/revenue", handler.ReportsHandler.GetRevenueReport)
	})

	return r
}