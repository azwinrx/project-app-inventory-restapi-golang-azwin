package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type WarehousesService interface {
	GetWarehousesById(id int) (*model.Warehouses, error)
	GetAllWarehouses(page, limit int) ([]model.Warehouses, int, error)
	CreateWarehouses(data *model.Warehouses) error
	UpdateWarehouses(id int, data *model.Warehouses) error
	DeleteWarehouses(id int) error
}

type warehousesService struct {
	Repo repository.WarehousesRepository
}

func NewWarehousesService(repo repository.WarehousesRepository) WarehousesService {
	return &warehousesService{Repo: repo}
}

func (s *warehousesService) GetWarehousesById(id int) (*model.Warehouses, error) {
	return s.Repo.GetWarehousesById(id)
}

func (s *warehousesService) GetAllWarehouses(page, limit int) ([]model.Warehouses, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	
	return s.Repo.GetAllWarehouses(page, limit)
}

func (s *warehousesService) CreateWarehouses(data *model.Warehouses) error {
	return s.Repo.CreateWarehouses(data)
}

func (s *warehousesService) UpdateWarehouses(id int, data *model.Warehouses) error {
	return s.Repo.UpdateWarehouses(id, data)
}

func (s *warehousesService) DeleteWarehouses(id int) error {
	return s.Repo.DeleteWarehouses(id)
}
