package service

import (
	"errors"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRacksRepository for testing
type MockRacksRepository struct {
	mock.Mock
}

func (m *MockRacksRepository) GetRacksById(id int) (*model.Racks, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Racks), args.Error(1)
}

func (m *MockRacksRepository) GetAllRacks(page, limit int) ([]model.Racks, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Racks), args.Int(1), args.Error(2)
}

func (m *MockRacksRepository) CreateRacks(data *model.Racks) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockRacksRepository) UpdateRacks(id int, data *model.Racks) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockRacksRepository) DeleteRacks(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRacksService_GetRacksById_Success(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	now := time.Now()
	expected := &model.Racks{
		Id:          1,
		WarehouseId: 1,
		Name:        "Rack A1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetRacksById", 1).Return(expected, nil)

	result, err := service.GetRacksById(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_GetAllRacks_Success(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	now := time.Now()
	racks := []model.Racks{
		{Id: 1, WarehouseId: 1, Name: "Rack A1", CreatedAt: now, UpdatedAt: now},
		{Id: 2, WarehouseId: 1, Name: "Rack A2", CreatedAt: now, UpdatedAt: now},
	}

	mockRepo.On("GetAllRacks", 1, 10).Return(racks, 2, nil)

	result, total, err := service.GetAllRacks(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_GetAllRacks_ValidationPage(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	racks := []model.Racks{}
	mockRepo.On("GetAllRacks", 1, 10).Return(racks, 0, nil)

	result, total, err := service.GetAllRacks(-1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_GetAllRacks_ValidationLimit(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	racks := []model.Racks{}
	mockRepo.On("GetAllRacks", 1, 100).Return(racks, 0, nil)

	result, total, err := service.GetAllRacks(1, 200)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_CreateRacks_Success(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "New Rack",
	}

	mockRepo.On("CreateRacks", rack).Return(nil)

	err := service.CreateRacks(rack)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_UpdateRacks_Success(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "Updated Rack",
	}

	mockRepo.On("UpdateRacks", 1, rack).Return(nil)

	err := service.UpdateRacks(1, rack)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_DeleteRacks_Success(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	mockRepo.On("DeleteRacks", 1).Return(nil)

	err := service.DeleteRacks(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRacksService_DeleteRacks_Error(t *testing.T) {
	mockRepo := new(MockRacksRepository)
	service := NewRacksService(mockRepo)

	mockRepo.On("DeleteRacks", 999).Return(errors.New("rack not found"))

	err := service.DeleteRacks(999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
