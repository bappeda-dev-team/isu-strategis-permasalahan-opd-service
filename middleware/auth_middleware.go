package middleware

import (
	"net/http"
	"permasalahanService/model/web"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-API-Key") == "RAHASIA" {
			return next(c)
		}

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		return c.JSON(http.StatusUnauthorized, webResponse)
	}
}
