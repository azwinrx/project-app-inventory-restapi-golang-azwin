package service

import (
	"errors"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWarehousesRepository for testing
type MockWarehousesRepository struct {
	mock.Mock
}

func (m *MockWarehousesRepository) GetWarehousesById(id int) (*model.Warehouses, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Warehouses), args.Error(1)
}

func (m *MockWarehousesRepository) GetAllWarehouses(page, limit int) ([]model.Warehouses, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Warehouses), args.Int(1), args.Error(2)
}

func (m *MockWarehousesRepository) CreateWarehouses(data *model.Warehouses) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockWarehousesRepository) UpdateWarehouses(id int, data *model.Warehouses) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockWarehousesRepository) DeleteWarehouses(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestWarehousesService_GetWarehousesById_Success(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	now := time.Now()
	expected := &model.Warehouses{
		Id:        1,
		Name:      "Main Warehouse",
		Location:  "Jakarta",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetWarehousesById", 1).Return(expected, nil)

	result, err := service.GetWarehousesById(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_GetAllWarehouses_Success(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	now := time.Now()
	warehouses := []model.Warehouses{
		{Id: 1, Name: "Main Warehouse", Location: "Jakarta", CreatedAt: now, UpdatedAt: now},
		{Id: 2, Name: "Secondary Warehouse", Location: "Bandung", CreatedAt: now, UpdatedAt: now},
	}

	mockRepo.On("GetAllWarehouses", 1, 10).Return(warehouses, 2, nil)

	result, total, err := service.GetAllWarehouses(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_GetAllWarehouses_ValidationPage(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	warehouses := []model.Warehouses{}
	mockRepo.On("GetAllWarehouses", 1, 10).Return(warehouses, 0, nil)

	result, total, err := service.GetAllWarehouses(0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_GetAllWarehouses_ValidationLimit(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	warehouses := []model.Warehouses{}
	mockRepo.On("GetAllWarehouses", 1, 100).Return(warehouses, 0, nil)

	result, total, err := service.GetAllWarehouses(1, 150)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_CreateWarehouses_Success(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	warehouse := &model.Warehouses{
		Name:     "New Warehouse",
		Location: "Surabaya",
	}

	mockRepo.On("CreateWarehouses", warehouse).Return(nil)

	err := service.CreateWarehouses(warehouse)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_UpdateWarehouses_Success(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	warehouse := &model.Warehouses{
		Name:     "Updated Warehouse",
		Location: "Medan",
	}

	mockRepo.On("UpdateWarehouses", 1, warehouse).Return(nil)

	err := service.UpdateWarehouses(1, warehouse)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_DeleteWarehouses_Success(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	mockRepo.On("DeleteWarehouses", 1).Return(nil)

	err := service.DeleteWarehouses(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWarehousesService_DeleteWarehouses_Error(t *testing.T) {
	mockRepo := new(MockWarehousesRepository)
	service := NewWarehousesService(mockRepo)

	mockRepo.On("DeleteWarehouses", 999).Return(errors.New("warehouse not found"))

	err := service.DeleteWarehouses(999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
