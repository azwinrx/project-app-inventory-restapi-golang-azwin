package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"

	"go.uber.org/zap"
)

type CategoriesRepository interface {
	GetCategoriesById(id int) (*model.Categories, error)
	GetAllCategories(page, limit int) ([]model.Categories, int, error)
	CreateCategories(data *model.Categories) error
	UpdateCategories(id int, data *model.Categories) error
	DeleteCategories(id int) error
}

type categoriesRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewCategoriesRepository(db database.PgxIface, log *zap.Logger) CategoriesRepository {
	return &categoriesRepository{db: db, Logger: log}
}

func (r *categoriesRepository) GetCategoriesById(id int) (*model.Categories, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	var c model.Categories
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&c.Id,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	return &c, err
}


func (r *categoriesRepository) GetAllCategories(page, limit int) ([]model.Categories, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM categories`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findall repo ", zap.Error(err))
		return nil, 0, err
	}

	// get data with pagination
	query := `
		SELECT id, name, created_at, updated_at
		FROM categories
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []model.Categories
	for rows.Next() {
		var c model.Categories
		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	return categories, total, nil
}

func (r *categoriesRepository) CreateCategories(data *model.Categories) error {
	query := `
		INSERT INTO categories (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, data.Name).Scan(&data.Id)
	return err
}

func (r *categoriesRepository) UpdateCategories(id int, data *model.Categories) error {
	query := `
		UPDATE categories
		SET name = $1, updated_at = NOW()
		WHERE id = $2`

	result, err := r.db.Exec(context.Background(), query, data.Name, id)
	if err != nil {
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("no rows affected")
	}
	return err
}

func (r *categoriesRepository) DeleteCategories(id int) error {
	query := `DELETE FROM categories WHERE id = $1`

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