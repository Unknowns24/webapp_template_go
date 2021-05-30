package controllers

import (
	models "app_template/src/app/Models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var UserController *usercontrollerController

type usercontrollerController struct {
}

type RequestUsercontroller struct {
	//	Account  string `json:"account" form:"account" query:"account"`
}

func (t *usercontrollerController) Add(c echo.Context) error {
	/*
		u := new(RequestUsercontroller)
		if err := c.Bind(u); err != nil {
			return err
		}
	*/

	models.User.Add("Unknowns", "myPassword", "unknowns0074@gmail.com")

	return c.JSON(http.StatusOK, "Registro Creado")
}

func (t *usercontrollerController) List(c echo.Context) error {
	list := models.User.List()
	return c.JSON(http.StatusOK, list)
}

func (t *usercontrollerController) Show(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	use, _ := models.User.Info(uint(id))
	return c.JSON(http.StatusOK, use)
}

func (t *usercontrollerController) Del(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	models.User.Del(uint(id))
	return c.JSON(http.StatusOK, "Registro eliminado")
}

func (t *usercontrollerController) Update(c echo.Context) error {
	/*
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		u := new(RequestUsercontroller)
		if err := c.Bind(u); err != nil {
			return err
		}

		models.User.Update(uint(id), u)
	*/
	return c.JSON(http.StatusOK, "Registro actualizado")
}
