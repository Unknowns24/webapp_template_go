package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var HomeController *homeController

type homeController struct {
}

func (t *homeController) ShowHome(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"name":      "Mordern Artist",
		"greetings": "I'm very pleased to see you",
	})
}
