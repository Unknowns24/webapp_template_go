package routes

import (
	controllers "app_template/src/app/Http/controllers/user"

	"github.com/labstack/echo/v4"
)

func CreateWebRoutes(app *echo.Echo) {
	app.GET("/", controllers.UserController.Add)
}
