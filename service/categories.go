package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type CategoriesService interface {
	GetCategoriesById(id int) (*model.Categories, error)
	GetAllCategories(page, limit int) ([]model.Categories, int, error)
	CreateCategories(data *model.Categories) error
	UpdateCategories(id int, data *model.Categories) error
	DeleteCategories(id int) error
}

type categoriesService struct {
	Repo repository.CategoriesRepository
}

func NewCategoriesService(repo repository.CategoriesRepository) CategoriesService {
	return &categoriesService{Repo: repo}
}

func (s *categoriesService) GetCategoriesById(id int) (*model.Categories, error) {
	return s.Repo.GetCategoriesById(id)
}

func (s *categoriesService) GetAllCategories(page, limit int) ([]model.Categories, int, error) {
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
	
	return s.Repo.GetAllCategories(page, limit)
}

func (s *categoriesService) CreateCategories(data *model.Categories) error {
	return s.Repo.CreateCategories(data)
}

func (s *categoriesService) UpdateCategories(id int, data *model.Categories) error {
	return s.Repo.UpdateCategories(id, data)
}

func (s *categoriesService) DeleteCategories(id int) error {
	return s.Repo.DeleteCategories(id)
}