package app

import (
	"permasalahanService/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter(permasalahanController controller.PermasalahanController, permasalahanTerpilihController controller.PermasalahanTerpilihController, isuStrategisController controller.IsuStrategisController) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/permasalahan/:kode_opd/:tahun", permasalahanController.FindAllPohonKinerja)
	e.POST("/permasalahan", permasalahanController.Create)
	e.PUT("/permasalahan/:id", permasalahanController.Update)
	e.DELETE("/permasalahan/:id", permasalahanController.Delete)
	e.GET("/permasalahan/:id", permasalahanController.FindById)

	e.POST("/permasalahan_terpilih/create", permasalahanTerpilihController.Create)
	e.GET("/permasalahan_terpilih/findall", permasalahanTerpilihController.FindAll)
	e.DELETE("/permasalahan/:id/hapus_permasalahan_terpilih", permasalahanTerpilihController.Delete)

	e.POST("/isu_strategis", isuStrategisController.Create)
	e.PUT("/isu_strategis/:id", isuStrategisController.Update)
	e.DELETE("/isu_strategis/:id", isuStrategisController.Delete)
	e.GET("/isu_strategis/:id", isuStrategisController.FindById)
	e.GET("/isu_strategis/:kode_opd/:tahun_awal/:tahun_akhir", isuStrategisController.FindAll)

	return e
}
