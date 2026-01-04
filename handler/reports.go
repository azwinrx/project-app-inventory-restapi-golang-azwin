package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
)

type ReportsHandler struct {
	ReportsHandlerService service.ReportsService
	config                utils.Configuration
}

func NewReportsHandler(reportsService service.ReportsService, config utils.Configuration) ReportsHandler {
	return ReportsHandler{
		ReportsHandlerService: reportsService,
		config:                config,
	}
}

func (h *ReportsHandler) GetItemsReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.ReportsHandlerService.GetItemsReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get items report", report)
}

func (h *ReportsHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.ReportsHandlerService.GetSalesReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get sales report", report)
}

func (h *ReportsHandler) GetRevenueReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.ReportsHandlerService.GetRevenueReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get revenue report", report)
}
