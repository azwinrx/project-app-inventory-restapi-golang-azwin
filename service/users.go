package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type UsersService interface {
	FindUsersByEmail(email string) (*model.Users, error)
	CreateUsers(data *model.Users) error
	FindAllUsers() ([]model.Users, error)
	GetUsersByID(id int) (model.Users, error)
	DeleteUsers(id int) error
}

type usersServiceImpl struct {
	Repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersServiceImpl{Repo: repo}
}

func (s *usersServiceImpl) FindUsersByEmail(email string) (*model.Users, error) {
	return s.Repo.FindUsersByEmail(email)
}

func (s *usersServiceImpl) CreateUsers(data *model.Users) error {
	return s.Repo.CreateUsers(data)
}

func (s *usersServiceImpl) FindAllUsers() ([]model.Users, error) {
	return s.Repo.FindAllUsers()
}

func (s *usersServiceImpl) GetUsersByID(id int) (model.Users, error) {
	return s.Repo.GetUsersByID(id)
}

func (s *usersServiceImpl) DeleteUsers(id int) error {
	return s.Repo.DeleteUsers(id)
}