package repository

import (
	"context"
	"project-app-inventory-restapi-golang-azwin/database"

	"go.uber.org/zap"
)

type ItemsReport struct {
	TotalItems    int `json:"total_items"`
	TotalStock    int `json:"total_stock"`
	LowStockItems int `json:"low_stock_items"`
}

type SalesReport struct {
	TotalTransactions int `json:"total_transactions"`
	TotalItemsSold    int `json:"total_items_sold"`
}

type RevenueReport struct {
	TotalRevenue float64 `json:"total_revenue"`
	AveragePerTransaction float64 `json:"average_per_transaction"`
}

type ReportsRepository interface {
	GetItemsReport() (*ItemsReport, error)
	GetSalesReport() (*SalesReport, error)
	GetRevenueReport() (*RevenueReport, error)
}

type reportsRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewReportsRepository(db database.PgxIface, log *zap.Logger) ReportsRepository {
	return &reportsRepository{db: db, Logger: log}
}

func (r *reportsRepository) GetItemsReport() (*ItemsReport, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM items) as total_items,
			(SELECT COALESCE(SUM(stock), 0) FROM items) as total_stock,
			(SELECT COUNT(*) FROM items WHERE stock < min_stock) as low_stock_items
	`

	var report ItemsReport
	err := r.db.QueryRow(context.Background(), query).Scan(
		&report.TotalItems,
		&report.TotalStock,
		&report.LowStockItems,
	)

	if err != nil {
		r.Logger.Error("failed to get items report", zap.Error(err))
		return nil, err
	}

	return &report, nil
}

func (r *reportsRepository) GetSalesReport() (*SalesReport, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM sales) as total_transactions,
			(SELECT COALESCE(SUM(quantity), 0) FROM sale_items) as total_items_sold
	`

	var report SalesReport
	err := r.db.QueryRow(context.Background(), query).Scan(
		&report.TotalTransactions,
		&report.TotalItemsSold,
	)

	if err != nil {
		r.Logger.Error("failed to get sales report", zap.Error(err))
		return nil, err
	}

	return &report, nil
}

func (r *reportsRepository) GetRevenueReport() (*RevenueReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			CASE 
				WHEN COUNT(*) > 0 THEN COALESCE(SUM(total_amount), 0) / COUNT(*)
				ELSE 0
			END as average_per_transaction
		FROM sales
	`

	var report RevenueReport
	err := r.db.QueryRow(context.Background(), query).Scan(
		&report.TotalRevenue,
		&report.AveragePerTransaction,
	)

	if err != nil {
		r.Logger.Error("failed to get revenue report", zap.Error(err))
		return nil, err
	}

	return &report, nil
}
