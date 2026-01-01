package repository

import (
	"project-app-inventory-restapi-golang-azwin/database"

	"go.uber.org/zap"
)

type Repository struct {
	ItemsRepo *itemsRepository
	CategoriesRepo *categoriesRepository
}

func NewRepository(db database.PgxIface, log *zap.Logger) Repository {
	return Repository{
		ItemsRepo: &itemsRepository{db: db, Logger: log},
		CategoriesRepo: &categoriesRepository{db: db, Logger: log},
	}
}