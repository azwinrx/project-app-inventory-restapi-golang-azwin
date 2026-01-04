package service

import (
	"project-app-inventory-restapi-golang-azwin/repository"
)

type ReportsService interface {
	GetItemsReport() (*repository.ItemsReport, error)
	GetSalesReport() (*repository.SalesReport, error)
	GetRevenueReport() (*repository.RevenueReport, error)
}

type reportsService struct {
	Repo repository.ReportsRepository
}

func NewReportsService(repo repository.ReportsRepository) ReportsService {
	return &reportsService{Repo: repo}
}

func (s *reportsService) GetItemsReport() (*repository.ItemsReport, error) {
	return s.Repo.GetItemsReport()
}

func (s *reportsService) GetSalesReport() (*repository.SalesReport, error) {
	return s.Repo.GetSalesReport()
}

func (s *reportsService) GetRevenueReport() (*repository.RevenueReport, error) {
	return s.Repo.GetRevenueReport()
}
