package repository

import (
	"project-app-inventory-restapi-golang-azwin/database"

	"go.uber.org/zap"
)

type Repository struct {
	ItemsRepo *itemsRepository
	CategoriesRepo *categoriesRepository
	RacksRepo *racksRepository
	WarehousesRepo *warehousesRepository
	UsersRepo *usersRepository
}

func NewRepository(db database.PgxIface, log *zap.Logger) Repository {
	return Repository{
		ItemsRepo: &itemsRepository{db: db, Logger: log},
		CategoriesRepo: &categoriesRepository{db: db, Logger: log},
		RacksRepo: &racksRepository{db: db, Logger: log},
		WarehousesRepo: &warehousesRepository{db: db, Logger: log},
		UsersRepo: &usersRepository{db: db, Logger: log},
	}
}