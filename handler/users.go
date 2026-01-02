package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-azwin/dto"
	"project-app-inventory-restapi-golang-azwin/model"
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	UsersHandlerService service.UsersService
	config              utils.Configuration
}

func NewUsersHandler(usersService service.UsersService, config utils.Configuration) UsersHandler {
	return UsersHandler{
		UsersHandlerService: usersService,
		config:              config,
	}
}

// GetUsersByID - Get user by ID
func (u *UsersHandler) GetUsersByID(w http.ResponseWriter, r *http.Request) {
	usersIDStr := chi.URLParam(r, "id")

	usersID, err := strconv.Atoi(usersIDStr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid id format", nil)
		return
	}

	response, err := u.UsersHandlerService.GetUsersByID(usersID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "user not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data user by id", response)
}

// GetAllUsers - Get all users with pagination
func (u *UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.UsersHandlerService.FindAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get all users", users)
}

// GetUsersByEmail - Get user by email
func (u *UsersHandler) GetUsersByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "email parameter is required", nil)
		return
	}

	user, err := u.UsersHandlerService.FindUsersByEmail(email)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "error finding user", nil)
		return
	}

	if user == nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, "user not found", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get user by email", user)
}

// CreateUsers - Create new user
func (u *UsersHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	var userReq dto.Usersrequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// Hash password before saving
	hashedPassword := utils.HashPassword(userReq.Password)

	// Map DTO to model
	users := model.Users{
		Username: userReq.Username,
		Email:    userReq.Email,
		Password: hashedPassword,
		Role:     userReq.Role,
	}

	err = u.UsersHandlerService.CreateUsers(&users)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "error creating user", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success create user", users)
}


