package service

import (
	"errors"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockItemsRepository is a mock implementation of the ItemsRepository interface
type MockItemsRepository struct {
	mock.Mock
}

func (m *MockItemsRepository) GetItemsById(id int) (*model.Items, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Items), args.Error(1)
}

func (m *MockItemsRepository) GetAllItems(page, limit int) ([]model.Items, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Items), args.Int(1), args.Error(2)
}

func (m *MockItemsRepository) GetLowStockItems(threshold int) ([]model.Items, error) {
	args := m.Called(threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Items), args.Error(1)
}

func (m *MockItemsRepository) CreateItems(data *model.Items) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockItemsRepository) UpdateItems(id int, data *model.Items) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockItemsRepository) DeleteItems(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetItemsById_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItem := &model.Items{
		Id:         1,
		CategoryId: 1,
		RackId:     1,
		Name:       "Test Item",
		Sku:        "TEST-001",
		Stock:      100,
		MinStock:   10,
		Price:      50000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetItemsById", 1).Return(expectedItem, nil)

	// Execute
	result, err := service.GetItemsById(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedItem.Id, result.Id)
	assert.Equal(t, expectedItem.Name, result.Name)
	assert.Equal(t, expectedItem.Sku, result.Sku)
	mockRepo.AssertExpectations(t)
}

func TestGetItemsById_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	// Mock expectations
	mockRepo.On("GetItemsById", 999).Return(nil, errors.New("item not found"))

	// Execute
	result, err := service.GetItemsById(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetAllItems_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{
			Id:        1,
			Name:      "Item 1",
			Sku:       "SKU-001",
			Stock:     100,
			MinStock:  10,
			Price:     10000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        2,
			Name:      "Item 2",
			Sku:       "SKU-002",
			Stock:     200,
			MinStock:  20,
			Price:     20000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock expectations
	mockRepo.On("GetAllItems", 1, 10).Return(expectedItems, 25, nil)

	// Execute
	items, total, err := service.GetAllItems(1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 25, total)
	assert.Len(t, items, 2)
	assert.Equal(t, "Item 1", items[0].Name)
	assert.Equal(t, "Item 2", items[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetAllItems_WithInvalidPage(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{Id: 1, Name: "Item 1"},
	}

	// Mock expectations - should default to page 1
	mockRepo.On("GetAllItems", 1, 10).Return(expectedItems, 10, nil)

	// Execute with invalid page (0)
	items, total, err := service.GetAllItems(0, 10)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetAllItems_WithInvalidLimit(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{Id: 1, Name: "Item 1"},
	}

	// Mock expectations - should default to limit 10
	mockRepo.On("GetAllItems", 1, 10).Return(expectedItems, 10, nil)

	// Execute with invalid limit (0)
	items, total, err := service.GetAllItems(1, 0)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetAllItems_WithLimitExceedsMax(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{Id: 1, Name: "Item 1"},
	}

	// Mock expectations - should cap at limit 100
	mockRepo.On("GetAllItems", 1, 100).Return(expectedItems, 10, nil)

	// Execute with limit exceeding max (150)
	items, total, err := service.GetAllItems(1, 150)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetLowStockItems_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{
			Id:        1,
			Name:      "Low Stock Item",
			Sku:       "LOW-001",
			Stock:     3,
			MinStock:  10,
			Price:     10000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock expectations
	mockRepo.On("GetLowStockItems", 10).Return(expectedItems, nil)

	// Execute
	items, err := service.GetLowStockItems(10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Low Stock Item", items[0].Name)
	assert.Equal(t, 3, items[0].Stock)
	mockRepo.AssertExpectations(t)
}

func TestGetLowStockItems_WithInvalidThreshold(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{Id: 1, Name: "Low Stock Item", Stock: 3},
	}

	// Mock expectations - should default to threshold 5
	mockRepo.On("GetLowStockItems", 5).Return(expectedItems, nil)

	// Execute with invalid threshold (0)
	items, err := service.GetLowStockItems(0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetLowStockItems_WithNegativeThreshold(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	expectedItems := []model.Items{
		{Id: 1, Name: "Low Stock Item", Stock: 3},
	}

	// Mock expectations - should default to threshold 5
	mockRepo.On("GetLowStockItems", 5).Return(expectedItems, nil)

	// Execute with negative threshold
	items, err := service.GetLowStockItems(-10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestCreateItems_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	newItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "New Item",
		Sku:        "NEW-001",
		Stock:      100,
		MinStock:   10,
		Price:      25000,
	}

	// Mock expectations
	mockRepo.On("CreateItems", newItem).Return(nil)

	// Execute
	err := service.CreateItems(newItem)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateItems_Error(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	newItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "New Item",
		Sku:        "NEW-001",
		Stock:      100,
		MinStock:   10,
		Price:      25000,
	}

	// Mock expectations
	mockRepo.On("CreateItems", newItem).Return(errors.New("database error"))

	// Execute
	err := service.CreateItems(newItem)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateItems_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	updateItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "Updated Item",
		Sku:        "UPD-001",
		Stock:      150,
		MinStock:   15,
		Price:      30000,
	}

	// Mock expectations
	mockRepo.On("UpdateItems", 1, updateItem).Return(nil)

	// Execute
	err := service.UpdateItems(1, updateItem)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateItems_Error(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	updateItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "Updated Item",
		Sku:        "UPD-001",
		Stock:      150,
		MinStock:   15,
		Price:      30000,
	}

	// Mock expectations
	mockRepo.On("UpdateItems", 999, updateItem).Return(errors.New("item not found"))

	// Execute
	err := service.UpdateItems(999, updateItem)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteItems_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	// Mock expectations
	mockRepo.On("DeleteItems", 1).Return(nil)

	// Execute
	err := service.DeleteItems(1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteItems_Error(t *testing.T) {
	// Setup
	mockRepo := new(MockItemsRepository)
	service := NewItemsService(mockRepo)

	// Mock expectations
	mockRepo.On("DeleteItems", 999).Return(errors.New("item not found"))

	// Execute
	err := service.DeleteItems(999)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	mockRepo.AssertExpectations(t)
}
