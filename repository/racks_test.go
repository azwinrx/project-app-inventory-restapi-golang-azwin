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

func TestGetRacksById_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	now := time.Now()
	expectedRack := &model.Racks{
		Id:          1,
		WarehouseId: 1,
		Name:        "Rack A1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = expectedRack.Id
		*dest[1].(*int) = expectedRack.WarehouseId
		*dest[2].(*string) = expectedRack.Name
		*dest[3].(*time.Time) = expectedRack.CreatedAt
		*dest[4].(*time.Time) = expectedRack.UpdatedAt
	}).Return(nil)

	result, err := repo.GetRacksById(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedRack.Id, result.Id)
	assert.Equal(t, expectedRack.Name, result.Name)
	mockDB.AssertExpectations(t)
}

func TestGetRacksById_NotFound(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(sql.ErrNoRows)

	result, err := repo.GetRacksById(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "rack not found", err.Error())
}

func TestGetAllRacks_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRowCount := new(MockRow)
	mockRows := new(MockRows)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

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
			*dest[1].(*int) = 1
			*dest[2].(*string) = "Rack A1"
			*dest[3].(*time.Time) = now
			*dest[4].(*time.Time) = now
		} else {
			*dest[0].(*int) = 2
			*dest[1].(*int) = 1
			*dest[2].(*string) = "Rack A2"
			*dest[3].(*time.Time) = now
			*dest[4].(*time.Time) = now
		}
	}).Return(nil)

	mockRows.On("Close").Return()
	mockRows.On("Err").Return(nil)

	racks, total, err := repo.GetAllRacks(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, racks, 2)
	assert.Equal(t, "Rack A1", racks[0].Name)
}

func TestCreateRacks_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "New Rack",
	}

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 1
	}).Return(nil)

	err := repo.CreateRacks(rack)

	assert.NoError(t, err)
	assert.Equal(t, 1, rack.Id)
}

func TestUpdateRacks_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "Updated Rack",
	}

	mockTag := MockCommandTag{rowsAffected: 1}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.UpdateRacks(1, rack)

	assert.NoError(t, err)
}

func TestUpdateRacks_NoRowsAffected(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "Updated Rack",
	}

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.UpdateRacks(999, rack)

	assert.Error(t, err)
	assert.Equal(t, "no rows affected", err.Error())
}

func TestDeleteRacks_Success(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 1}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.DeleteRacks(1)

	assert.NoError(t, err)
}

func TestDeleteRacks_NoRowsAffected(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, nil)

	err := repo.DeleteRacks(999)

	assert.Error(t, err)
	assert.Equal(t, "no rows affected", err.Error())
}

func TestCreateRacks_Error(t *testing.T) {
	mockDB := new(MockPgxIface)
	mockRow := new(MockRow)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "New Rack",
	}

	// Mock database error
	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("foreign key constraint failed"))

	err := repo.CreateRacks(rack)

	assert.Error(t, err)
	assert.Equal(t, "foreign key constraint failed", err.Error())
	mockDB.AssertExpectations(t)
}

func TestUpdateRacks_Error(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	rack := &model.Racks{
		WarehouseId: 1,
		Name:        "Updated Rack",
	}

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, errors.New("database connection error"))

	err := repo.UpdateRacks(1, rack)

	assert.Error(t, err)
	assert.Equal(t, "database connection error", err.Error())
	mockDB.AssertExpectations(t)
}

func TestDeleteRacks_Error(t *testing.T) {
	mockDB := new(MockPgxIface)
	logger, _ := zap.NewDevelopment()
	repo := NewRacksRepository(mockDB, logger)

	mockTag := MockCommandTag{rowsAffected: 0}
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockTag, errors.New("foreign key violation"))

	err := repo.DeleteRacks(1)

	assert.Error(t, err)
	assert.Equal(t, "foreign key violation", err.Error())
	mockDB.AssertExpectations(t)
}
