package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"

	"go.uber.org/zap"
)

type RacksRepository interface {
	GetRacksById(id int) (*model.Racks, error)
	GetAllRacks(page, limit int) ([]model.Racks, int, error)
	CreateRacks(data *model.Racks) error
	UpdateRacks(id int, data *model.Racks) error
	DeleteRacks(id int) error
}

type racksRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewRacksRepository(db database.PgxIface, log *zap.Logger) RacksRepository {
	return &racksRepository{db: db, Logger: log}
}

func (r *racksRepository) GetRacksById(id int) (*model.Racks, error) {
	query := `
		SELECT id, warehouse_id, name, created_at, updated_at
		FROM racks
		WHERE id = $1
	`
	var rack model.Racks
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&rack.Id,
		&rack.WarehouseId,
		&rack.Name,
		&rack.CreatedAt,
		&rack.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("rack not found")
	}
	return &rack, err
}

func (r *racksRepository) GetAllRacks(page, limit int) ([]model.Racks, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM racks`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findall repo ", zap.Error(err))
		return nil, 0, err
	}

	// get data with pagination
	query := `
		SELECT id, warehouse_id, name, created_at, updated_at
		FROM racks
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var racks []model.Racks
	for rows.Next() {
		var rack model.Racks
		err := rows.Scan(
			&rack.Id,
			&rack.WarehouseId,
			&rack.Name,
			&rack.CreatedAt,
			&rack.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		racks = append(racks, rack)
	}

	return racks, total, nil
}

func (r *racksRepository) CreateRacks(data *model.Racks) error {
	query := `
		INSERT INTO racks (warehouse_id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, data.WarehouseId, data.Name).Scan(&data.Id)
	return err
}

func (r *racksRepository) UpdateRacks(id int, data *model.Racks) error {
	query := `
		UPDATE racks
		SET warehouse_id = $1, name = $2, updated_at = NOW()
		WHERE id = $3`

	result, err := r.db.Exec(context.Background(), query, data.WarehouseId, data.Name, id)
	if err != nil {
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("no rows affected")
	}
	return err
}

func (r *racksRepository) DeleteRacks(id int) error {
	query := `DELETE FROM racks WHERE id = $1`

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
