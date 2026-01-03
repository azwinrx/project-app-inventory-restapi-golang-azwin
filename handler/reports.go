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

func (h *ReportsHandler) GetDashboardReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.ReportsHandlerService.GetDashboardReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get dashboard report", report)
}

func (h *ReportsHandler) GetItemSalesReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.ReportsHandlerService.GetItemSalesReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get item sales report", report)
}
