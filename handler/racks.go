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

type RacksHandler struct {
	RacksHandlerService service.RacksService
	config              utils.Configuration
}

func NewRacksHandler(racksService service.RacksService, config utils.Configuration) RacksHandler {
	return RacksHandler{
		RacksHandlerService: racksService,
		config:              config,
	}
}

func (h *RacksHandler) GetRacksById(w http.ResponseWriter, r *http.Request) {
	racksIDstr := chi.URLParam(r, "id")

	racksID, err := strconv.Atoi(racksIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	response, err := h.RacksHandlerService.GetRacksById(racksID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "rack not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data rack by id", response)
}

func (h *RacksHandler) GetAllRacks(w http.ResponseWriter, r *http.Request) {
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

	racks, total, err := h.RacksHandlerService.GetAllRacks(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": racks,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RacksHandler) CreateRacks(w http.ResponseWriter, r *http.Request) {
	var newRacks dto.RacksRequest
	if err := json.NewDecoder(r.Body).Decode(&newRacks); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data", nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(newRacks)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model
	racks := model.Racks{
		WarehouseId: newRacks.WarehouseId,
		Name:        newRacks.Name,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// create service
	err = h.RacksHandlerService.CreateRacks(&racks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error creating rack",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success created rack", racks)
}

func (h *RacksHandler) UpdateRacks(w http.ResponseWriter, r *http.Request) {
	racksIDstr := chi.URLParam(r, "id")

	racksID, err := strconv.Atoi(racksIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error param rack id :"+err.Error(), nil)
		return
	}

	var newRacks dto.RacksRequest
	if err := json.NewDecoder(r.Body).Decode(&newRacks); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data :"+err.Error(), nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(newRacks)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model
	racks := model.Racks{
		WarehouseId: newRacks.WarehouseId,
		Name:        newRacks.Name,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.RacksHandlerService.UpdateRacks(racksID, &racks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error updating rack",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success update rack", racks)
}

func (h *RacksHandler) DeleteRacks(w http.ResponseWriter, r *http.Request) {
	racksIDstr := chi.URLParam(r, "id")

	racksID, err := strconv.Atoi(racksIDstr)
	if err != nil {
		return
	}

	err = h.RacksHandlerService.DeleteRacks(racksID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error deleting rack",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success delete rack", nil)
}
