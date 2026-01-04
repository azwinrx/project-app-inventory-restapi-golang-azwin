package service

import (
	"project-app-inventory-restapi-golang-azwin/dto"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSalesRepository for testing
type MockSalesRepository struct {
	mock.Mock
}

func (m *MockSalesRepository) GetSalesById(id int) (*model.Sales, []model.SaleItems, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*model.Sales), args.Get(1).([]model.SaleItems), args.Error(2)
}

func (m *MockSalesRepository) GetAllSales(page, limit int) ([]model.Sales, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Sales), args.Int(1), args.Error(2)
}

func (m *MockSalesRepository) CreateSales(sale *model.Sales, items []model.SaleItems) error {
	args := m.Called(sale, items)
	return args.Error(0)
}

func (m *MockSalesRepository) UpdateSales(id int, data *model.Sales) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockSalesRepository) DeleteSales(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSalesService_GetSalesById_Success(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	now := time.Now()
	sale := &model.Sales{
		Id:          1,
		UserId:      1,
		TotalAmount: 100.50,
		CreatedAt:   now,
	}
	items := []model.SaleItems{
		{Id: 1, SaleId: 1, ItemId: 1, Quantity: 2, Price: 50.25, Subtotal: 100.50},
	}

	mockRepo.On("GetSalesById", 1).Return(sale, items, nil)

	result, err := service.GetSalesById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sale.Id, result.Id)
	assert.Equal(t, sale.TotalAmount, result.TotalAmount)
	assert.Len(t, result.Items, 1)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_GetSalesById_NotFound(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	mockRepo.On("GetSalesById", 999).Return(nil, nil, assert.AnError)

	result, err := service.GetSalesById(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_GetAllSales_Success(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	now := time.Now()
	sales := []model.Sales{
		{
			Id:          1,
			UserId:      1,
			TotalAmount: 100.50,
			CreatedAt:   now,
			Items: []model.SaleItems{
				{Id: 1, SaleId: 1, ItemId: 1, Quantity: 2, Price: 50.25, Subtotal: 100.50},
			},
		},
	}

	mockRepo.On("GetAllSales", 1, 10).Return(sales, 1, nil)

	result, total, err := service.GetAllSales(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].Items, 1)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_GetAllSales_ValidationPage(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	sales := []model.Sales{}
	mockRepo.On("GetAllSales", 1, 10).Return(sales, 0, nil)

	result, total, err := service.GetAllSales(0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_GetAllSales_ValidationLimit(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	sales := []model.Sales{}
	mockRepo.On("GetAllSales", 1, 100).Return(sales, 0, nil)

	result, total, err := service.GetAllSales(1, 150)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_CreateSales_Success(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 1,
		Items: []dto.SaleItemRequest{
			{ItemId: 1, Quantity: 2, Price: 50.25},
		},
	}

	mockRepo.On("CreateSales", mock.AnythingOfType("*model.Sales"), mock.AnythingOfType("[]model.SaleItems")).Return(nil)

	err := service.CreateSales(request)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_CreateSales_ValidationUserIdRequired(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 0, // Invalid
		Items: []dto.SaleItemRequest{
			{ItemId: 1, Quantity: 2, Price: 50.25},
		},
	}

	err := service.CreateSales(request)

	assert.Error(t, err)
	assert.Equal(t, "user_id is required", err.Error())
}

func TestSalesService_CreateSales_ValidationItemsRequired(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 1,
		Items:  []dto.SaleItemRequest{}, // Empty
	}

	err := service.CreateSales(request)

	assert.Error(t, err)
	assert.Equal(t, "at least one item is required", err.Error())
}

func TestSalesService_CreateSales_ValidationQuantityInvalid(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 1,
		Items: []dto.SaleItemRequest{
			{ItemId: 1, Quantity: 0, Price: 50.25}, // Invalid quantity
		},
	}

	err := service.CreateSales(request)

	assert.Error(t, err)
	assert.Equal(t, "quantity must be greater than 0", err.Error())
}

func TestSalesService_CreateSales_ValidationPriceInvalid(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 1,
		Items: []dto.SaleItemRequest{
			{ItemId: 1, Quantity: 2, Price: 0}, // Invalid price
		},
	}

	err := service.CreateSales(request)

	assert.Error(t, err)
	assert.Equal(t, "price must be greater than 0", err.Error())
}

func TestSalesService_UpdateSales_Success(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 1,
		Items: []dto.SaleItemRequest{
			{ItemId: 1, Quantity: 3, Price: 50.25},
		},
	}

	mockRepo.On("UpdateSales", 1, mock.AnythingOfType("*model.Sales")).Return(nil)

	err := service.UpdateSales(1, request)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_UpdateSales_ValidationUserIdRequired(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	request := &dto.SalesRequest{
		UserId: 0, // Invalid
		Items: []dto.SaleItemRequest{
			{ItemId: 1, Quantity: 2, Price: 50.25},
		},
	}

	err := service.UpdateSales(1, request)

	assert.Error(t, err)
	assert.Equal(t, "user_id is required", err.Error())
}

func TestSalesService_DeleteSales_Success(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	mockRepo.On("DeleteSales", 1).Return(nil)

	err := service.DeleteSales(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSalesService_DeleteSales_Error(t *testing.T) {
	mockRepo := new(MockSalesRepository)
	service := NewSalesService(mockRepo)

	mockRepo.On("DeleteSales", 999).Return(assert.AnError)

	err := service.DeleteSales(999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
