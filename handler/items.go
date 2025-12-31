package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
	"strconv"
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

func (h *ItemsHandler) GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
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

	items, total, err := h.ItemsHandlerService.GetAllItems(page, limit)
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