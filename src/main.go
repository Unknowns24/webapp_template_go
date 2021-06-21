package main

import (
	"app_template/src/libs"
	"app_template/src/routes"
	"app_template/src/utils"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
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

	/////////////////////
	// Template loader //
	/////////////////////

	files, err := ioutil.ReadDir("./src/resources/layouts/") // Get folders in the layout dir

	// Check if there is an error
	if err != nil {
		fmt.Println("FATAL ERROR!!!\n Could not open layout folder", err)
		return
	}

	render_htmls := NewTemplate() // Create a new template list

	// List files and register the templates
	for _, rootFile := range files {
		if rootFile.IsDir() {
			subDirFiles, subDirErr := ioutil.ReadDir("./src/resources/layouts/" + rootFile.Name()) // List subdir files

			// Check if an error ocurred
			if subDirErr != nil {
				fmt.Println("FATAL ERROR!!!\n Could not open "+rootFile.Name()+" layout folder", err)
				return
			}

			var fileList []string

			// List files
			for _, file := range subDirFiles {
				if !file.IsDir() {
					fileList = append(fileList, "./src/resources/layouts/"+rootFile.Name()+"/"+file.Name())
				}
			}

			viewsFiles, viewsErr := ioutil.ReadDir("./src/resources/views/" + rootFile.Name()) // List subdir files

			// Check if an error ocurred
			if viewsErr != nil {
				fmt.Println("FATAL ERROR!!!\n Could not open "+rootFile.Name()+" views folder", err)
				return
			}

			// List views files of the layout
			for _, vfile := range viewsFiles {
				if !vfile.IsDir() {
					var params []string
					params = append(params, fileList...)
					params = append(params, "./src/resources/views/"+rootFile.Name()+"/"+vfile.Name())

					render_htmls.Add(vfile.Name(), template.Must(template.ParseFiles(params...)))
				}
			}
		}
	}

	app.Renderer = render_htmls // Set app renderer

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
