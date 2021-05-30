package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateWebRoutes(app *echo.Echo) {
	app.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!")
	})
}
