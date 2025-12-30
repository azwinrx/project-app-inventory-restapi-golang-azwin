package main

import (
	"fmt"
	"project-app-inventory-restapi-golang-azwin/database"
	"project-app-inventory-restapi-golang-azwin/utils"
)

func main() {
	loadConfig, err := utils.ReadConfigration()
	if err != nil {
		fmt.Printf("Failed to read configuration: %v\n", err)
		return
	}
	
	db, err := database.InitDB(*loadConfig)
	_ = db 
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
}