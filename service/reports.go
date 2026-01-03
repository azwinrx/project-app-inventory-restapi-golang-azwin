package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type ReportsService interface {
	GetDashboardReport() (*repository.ReportData, error)
	GetItemSalesReport() ([]repository.ItemSalesReport, error)
}

type reportsService struct {
	Repo repository.ReportsRepository
}

func NewReportsService(repo repository.ReportsRepository) ReportsService {
	return &reportsService{Repo: repo}
}

func (s *reportsService) GetDashboardReport() (*repository.ReportData, error) {
	return s.Repo.GetDashboardReport()
}

func (s *reportsService) GetItemSalesReport() ([]repository.ItemSalesReport, error) {
	return s.Repo.GetItemSalesReport()
}
