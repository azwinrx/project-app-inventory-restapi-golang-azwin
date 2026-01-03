package repository

import (
	"context"
	"project-app-inventory-restapi-golang-azwin/database"

	"go.uber.org/zap"
)

type ReportData struct {
	TotalItems    int     `json:"total_items"`
	TotalStock    int     `json:"total_stock"`
	TotalSales    int     `json:"total_sales"`
	TotalRevenue  float64 `json:"total_revenue"`
	LowStockItems int     `json:"low_stock_items"`
}

type ItemSalesReport struct {
	ItemId       int     `json:"item_id"`
	ItemName     string  `json:"item_name"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type ReportsRepository interface {
	GetDashboardReport() (*ReportData, error)
	GetItemSalesReport() ([]ItemSalesReport, error)
}

type reportsRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewReportsRepository(db database.PgxIface, log *zap.Logger) ReportsRepository {
	return &reportsRepository{db: db, Logger: log}
}

func (r *reportsRepository) GetDashboardReport() (*ReportData, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM items) as total_items,
			(SELECT COALESCE(SUM(stock), 0) FROM items) as total_stock,
			(SELECT COUNT(*) FROM sales) as total_sales,
			(SELECT COALESCE(SUM(total_amount), 0) FROM sales) as total_revenue,
			(SELECT COUNT(*) FROM items WHERE stock < min_stock) as low_stock_items
	`

	var report ReportData
	err := r.db.QueryRow(context.Background(), query).Scan(
		&report.TotalItems,
		&report.TotalStock,
		&report.TotalSales,
		&report.TotalRevenue,
		&report.LowStockItems,
	)

	if err != nil {
		r.Logger.Error("failed to get dashboard report", zap.Error(err))
		return nil, err
	}

	return &report, nil
}

func (r *reportsRepository) GetItemSalesReport() ([]ItemSalesReport, error) {
	query := `
		SELECT 
			i.id as item_id,
			i.name as item_name,
			COALESCE(SUM(si.quantity), 0) as total_sold,
			COALESCE(SUM(si.subtotal), 0) as total_revenue
		FROM items i
		LEFT JOIN sale_items si ON i.id = si.item_id
		GROUP BY i.id, i.name
		ORDER BY total_revenue DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		r.Logger.Error("failed to get item sales report", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var reports []ItemSalesReport
	for rows.Next() {
		var report ItemSalesReport
		err := rows.Scan(
			&report.ItemId,
			&report.ItemName,
			&report.TotalSold,
			&report.TotalRevenue,
		)
		if err != nil {
			r.Logger.Error("failed to scan item sales report", zap.Error(err))
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}
