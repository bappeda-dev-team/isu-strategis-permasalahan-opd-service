package controller

import (
	"net/http"
	"permasalahanService/model/web"
	"permasalahanService/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PermasalahanTerpilihControllerImpl struct {
	permasalahanTerpilihService service.PermasalahanTerpilihService
}

func NewPermasalahanTerpilihControllerImpl(permasalahanTerpilihService service.PermasalahanTerpilihService) *PermasalahanTerpilihControllerImpl {
	return &PermasalahanTerpilihControllerImpl{permasalahanTerpilihService: permasalahanTerpilihService}
}

func (controller *PermasalahanTerpilihControllerImpl) Create(c echo.Context) error {
	request := web.PermasalahanTerpilihRequest{}
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Create Permasalahan Terpilih",
			Data:   err.Error(),
		})
	}

	permasalahanTerpilihResponse, err := controller.permasalahanTerpilihService.Create(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Create Permasalahan Terpilih",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Create Permasalahan Terpilih",
		Data:   permasalahanTerpilihResponse,
	})
}

func (controller *PermasalahanTerpilihControllerImpl) FindAll(c echo.Context) error {
	kodeOpd := c.Request().URL.Query().Get("kode_opd")
	tahun := c.Request().URL.Query().Get("tahun")

	permasalahanTerpilihResponse, err := controller.permasalahanTerpilihService.FindAll(c.Request().Context(), kodeOpd, tahun)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Find All Permasalahan Terpilih",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Find All Permasalahan Terpilih",
		Data:   permasalahanTerpilihResponse,
	})
}

func (controller *PermasalahanTerpilihControllerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Invalid ID",
			Data:   err.Error(),
		})
	}
	err = controller.permasalahanTerpilihService.Delete(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Delete Permasalahan Terpilih",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Delete Permasalahan Terpilih",
		Data:   nil,
	})
}
