package routes

import (
	"app_template/src/app/Http/controllers"

	"github.com/labstack/echo/v4"
)

func CreateWebRoutes(app *echo.Echo) {
	app.GET("/", controllers.HomeController.ShowHome)
}
