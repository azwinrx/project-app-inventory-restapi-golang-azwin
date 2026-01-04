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

func TestGetCategoriesById_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	now := time.Now()
	expectedCategory := &model.Categories{
		Id:        1,
		Name:      "Electronics",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = expectedCategory.Id
		*dest[1].(*string) = expectedCategory.Name
		*dest[2].(*time.Time) = expectedCategory.CreatedAt
		*dest[3].(*time.Time) = expectedCategory.UpdatedAt
	}).Return(nil)

	result, err := repo.GetCategoriesById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCategory.Id, result.Id)
	assert.Equal(t, expectedCategory.Name, result.Name)
	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestGetCategoriesById_NotFound(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(sql.ErrNoRows)

	result, err := repo.GetCategoriesById(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "category not found", err.Error())
	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestGetAllCategories_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRowCount := new(MockRow)
	mockRows := new(MockRows)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	now := time.Now()
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowCount)
	mockRowCount.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 2
	}).Return(nil)

	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)

	callCount := 0
	mockRows.On("Next").Return(true).Times(2)
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		callCount++
		dest := args.Get(0).([]any)
		if callCount == 1 {
			*dest[0].(*int) = 1
			*dest[1].(*string) = "Electronics"
			*dest[2].(*time.Time) = now
			*dest[3].(*time.Time) = now
		} else if callCount == 2 {
			*dest[0].(*int) = 2
			*dest[1].(*string) = "Furniture"
			*dest[2].(*time.Time) = now
			*dest[3].(*time.Time) = now
		}
	}).Return(nil)

	mockRows.On("Close").Return()
	mockRows.On("Err").Return(nil)

	categories, total, err := repo.GetAllCategories(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Electronics", categories[0].Name)
	assert.Equal(t, "Furniture", categories[1].Name)
	mockDB.AssertExpectations(t)
}

func TestCreateCategories_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	category := &model.Categories{
		Name: "New Category",
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 1
	}).Return(nil)

	err := repo.CreateCategories(category)

	assert.NoError(t, err)
	assert.Equal(t, 1, category.Id)
	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestUpdateCategories_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	category := &model.Categories{
		Name: "Updated Category",
	}

	mockTag := MockCommandTag{rowsAffected: 1}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.UpdateCategories(1, category)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateCategories_NoRowsAffected(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	category := &model.Categories{
		Name: "Updated Category",
	}

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.UpdateCategories(999, category)

	assert.Error(t, err)
	assert.Equal(t, "no rows affected", err.Error())
	mockDB.AssertExpectations(t)
}

func TestDeleteCategories_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 1}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.DeleteCategories(1)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestDeleteCategories_NoRowsAffected(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.DeleteCategories(999)

	assert.Error(t, err)
	assert.Equal(t, "no rows affected", err.Error())
	mockDB.AssertExpectations(t)
}

func TestDeleteCategories_Error(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, errors.New("database error"))

	err := repo.DeleteCategories(1)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestGetAllCategories_QueryError(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRowCount := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	// Mock count query success
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowCount)
	mockRowCount.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 10
	}).Return(nil)

	// Mock data query failure
	mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("query failed"))

	categories, total, err := repo.GetAllCategories(1, 10)

	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Equal(t, 0, total)
	assert.Equal(t, "query failed", err.Error())
	mockDB.AssertExpectations(t)
}

func TestCreateCategories_Error(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	category := &model.Categories{
		Name: "New Category",
	}

	// Mock database error during creation
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("constraint violation"))

	err := repo.CreateCategories(category)

	assert.Error(t, err)
	assert.Equal(t, "constraint violation", err.Error())
	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestUpdateCategories_Error(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	category := &model.Categories{
		Name: "Updated Category",
	}

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, errors.New("connection lost"))

	err := repo.UpdateCategories(1, category)

	assert.Error(t, err)
	assert.Equal(t, "connection lost", err.Error())
	mockDB.AssertExpectations(t)
}

func TestGetAllCategories_CountQueryError(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRowCount := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewCategoriesRepository(mockDB, logger)

	// Mock count query error
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRowCount)
	mockRowCount.On("Scan", mock.Anything).Return(errors.New("count query failed"))

	categories, total, err := repo.GetAllCategories(1, 10)

	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Equal(t, 0, total)
	assert.Equal(t, "count query failed", err.Error())
	mockDB.AssertExpectations(t)
}
