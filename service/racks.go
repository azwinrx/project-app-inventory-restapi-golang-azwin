package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type RacksService interface {
	GetRacksById(id int) (*model.Racks, error)
	GetAllRacks(page, limit int) ([]model.Racks, int, error)
	CreateRacks(data *model.Racks) error
	UpdateRacks(id int, data *model.Racks) error
	DeleteRacks(id int) error
}

type racksService struct {
	Repo repository.RacksRepository
}

func NewRacksService(repo repository.RacksRepository) RacksService {
	return &racksService{Repo: repo}
}

func (s *racksService) GetRacksById(id int) (*model.Racks, error) {
	return s.Repo.GetRacksById(id)
}

func (s *racksService) GetAllRacks(page, limit int) ([]model.Racks, int, error) {
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
	
	return s.Repo.GetAllRacks(page, limit)
}

func (s *racksService) CreateRacks(data *model.Racks) error {
	return s.Repo.CreateRacks(data)
}

func (s *racksService) UpdateRacks(id int, data *model.Racks) error {
	return s.Repo.UpdateRacks(id, data)
}

func (s *racksService) DeleteRacks(id int) error {
	return s.Repo.DeleteRacks(id)
}
