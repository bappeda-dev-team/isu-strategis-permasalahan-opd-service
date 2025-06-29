package controller

import "github.com/labstack/echo/v4"

type IsuStrategisController interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	FindById(c echo.Context) error
	FindAll(c echo.Context) error
}
