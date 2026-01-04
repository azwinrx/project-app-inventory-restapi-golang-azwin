package repository

import (
	"database/sql"
	"errors"
	"project-app-inventory-restapi-golang-azwin/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestGetItemsById_Success(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

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
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		// Fill the scan arguments with test data
		dest := args.Get(0).([]any)
		*dest[0].(*int) = expectedItem.Id
		*dest[1].(*int) = expectedItem.CategoryId
		*dest[2].(*int) = expectedItem.RackId
		*dest[3].(*string) = expectedItem.Name
		*dest[4].(*string) = expectedItem.Sku
		*dest[5].(*int) = expectedItem.Stock
		*dest[6].(*int) = expectedItem.MinStock
		*dest[7].(*float64) = expectedItem.Price
		*dest[8].(*time.Time) = expectedItem.CreatedAt
		*dest[9].(*time.Time) = expectedItem.UpdatedAt
	}).Return(nil)

	// Execute
	result, err := repo.GetItemsById(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedItem.Id, result.Id)
	assert.Equal(t, expectedItem.Name, result.Name)
	assert.Equal(t, expectedItem.Sku, result.Sku)
	mockDB.AssertExpectations(t)
}

func TestGetItemsById_NotFound(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	// Mock expectations - simulate no rows found
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(sql.ErrNoRows)

	// Execute
	result, err := repo.GetItemsById(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	// Note: pgx.ErrNoRows is handled in the repository
	mockDB.AssertExpectations(t)
}

func TestGetAllItems_Success(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockCountRow := new(MockRow)
	mockRows := new(MockRows)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	totalCount := 25
	page := 1
	limit := 10

	// Mock count query - need to match exact args (ctx, query, args)
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockCountRow)
	mockCountRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = totalCount
	}).Return(nil)

	// Mock data query
	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
	
	// Simulate 2 rows returned
	call1 := mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 1
		*dest[1].(*int) = 1
		*dest[2].(*int) = 1
		*dest[3].(*string) = "Item 1"
		*dest[4].(*string) = "SKU-001"
		*dest[5].(*int) = 50
		*dest[6].(*int) = 10
		*dest[7].(*float64) = 10000.0
		*dest[8].(*time.Time) = time.Now()
		*dest[9].(*time.Time) = time.Now()
	}).Return(nil).Once()
	
	call2 := mockRows.On("Next").Return(true).Once().NotBefore(call1)
	mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 2
		*dest[1].(*int) = 1
		*dest[2].(*int) = 1
		*dest[3].(*string) = "Item 2"
		*dest[4].(*string) = "SKU-002"
		*dest[5].(*int) = 75
		*dest[6].(*int) = 15
		*dest[7].(*float64) = 15000.0
		*dest[8].(*time.Time) = time.Now()
		*dest[9].(*time.Time) = time.Now()
	}).Return(nil).Once().NotBefore(call2)
	
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return()

	// Execute
	items, total, err := repo.GetAllItems(page, limit)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, totalCount, total)
	assert.Len(t, items, 2)
	assert.Equal(t, "Item 1", items[0].Name)
	assert.Equal(t, "Item 2", items[1].Name)
	mockDB.AssertExpectations(t)
}

func TestGetLowStockItems_Success(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockRows := new(MockRows)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	threshold := 20

	// Mock query
	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
	
	// Simulate 1 row returned
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 1
		*dest[1].(*int) = 1
		*dest[2].(*int) = 1
		*dest[3].(*string) = "Low Stock Item"
		*dest[4].(*string) = "LOW-001"
		*dest[5].(*int) = 5
		*dest[6].(*int) = 10
		*dest[7].(*float64) = 10000.0
		*dest[8].(*time.Time) = time.Now()
		*dest[9].(*time.Time) = time.Now()
	}).Return(nil).Once()
	
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return()
	mockRows.On("Err").Return(nil)

	// Execute
	items, err := repo.GetLowStockItems(threshold)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Low Stock Item", items[0].Name)
	assert.Equal(t, 5, items[0].Stock)
	mockDB.AssertExpectations(t)
}

func TestCreateItems_Success(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

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
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 1 // Return ID
	}).Return(nil)

	// Execute
	err := repo.CreateItems(newItem)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, newItem.Id)
	mockDB.AssertExpectations(t)
}

func TestUpdateItems_Success(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	updateItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "Updated Item",
		Sku:        "UPD-001",
		Stock:      150,
		MinStock:   15,
		Price:      30000,
	}

	mockTag := MockCommandTag{rowsAffected: 1}

	// Mock expectations
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	// Execute
	err := repo.UpdateItems(1, updateItem)

	// Assert
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateItems_NoRowsAffected(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	updateItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "Updated Item",
		Sku:        "UPD-001",
		Stock:      150,
		MinStock:   15,
		Price:      30000,
	}

	mockTag := MockCommandTag{rowsAffected: 0}

	// Mock expectations
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	// Execute
	err := repo.UpdateItems(999, updateItem)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "no rows affected", err.Error())
	mockDB.AssertExpectations(t)
}

func TestDeleteItems_Success(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 1}

	// Mock expectations
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	// Execute
	err := repo.DeleteItems(1)

	// Assert
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestDeleteItems_NoRowsAffected(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 0}

	// Mock expectations
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	// Execute
	err := repo.DeleteItems(999)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "no rows affected", err.Error())
	mockDB.AssertExpectations(t)
}

func TestGetAllItems_CountQueryError(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	// Mock count query error
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("database error"))

	// Execute
	items, total, err := repo.GetAllItems(1, 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, 0, total)
	mockDB.AssertExpectations(t)
}

func TestCreateItems_Error(t *testing.T) {
	// Setup
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewItemsRepository(mockDB, logger)

	newItem := &model.Items{
		CategoryId: 1,
		RackId:     1,
		Name:       "New Item",
		Sku:        "NEW-001",
		Stock:      100,
		MinStock:   10,
		Price:      25000,
	}

	// Mock expectations - simulate database error
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("duplicate key error"))

	// Execute
	err := repo.CreateItems(newItem)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "duplicate key error", err.Error())
	mockDB.AssertExpectations(t)
}
