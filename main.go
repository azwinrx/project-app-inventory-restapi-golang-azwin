package main

import (
	"fmt"
	"log"
	"net/http"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/handler"
	"project-app-inventory-restapi-golang-azwin/repository"
	"project-app-inventory-restapi-golang-azwin/router"
	"project-app-inventory-restapi-golang-azwin/service"
	"project-app-inventory-restapi-golang-azwin/utils"

	"go.uber.org/zap"
)

func main() {
	loadConfig, err := utils.ReadConfigration()
	if err != nil {
		fmt.Printf("Failed to read configuration: %v\n", err)
		return
	}

	db, err := database.InitDB(*loadConfig)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// Initialize repository
	repo := repository.NewRepository(db, logger)
	service := service.NewService(repo)
	handler := handler.NewHandler(service, *loadConfig)


	// Initialize router
	r := router.NewRouter(handler, service, logger)

	// Start server
	addr := fmt.Sprintf(":%d", loadConfig.Port)
	fmt.Printf("Server starting on port %d\n", loadConfig.Port)
	log.Fatal(http.ListenAndServe(addr, r))
}