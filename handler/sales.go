package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-azwin/dto"
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SalesHandler struct {
	SalesHandlerService service.SalesService
	config              utils.Configuration
}

func NewSalesHandler(salesService service.SalesService, config utils.Configuration) SalesHandler {
	return SalesHandler{
		SalesHandlerService: salesService,
		config:              config,
	}
}

func (h *SalesHandler) GetSalesById(w http.ResponseWriter, r *http.Request) {
	saleIDstr := chi.URLParam(r, "id")

	saleID, err := strconv.Atoi(saleIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	response, err := h.SalesHandlerService.GetSalesById(saleID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "sale not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data sale by id", response)
}

func (h *SalesHandler) GetAllSales(w http.ResponseWriter, r *http.Request) {
	// Get query param page and limit
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := h.config.Limit
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = h.config.Limit
		}
	}

	sales, total, err := h.SalesHandlerService.GetAllSales(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": sales,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SalesHandler) CreateSales(w http.ResponseWriter, r *http.Request) {
	var newSale dto.SalesRequest
	if err := json.NewDecoder(r.Body).Decode(&newSale); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(newSale)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Create sale
	err = h.SalesHandlerService.CreateSales(&newSale)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error creating sale",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success create sale", nil)
}

func (h *SalesHandler) UpdateSales(w http.ResponseWriter, r *http.Request) {
	saleIDstr := chi.URLParam(r, "id")

	saleID, err := strconv.Atoi(saleIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	var updateSale dto.SalesRequest
	if err := json.NewDecoder(r.Body).Decode(&updateSale); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(updateSale)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Update sale
	err = h.SalesHandlerService.UpdateSales(saleID, &updateSale)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error updating sale",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success update sale", nil)
}

func (h *SalesHandler) DeleteSales(w http.ResponseWriter, r *http.Request) {
	saleIDstr := chi.URLParam(r, "id")

	saleID, err := strconv.Atoi(saleIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	err = h.SalesHandlerService.DeleteSales(saleID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error deleting sale",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success delete sale", nil)
}
