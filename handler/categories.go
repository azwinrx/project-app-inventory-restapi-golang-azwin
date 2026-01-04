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

type CategoriesHandler struct {
	CategoriesHandlerService service.CategoriesService
	config 			utils.Configuration
}

func NewCategoriesHandler(categoriesService service.CategoriesService, config utils.Configuration) CategoriesHandler {
	return CategoriesHandler{
		CategoriesHandlerService: categoriesService,
		config:				config,
	}
}

func (c *CategoriesHandler) GetCategoriesById(w http.ResponseWriter, r *http.Request) {
	categoriesIDstr := chi.URLParam(r, "id")

	categoriesID, err := strconv.Atoi(categoriesIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	response, err := c.CategoriesHandlerService.GetCategoriesById(categoriesID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "category not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data category by id", response)
}

func (c *CategoriesHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Ambil query param page dan limit
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := c.config.Limit
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
			limit = c.config.Limit
		}
	}

	categories, total, err := c.CategoriesHandlerService.GetAllCategories(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": categories,
		"pagination": map[string]interface{}{
			"page": page,
			"limit": limit,
			"total": total,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *CategoriesHandler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	var newCategories dto.CategoriesRequest
	if err := json.NewDecoder(r.Body).Decode(&newCategories); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data", nil)
		return
	}

	
	// validation
	messages, err := utils.ValidateErrors(newCategories)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model assignment
	categories := model.Categories{
		Name: newCategories.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// create assignment service
	err = c.CategoriesHandlerService.CreateCategories(&categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error creating category",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success created category", categories)
}

func (c *CategoriesHandler) UpdateCategories(w http.ResponseWriter, r *http.Request) {
	categoriesIDstr := chi.URLParam(r, "id")

	categoriesID, err := strconv.Atoi(categoriesIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error param assignment id :"+err.Error(), nil)
		return
	}

	var newCategories dto.CategoriesRequest
	if err := json.NewDecoder(r.Body).Decode(&newCategories); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data :"+err.Error(), nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(newCategories)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing to model assignment
	categories := model.Categories{
		Name: newCategories.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}


	err = c.CategoriesHandlerService.UpdateCategories(categoriesID, &categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error updating category",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success update category", categories)
}

func (c *CategoriesHandler) DeleteCategories(w http.ResponseWriter, r *http.Request) {
	categoriesIDstr := chi.URLParam(r, "id")

	categoriesID, err := strconv.Atoi(categoriesIDstr)
	if err != nil {
		return
	}

	err = c.CategoriesHandlerService.DeleteCategories(categoriesID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "error deleting category",
			"error":   err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success delete category", nil)
}