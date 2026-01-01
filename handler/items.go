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
)

type ItemsHandler struct {
	ItemsHandlerService service.ItemsService
	config 			utils.Configuration
}

func NewItemsHandler(itemsService service.ItemsService, config utils.Configuration) ItemsHandler {
	return ItemsHandler{
		ItemsHandlerService: itemsService,
		config:				config,
	}
}

func (i *ItemsHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	// Ambil query param page dan limit
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := i.config.Limit
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
			limit = i.config.Limit
		}
	}

	items, total, err := i.ItemsHandlerService.GetAllItems(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": items,
		"pagination": map[string]interface{}{
			"page": page,
			"limit": limit,
			"total": total,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (i *ItemsHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem dto.ItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data", nil)
		return
	}

	
	// validation
	messages, err := utils.ValidateErrors(newItem)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model assignment
	items := model.Items{
		Id: newItem.Id,
		CategoryId: newItem.CategoryId,
		RackId: newItem.RackId,
		Name: newItem.Name,
		Sku: newItem.Sku,
		Stock: newItem.Stock,
		MinStock: newItem.MinStock,
		Price: newItem.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// create assignment service
	err = i.ItemsHandlerService.Create(&items)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success created item", nil)
}