package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"
	"strings"

	"go.uber.org/zap"
)

type SalesRepository interface {
	GetSalesById(id int) (*model.Sales, []model.SaleItems, error)
	GetAllSales(page, limit int) ([]model.Sales, int, error)
	CreateSales(sale *model.Sales, items []model.SaleItems) error
	UpdateSales(id int, data *model.Sales) error
	DeleteSales(id int) error
}

type salesRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewSalesRepository(db database.PgxIface, log *zap.Logger) SalesRepository {
	return &salesRepository{db: db, Logger: log}
}

func (r *salesRepository) GetSalesById(id int) (*model.Sales, []model.SaleItems, error) {
	// Get sales data
	queryS := `
		SELECT id, user_id, total_amount, created_at
		FROM sales
		WHERE id = $1
	`
	var s model.Sales
	err := r.db.QueryRow(context.Background(), queryS, id).Scan(
		&s.Id,
		&s.UserId,
		&s.TotalAmount,
		&s.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil, errors.New("sale not found")
	}
	if err != nil {
		return nil, nil, err
	}

	// Get sale items
	queryItems := `
		SELECT id, sale_id, item_id, quantity, price, subtotal
		FROM sale_items
		WHERE sale_id = $1
	`
	rows, err := r.db.Query(context.Background(), queryItems, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var items []model.SaleItems
	for rows.Next() {
		var item model.SaleItems
		err := rows.Scan(
			&item.Id,
			&item.SaleId,
			&item.ItemId,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
		)
		if err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}

	return &s, items, nil
}

func (r *salesRepository) GetAllSales(page, limit int) ([]model.Sales, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM sales`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query count sales", zap.Error(err))
		return nil, 0, err
	}

	// get data with pagination
	query := `
		SELECT id, user_id, total_amount, created_at
		FROM sales
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sales []model.Sales
	for rows.Next() {
		var s model.Sales
		err := rows.Scan(
			&s.Id,
			&s.UserId,
			&s.TotalAmount,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		sales = append(sales, s)
	}

	// Fetch sale items for each sale
	for i := range sales {
		itemsQuery := `
			SELECT id, sale_id, item_id, quantity, price, subtotal
			FROM sale_items
			WHERE sale_id = $1
		`
		itemRows, err := r.db.Query(context.Background(), itemsQuery, sales[i].Id)
		if err != nil {
			r.Logger.Error("error querying sale items", zap.Error(err))
			continue
		}

		var items []model.SaleItems
		for itemRows.Next() {
			var item model.SaleItems
			err := itemRows.Scan(
				&item.Id,
				&item.SaleId,
				&item.ItemId,
				&item.Quantity,
				&item.Price,
				&item.Subtotal,
			)
			if err != nil {
				r.Logger.Error("error scanning sale item", zap.Error(err))
				continue
			}
			items = append(items, item)
		}
		itemRows.Close()
		sales[i].Items = items
	}

	return sales, total, nil
}

func (r *salesRepository) CreateSales(sale *model.Sales, items []model.SaleItems) error {
	// Start Transaction
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		r.Logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
			r.Logger.Error("transaction rolled back", zap.Error(err))
		}
	}()

	// Insert Sales
	querySales := `
		INSERT INTO sales (user_id, total_amount, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id
	`
	var saleId int
	err = tx.QueryRow(context.Background(), querySales,
		sale.UserId,
		sale.TotalAmount,
	).Scan(&saleId)

	if err != nil {
		r.Logger.Error("failed to insert sales", zap.Error(err))
		return err
	}

	// Batch INSERT sale_items
	var valueStrings []string
	var valueArgs []interface{}
	argPosition := 1

	for _, item := range items {
		subtotal := float64(item.Quantity) * item.Price
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				argPosition, argPosition+1, argPosition+2,
				argPosition+3, argPosition+4))

		valueArgs = append(valueArgs, saleId, item.ItemId,
			item.Quantity, item.Price, subtotal)
		argPosition += 5
	}

	querySaleItems := fmt.Sprintf(`
		INSERT INTO sale_items (sale_id, item_id, quantity, price, subtotal)
		VALUES %s
	`, strings.Join(valueStrings, ", "))

	_, err = tx.Exec(context.Background(), querySaleItems, valueArgs...)
	if err != nil {
		r.Logger.Error("failed to batch insert sale items", zap.Error(err))
		return err
	}

	// Batch UPDATE stock
	var itemIds []int
	var quantities []int

	for _, item := range items {
		itemIds = append(itemIds, item.ItemId)
		quantities = append(quantities, item.Quantity)
	}

	queryUpdateStock := `
		WITH updated AS (
			UPDATE items 
			SET stock = stock - data.qty
			FROM (
				SELECT unnest($1::int[]) as item_id, 
				       unnest($2::int[]) as qty
			) AS data
			WHERE items.id = data.item_id 
			  AND items.stock >= data.qty
			RETURNING items.id
		)
		SELECT COUNT(*) FROM updated
	`

	var updatedCount int
	err = tx.QueryRow(context.Background(), queryUpdateStock,
		itemIds, quantities).Scan(&updatedCount)

	if err != nil {
		r.Logger.Error("failed to batch update stock", zap.Error(err))
		return err
	}

	// Validate all items updated
	if updatedCount != len(items) {
		err = errors.New("insufficient stock for one or more items")
		r.Logger.Error("stock validation failed",
			zap.Int("expected", len(items)),
			zap.Int("updated", updatedCount))
		return err
	}

	// Commit Transaction
	err = tx.Commit(context.Background())
	if err != nil {
		r.Logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	sale.Id = saleId
	r.Logger.Info("sales created successfully",
		zap.Int("sale_id", saleId),
		zap.Int("items_count", len(items)))
	return nil
}

func (r *salesRepository) UpdateSales(id int, data *model.Sales) error {
	query := `
		UPDATE sales
		SET user_id = $1, total_amount = $2
		WHERE id = $3`

	result, err := r.db.Exec(context.Background(), query, data.UserId, data.TotalAmount, id)
	if err != nil {
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("no rows affected")
	}
	return err
}

func (r *salesRepository) DeleteSales(id int) error {
	// Start Transaction
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		r.Logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Delete sale_items first (foreign key constraint)
	queryItems := `DELETE FROM sale_items WHERE sale_id = $1`
	_, err = tx.Exec(context.Background(), queryItems, id)
	if err != nil {
		r.Logger.Error("failed to delete sale items", zap.Error(err))
		return err
	}

	// Delete sales
	querySales := `DELETE FROM sales WHERE id = $1`
	result, err := tx.Exec(context.Background(), querySales, id)
	if err != nil {
		r.Logger.Error("failed to delete sales", zap.Error(err))
		return err
	}

	rowAffected := result.RowsAffected()
	if rowAffected == 0 {
		err = errors.New("no rows affected")
		return err
	}

	// Commit
	err = tx.Commit(context.Background())
	if err != nil {
		r.Logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}