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
	GetUsersByEmail(email string) (*model.Users, error)
	CreateUsers(data *model.Users) error
	GetAllUsers() ([]model.Users, error)
	GetUsersByID(id int) (model.Users, error)
	UpdateUsers(id int, data *model.Users) error
	DeleteUsers(id int) error
}

type usersRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewUsersRepository(db database.PgxIface, log *zap.Logger) UsersRepository {
	return &usersRepository{db: db, Logger: log}
}

func (r *usersRepository) GetUsersByEmail(email string) (*model.Users, error) {
	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
		`
	var user model.Users
	err := r.db.QueryRow(context.Background(), query, email).Scan(
			&user.Id, &user.Username, &user.Email, &user.Password, &user.Role,  &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		r.Logger.Debug("user not found by email", zap.String("email", email))
		return nil, nil // user tidak ditemukan
	}

	if err != nil {
		r.Logger.Error("failed to get user by email",
			zap.String("email", email),
			zap.Error(err),
		)
	}

	return &user, err
}


func (r *usersRepository) GetAllUsers() ([]model.Users, error) {
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
	err := r.db.QueryRow(context.Background(), query, data.Username, data.Email, data.Password, data.Role).Scan(&data.Id)
	if err != nil {
		r.Logger.Error("failed to create user",
			zap.String("username", data.Username),
			zap.String("email", data.Email),
			zap.Error(err),
		)
		return err
	}
	r.Logger.Info("user created successfully",
		zap.Int("user_id", data.Id),
		zap.String("username", data.Username),
		zap.String("email", data.Email),
	)
	return nil
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

func (r *usersRepository) UpdateUsers(id int, data *model.Users) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, role = $4, updated_at = NOW()
		WHERE id = $5
	`
	result, err := r.db.Exec(context.Background(), query, data.Username, data.Email, data.Password, data.Role, id)
	if err != nil {
		r.Logger.Error("failed to update user",
			zap.Int("user_id", id),
			zap.Error(err),
		)
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Warn("user not found for update", zap.Int("user_id", id))
		return errors.New("user not found or already deleted")
	}
	r.Logger.Info("user updated successfully",
		zap.Int("user_id", id),
		zap.String("username", data.Username),
	)
	return nil
}

func (r *usersRepository) DeleteUsers(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("failed to delete user",
			zap.Int("user_id", id),
			zap.Error(err),
		)
		return err
	}
	rowAffected := result.RowsAffected()

	if rowAffected == 0 {
		r.Logger.Warn("user not found for deletion", zap.Int("user_id", id))
		return errors.New("no rows affected")
	}

	r.Logger.Info("user deleted successfully", zap.Int("user_id", id))
	return nil
}