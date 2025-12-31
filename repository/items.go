package repository

import (
	"context"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"

	"go.uber.org/zap"
)

type ItemsRepository interface {
	GetAllItems(page, limit int) ([]model.Items, int, error)
}

type itemsRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewItemsRepository(db database.PgxIface, log *zap.Logger) ItemsRepository {
	return &itemsRepository{db: db, Logger: log}
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
