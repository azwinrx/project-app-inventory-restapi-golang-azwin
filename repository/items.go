package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"

	"go.uber.org/zap"
)

type ItemsRepository interface {
	GetItemsById(id int) (*model.Items, error) 
	GetAllItems(page, limit int) ([]model.Items, int, error)
	GetLowStockItems(threshold int) ([]model.Items, error)
	CreateItems(data *model.Items) error
	UpdateItems(id int, data *model.Items) error
	DeleteItems(id int) error
}

type itemsRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewItemsRepository(db database.PgxIface, log *zap.Logger) ItemsRepository {
	return &itemsRepository{db: db, Logger: log}
}

func (r *itemsRepository) GetItemsById(id int) (*model.Items, error) {
	query := `
		SELECT id, category_id, rack_id, name, sku, stock, min_stock, price, created_at, updated_at
		FROM items
		WHERE id = $1

	`
	var i model.Items
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&i.Id,
		&i.CategoryId,
		&i	.RackId,
		&i.Name,
		&i.Sku,
		&i.Stock,
		&i.MinStock,
		&i.Price,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("item not found")
	}
	return &i, err
}

func (r *itemsRepository) GetAllItems(page, limit int) ([]model.Items, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM items`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findall repo ", zap.Error(err))
		return nil, 0, err
	}

	// get data with pagination
	query := `
		SELECT id, category_id, rack_id, name, sku, stock, min_stock, price, created_at, updated_at
		FROM items
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.Items
	for rows.Next() {
		var i model.Items
		err := rows.Scan(
			&i.Id,
			&i.CategoryId,
			&i.RackId,
			&i.Name,
			&i.Sku,
			&i.Stock,
			&i.MinStock,
			&i.Price,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, i)
	}

	return items, total, nil
}

func (r *itemsRepository) GetLowStockItems(threshold int) ([]model.Items, error) {
	query := `
		SELECT id, category_id, rack_id, name, sku, stock, min_stock, price, created_at, updated_at
		FROM items
		WHERE stock < $1
		ORDER BY stock ASC, name ASC
	`

	rows, err := r.db.Query(context.Background(), query, threshold)
	if err != nil {
		r.Logger.Error("failed to get low stock items",
			zap.Int("threshold", threshold),
			zap.Error(err),
		)
		return nil, err
	}
	defer rows.Close()

	var items []model.Items
	for rows.Next() {
		var i model.Items
		err := rows.Scan(
			&i.Id,
			&i.CategoryId,
			&i.RackId,
			&i.Name,
			&i.Sku,
			&i.Stock,
			&i.MinStock,
			&i.Price,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			r.Logger.Error("failed to scan low stock item", zap.Error(err))
			return nil, err
		}
		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		r.Logger.Error("error iterating low stock items", zap.Error(err))
		return nil, err
	}

	r.Logger.Info("low stock items retrieved",
		zap.Int("count", len(items)),
		zap.Int("threshold", threshold),
	)

	return items, nil
}

func (r *itemsRepository) CreateItems(data *model.Items) error {
	query := `
		INSERT INTO items (category_id, rack_id, name, sku, stock, min_stock, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, data.CategoryId, data.RackId, data.Name, data.Sku, data.Stock, data.MinStock, data.Price).Scan(&data.Id)
	return err
}

func (r *itemsRepository) UpdateItems(id int, data *model.Items) error {
	query := `
		UPDATE items
		SET category_id = $1, rack_id = $2, name = $3, sku = $4, stock = $5, min_stock = $6, price = $7, updated_at = NOW()
		WHERE id = $8`

	result, err := r.db.Exec(context.Background(), query, data.CategoryId, data.RackId, data.Name, data.Sku, data.Stock, data.MinStock, data.Price, id)
	if err != nil {
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("no rows affected")
	}
	return err
}

func (r *itemsRepository) DeleteItems(id int) error {
	query := `DELETE FROM items WHERE id = $1`

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