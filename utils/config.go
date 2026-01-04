package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	AppName     string
	Port        int
	Debug       bool
	Limit       int
	PathLogging string
	DB          DatabaseCofig
}

type DatabaseCofig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     int
	MaxConn  int32
}

func ReadConfigration() (*Configuration, error) {
	//Load env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")


	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error membaca file konfigurasi: %s\n", err)
		return nil, err
	}
	
	//Initialize configuration variables
	appName := viper.GetString("APP_NAME")
	port := viper.GetInt("PORT")
	debug := viper.GetBool("DEBUG")
	limit := viper.GetInt("LIMIT")
	pathLogging := viper.GetString("PATH_LOGGING")

	// Default values
	if limit == 0 {
		limit = 10
	}
	if pathLogging == "" {
		pathLogging = "./logs/"
	}

	dbUser := viper.GetString("DATABASE_USERNAME")
	dbPassword := viper.GetString("DATABASE_PASSWORD")
	dbHost := viper.GetString("DATABASE_HOST")
	dbPort := viper.GetInt("DATABASE_PORT")
	dbName := viper.GetString("DATABASE_NAME")
	maxConn := viper.GetInt32("DATABASE_MAX_CONN")

	return &Configuration{
		AppName: appName,
		Port:    port,
		Debug:   debug,
		Limit:   limit,
		PathLogging: pathLogging,
		DB: DatabaseCofig{
			Name:     dbName,
			Username: dbUser,
			Password: dbPassword,
			Host:     dbHost,
			Port:     dbPort,
			MaxConn:  maxConn,
		},
	}, nil
}