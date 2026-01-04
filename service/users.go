package service

import (
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/repository"
)

type UsersService interface {
	GetUsersByEmail(email string) (*model.Users, error)
	GetUsersByID(id int) (model.Users, error)
	GetAllUsers()([]model.Users, error)
	CreateUsers(data *model.Users) error
	UpdateUsers(id int, data *model.Users) error
	DeleteUsers(id int) error
}

type usersServiceImpl struct {
	Repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersServiceImpl{Repo: repo}
}

func (s *usersServiceImpl) GetUsersByEmail(email string) (*model.Users, error) {
	return s.Repo.GetUsersByEmail(email)
}

func (s *usersServiceImpl) GetAllUsers() ([]model.Users, error) {
	return s.Repo.GetAllUsers()
}

func (s *usersServiceImpl) GetUsersByID(id int) (model.Users, error) {
	return s.Repo.GetUsersByID(id)
}

func (s *usersServiceImpl) CreateUsers(data *model.Users) error {
	return s.Repo.CreateUsers(data)
}

func (s *usersServiceImpl) UpdateUsers(id int, data *model.Users) error {
	return s.Repo.UpdateUsers(id, data)
}

func (s *usersServiceImpl) DeleteUsers(id int) error {
	return s.Repo.DeleteUsers(id)
}