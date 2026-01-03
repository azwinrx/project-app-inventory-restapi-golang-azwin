package database

import (
	"context"
	"fmt"
	"project-app-inventory-restapi-golang-azwin/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxIface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

func InitDB(config utils.Configuration) (PgxIface, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Host, config.DB.Port)

	conn, err := pgx.Connect(context.Background(), connStr)
	
	//test connection
	ping := conn.Ping(context.Background())
	if ping != nil {
		fmt.Printf("Gagal terhubung ke database: %s\n", ping)
		return nil, ping
	}
	fmt.Println("Berhasil terhubung ke database")
	return conn, err
	

}