package main

import (
	"app_template/src/libs"
	"app_template/src/routes"
	"app_template/src/utils"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		fmt.Println("FATAL ERROR!!!\nCannot load config file:", err)
		return
	}

	// Setting database connection data
	dbCfg := libs.DbConfig{
		Host:         config.DBHost,
		Port:         config.DBPort,
		Database:     config.DBName,
		User:         config.DBUser,
		Password:     config.DBPass,
		Charset:      config.DBChar,
		MaxIdleConns: int(config.DBMaxIdleConns),
		MaxOpenConns: int(config.DBMaxOpenConns),
		TimeZone:     "",
		Print_log:    config.DBLog,
	}

	libs.DB = dbCfg.InitMysqlDB()

	app := echo.New()

	routes.CreateAuthRoutes(app) // Create login, logoute, register, verificate Routes
	routes.CreateWebRoutes(app)  // Create web routes f.e. /, contact
	routes.CreateApiRoutes(app)  // create api routes

	app.Logger.Fatal(app.Start(config.APPPort))
}
