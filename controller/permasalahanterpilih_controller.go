package controller

import "github.com/labstack/echo/v4"

type PermasalahanTerpilihController interface {
	Create(c echo.Context) error
	FindAll(c echo.Context) error
	Delete(c echo.Context) error
}
