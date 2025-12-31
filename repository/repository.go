package repository

import (
	"project-app-inventory-restapi-golang-azwin/database"
)

type Repository struct {
	ItemsRepo itemsRepository
}

func NewRepository(db database.PgxIface) Repository {
	return Repository{
		ItemsRepo: itemsRepository{db: db},
	}
}