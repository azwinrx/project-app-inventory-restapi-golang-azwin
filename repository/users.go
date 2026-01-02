package repository

import (
	"context"
	"database/sql"
	"errors"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/model"

	"go.uber.org/zap"
)

type UsersRepository interface {
	FindUsersByEmail(email string) (*model.Users, error)
	CreateUsers(data *model.Users) error
	FindAllUsers() ([]model.Users, error)
	GetUsersByID(id int) (model.Users, error)
	DeleteUsers(id int) error
}

type usersRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewUsersRepository(db database.PgxIface, log *zap.Logger) UsersRepository {
	return &usersRepository{db: db, Logger: log}
}

func (r *usersRepository) FindUsersByEmail(email string) (*model.Users, error) {
	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
		`
	var user model.Users
	err := r.db.QueryRow(context.Background(), query, email).Scan(
			&user.Id, &user.Username, &user.Email, &user.Password, &user.Role,  &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // user tidak ditemukan
	}

	return &user, err
}


func (r *usersRepository) FindAllUsers() ([]model.Users, error) {
	rows, err := r.db.Query(context.Background(), `SELECT id, username, email, password, role, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Users
	for rows.Next() {
		var u model.Users
		err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		students = append(students, u)
	}
	return students, nil
}


func (r *usersRepository) CreateUsers(data *model.Users) error {
	query := `
		INSERT INTO users (username, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`
	return r.db.QueryRow(context.Background(), query, data.Username, data.Email, data.Password, data.Role).Scan(&data.Id)
}



func (r *usersRepository) GetUsersByID(id int) (model.Users, error) {
	var user model.Users
	query := "SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE id = $1"

	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *usersRepository) DeleteUsers(id int) error {
	query := `DELETE FROM users WHERE id = $1`

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