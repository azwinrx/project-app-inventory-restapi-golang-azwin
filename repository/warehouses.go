package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"

	"go.uber.org/zap"
)

type WarehousesRepository interface {
	GetWarehousesById(id int) (*model.Warehouses, error)
	GetAllWarehouses(page, limit int) ([]model.Warehouses, int, error)
	CreateWarehouses(data *model.Warehouses) error
	UpdateWarehouses(id int, data *model.Warehouses) error
	DeleteWarehouses(id int) error
}

type warehousesRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewWarehousesRepository(db database.PgxIface, log *zap.Logger) WarehousesRepository {
	return &warehousesRepository{db: db, Logger: log}
}

func (r *warehousesRepository) GetWarehousesById(id int) (*model.Warehouses, error) {
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM warehouses
		WHERE id = $1
	`
	var warehouse model.Warehouses
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&warehouse.Id,
		&warehouse.Name,
		&warehouse.Location,
		&warehouse.CreatedAt,
		&warehouse.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("warehouse not found")
	}
	return &warehouse, err
}

func (r *warehousesRepository) GetAllWarehouses(page, limit int) ([]model.Warehouses, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM warehouses`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findall repo ", zap.Error(err))
		return nil, 0, err
	}

	// get data with pagination
	query := `
		SELECT id, name, location, created_at, updated_at
		FROM warehouses
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var warehouses []model.Warehouses
	for rows.Next() {
		var warehouse model.Warehouses
		err := rows.Scan(
			&warehouse.Id,
			&warehouse.Name,
			&warehouse.Location,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, total, nil
}

func (r *warehousesRepository) CreateWarehouses(data *model.Warehouses) error {
	query := `
		INSERT INTO warehouses (name, location, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, data.Name, data.Location).Scan(&data.Id)
	return err
}

func (r *warehousesRepository) UpdateWarehouses(id int, data *model.Warehouses) error {
	query := `
		UPDATE warehouses
		SET name = $1, location = $2, updated_at = NOW()
		WHERE id = $3`

	result, err := r.db.Exec(context.Background(), query, data.Name, data.Location, id)
	if err != nil {
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("no rows affected")
	}
	return err
}

func (r *warehousesRepository) DeleteWarehouses(id int) error {
	query := `DELETE FROM warehouses WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("no rows affected")
	}

	return err
}
