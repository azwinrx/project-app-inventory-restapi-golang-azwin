package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-azwin/dto"
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type WarehousesHandler struct {
	WarehousesHandlerService service.WarehousesService
	config                   utils.Configuration
}

func NewWarehousesHandler(warehousesService service.WarehousesService, config utils.Configuration) WarehousesHandler {
	return WarehousesHandler{
		WarehousesHandlerService: warehousesService,
		config:                   config,
	}
}

func (h *WarehousesHandler) GetWarehousesById(w http.ResponseWriter, r *http.Request) {
	warehousesIDstr := chi.URLParam(r, "id")

	warehousesID, err := strconv.Atoi(warehousesIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	response, err := h.WarehousesHandlerService.GetWarehousesById(warehousesID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "warehouse not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data warehouse by id", response)
}

func (h *WarehousesHandler) GetAllWarehouses(w http.ResponseWriter, r *http.Request) {
	// Ambil query param page dan limit
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

	warehouses, total, err := h.WarehousesHandlerService.GetAllWarehouses(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": warehouses,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *WarehousesHandler) CreateWarehouses(w http.ResponseWriter, r *http.Request) {
	var newWarehouses dto.WarehousesRequest
	if err := json.NewDecoder(r.Body).Decode(&newWarehouses); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data", nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(newWarehouses)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model
	warehouses := model.Warehouses{
		Name:      newWarehouses.Name,
		Location:  newWarehouses.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// create service
	err = h.WarehousesHandlerService.CreateWarehouses(&warehouses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error creating warehouse",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success created warehouse", warehouses)
}

func (h *WarehousesHandler) UpdateWarehouses(w http.ResponseWriter, r *http.Request) {
	warehousesIDstr := chi.URLParam(r, "id")

	warehousesID, err := strconv.Atoi(warehousesIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error param warehouse id :"+err.Error(), nil)
		return
	}

	var newWarehouses dto.WarehousesRequest
	if err := json.NewDecoder(r.Body).Decode(&newWarehouses); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data :"+err.Error(), nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(newWarehouses)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model
	warehouses := model.Warehouses{
		Name:      newWarehouses.Name,
		Location:  newWarehouses.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.WarehousesHandlerService.UpdateWarehouses(warehousesID, &warehouses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error updating warehouse",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success update warehouse", warehouses)
}

func (h *WarehousesHandler) DeleteWarehouses(w http.ResponseWriter, r *http.Request) {
	warehousesIDstr := chi.URLParam(r, "id")

	warehousesID, err := strconv.Atoi(warehousesIDstr)
	if err != nil {
		return
	}

	err = h.WarehousesHandlerService.DeleteWarehouses(warehousesID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error deleting warehouse",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success delete warehouse", nil)
}
