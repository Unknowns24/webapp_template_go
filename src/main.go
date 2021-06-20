package main

import (
	"app_template/src/libs"
	"app_template/src/routes"
	"app_template/src/utils"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates map[string]*template.Template
}

func NewTemplate() *Template {
	return &Template{
		templates: make(map[string]*template.Template),
	}
}

func (t *Template) Render(w io.Writer, html_name string, data interface{}, c echo.Context) error {
	if tmpl, exist := t.templates[html_name]; exist { //Check existence of the t.templates[html_name]
		return tmpl.ExecuteTemplate(w, "base", data) // ** It wll execute the map[string]interface{} data
	} else {
		return errors.New("There is no " + html_name + " in Template map.")
	}

}

func (tmpl *Template) Add(html_name string, template *template.Template) {
	tmpl.templates[html_name] = template
}

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

	/////////////////////////
	// Serve statics files //
	/////////////////////////

	staticBox, findBoxErr := rice.FindBox("./resources/static")
	if findBoxErr != nil {
		fmt.Println("FATAL ERROR!!!\n Could not find box", err)
		return
	}

	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
	app.GET("/static/*", echo.WrapHandler(staticFileServer))

	///////////////////////
	// Getting templates //
	///////////////////////

	render_htmls := NewTemplate()
	render_htmls.Add("home.html", template.Must(template.ParseFiles("./src/resources/views/home/home.html", "./src/resources/layouts/home/homepage.html")))

	app.Renderer = render_htmls

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
