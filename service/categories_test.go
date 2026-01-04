package service

import (
	"errors"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoriesRepository for testing
type MockCategoriesRepository struct {
	mock.Mock
}

func (m *MockCategoriesRepository) GetCategoriesById(id int) (*model.Categories, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Categories), args.Error(1)
}

func (m *MockCategoriesRepository) GetAllCategories(page, limit int) ([]model.Categories, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Categories), args.Int(1), args.Error(2)
}

func (m *MockCategoriesRepository) CreateCategories(data *model.Categories) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockCategoriesRepository) UpdateCategories(id int, data *model.Categories) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockCategoriesRepository) DeleteCategories(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCategoriesService_GetCategoriesById_Success(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	now := time.Now()
	expected := &model.Categories{
		Id:        1,
		Name:      "Electronics",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetCategoriesById", 1).Return(expected, nil)

	result, err := service.GetCategoriesById(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_GetCategoriesById_Error(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	mockRepo.On("GetCategoriesById", 999).Return(nil, errors.New("category not found"))

	result, err := service.GetCategoriesById(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_GetAllCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	now := time.Now()
	categories := []model.Categories{
		{Id: 1, Name: "Electronics", CreatedAt: now, UpdatedAt: now},
		{Id: 2, Name: "Furniture", CreatedAt: now, UpdatedAt: now},
	}

	mockRepo.On("GetAllCategories", 1, 10).Return(categories, 2, nil)

	result, total, err := service.GetAllCategories(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_GetAllCategories_ValidationPage(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	categories := []model.Categories{}
	mockRepo.On("GetAllCategories", 1, 10).Return(categories, 0, nil)

	// Test with invalid page (should default to 1)
	result, total, err := service.GetAllCategories(0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_GetAllCategories_ValidationLimit(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	categories := []model.Categories{}
	
	// Test with limit > 100 (should cap at 100)
	mockRepo.On("GetAllCategories", 1, 100).Return(categories, 0, nil)
	result, total, err := service.GetAllCategories(1, 150)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_GetAllCategories_ValidationLimitMin(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	categories := []model.Categories{}
	
	// Test with limit < 1 (should default to 10)
	mockRepo.On("GetAllCategories", 1, 10).Return(categories, 0, nil)
	result, total, err := service.GetAllCategories(1, 0)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_CreateCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	category := &model.Categories{
		Name: "New Category",
	}

	mockRepo.On("CreateCategories", category).Return(nil)

	err := service.CreateCategories(category)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_CreateCategories_Error(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	category := &model.Categories{
		Name: "New Category",
	}

	mockRepo.On("CreateCategories", category).Return(errors.New("database error"))

	err := service.CreateCategories(category)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_UpdateCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	category := &model.Categories{
		Name: "Updated Category",
	}

	mockRepo.On("UpdateCategories", 1, category).Return(nil)

	err := service.UpdateCategories(1, category)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoriesService_DeleteCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoriesRepository)
	service := NewCategoriesService(mockRepo)

	mockRepo.On("DeleteCategories", 1).Return(nil)

	err := service.DeleteCategories(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
