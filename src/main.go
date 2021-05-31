package main

import (
	"app_template/src/libs"
	"app_template/src/routes"
	"app_template/src/utils"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	/////////////////////////////
	// Loading the config data //
	/////////////////////////////

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

	// Opening the database connection
	libs.DB = dbCfg.InitMysqlDB()

	// Initializing web app
	app := echo.New()

	//////////////////////////////////
	// Using importants middlewares //
	//////////////////////////////////

	// Recover middleware
	app.Use(middleware.Recover())

	// Logging middleware
	myLogFile, err := os.OpenFile("./src/logs/requests.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("FATAL ERROR!!!\n Cannot open or create log file:", err)
		return
	}
	defer myLogFile.Close()

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: myLogFile,
	}))

	///////////////////
	// Adding routes //
	///////////////////

	routes.CreateAuthRoutes(app) // Create login, logoute, register, verificate Routes
	routes.CreateWebRoutes(app)  // Create web routes f.e. /, contact
	routes.CreateApiRoutes(app)  // create api routes

	//////////////////
	// Starting app //
	//////////////////

	app.Logger.Fatal(app.Start(config.APPPort))
}
