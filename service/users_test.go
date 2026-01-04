package service

import (
	"errors"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsersRepository for testing
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) GetUsersByEmail(email string) (*model.Users, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockUsersRepository) CreateUsers(data *model.Users) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockUsersRepository) GetAllUsers() ([]model.Users, error) {
	args := m.Called()
	return args.Get(0).([]model.Users), args.Error(1)
}

func (m *MockUsersRepository) GetUsersByID(id int) (model.Users, error) {
	args := m.Called(id)
	return args.Get(0).(model.Users), args.Error(1)
}

func (m *MockUsersRepository) UpdateUsers(id int, data *model.Users) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockUsersRepository) DeleteUsers(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUsersService_GetUsersByEmail_Success(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	now := time.Now()
	expected := &model.Users{
		Id:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetUsersByEmail", "test@example.com").Return(expected, nil)

	result, err := service.GetUsersByEmail("test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_GetUsersByEmail_NotFound(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	mockRepo.On("GetUsersByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

	result, err := service.GetUsersByEmail("notfound@example.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_GetAllUsers_Success(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	now := time.Now()
	users := []model.Users{
		{Id: 1, Username: "user1", Email: "user1@example.com", Role: "admin", CreatedAt: now, UpdatedAt: now},
		{Id: 2, Username: "user2", Email: "user2@example.com", Role: "user", CreatedAt: now, UpdatedAt: now},
	}

	mockRepo.On("GetAllUsers").Return(users, nil)

	result, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_GetUsersByID_Success(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	now := time.Now()
	expectedUser := model.Users{
		Id:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetUsersByID", 1).Return(expectedUser, nil)

	result, err := service.GetUsersByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Username, result.Username)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_CreateUsers_Success(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	user := &model.Users{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "hashedpassword",
		Role:     "user",
	}

	mockRepo.On("CreateUsers", user).Return(nil)

	err := service.CreateUsers(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_UpdateUsers_Success(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	user := &model.Users{
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword",
		Role:     "admin",
	}

	mockRepo.On("UpdateUsers", 1, user).Return(nil)

	err := service.UpdateUsers(1, user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_DeleteUsers_Success(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	mockRepo.On("DeleteUsers", 1).Return(nil)

	err := service.DeleteUsers(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUsersService_DeleteUsers_Error(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	service := NewUsersService(mockRepo)

	mockRepo.On("DeleteUsers", 999).Return(errors.New("user not found"))

	err := service.DeleteUsers(999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
