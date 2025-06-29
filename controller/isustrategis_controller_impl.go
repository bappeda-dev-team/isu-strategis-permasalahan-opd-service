package controller

import (
	"net/http"
	"permasalahanService/model/web"
	"permasalahanService/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IsuStrategisControllerImpl struct {
	IsuStrategisService service.IsuStrategisService
}

func NewIsuStrategisControllerImpl(isuStrategisService service.IsuStrategisService) *IsuStrategisControllerImpl {
	return &IsuStrategisControllerImpl{IsuStrategisService: isuStrategisService}
}

func (controller *IsuStrategisControllerImpl) Create(c echo.Context) error {
	request := web.IsuStrategisCreateRequest{}
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Create Isu Strategis",
			Data:   err.Error(),
		})
	}

	isuStrategisResponse, err := controller.IsuStrategisService.Create(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Create Isu Strategis",
			Data:   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Create Isu Strategis",
		Data:   isuStrategisResponse,
	})
}

func (controller *IsuStrategisControllerImpl) Update(c echo.Context) error {
	request := web.IsuStrategisUpdateRequest{}
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Update Isu Strategis",
			Data:   err.Error(),
		})
	}
	request.Id = idInt

	err = c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Update Isu Strategis",
			Data:   err.Error(),
		})
	}
	isuStrategisResponse, err := controller.IsuStrategisService.Update(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Update Isu Strategis",
			Data:   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Update Isu Strategis",
		Data:   isuStrategisResponse,
	})
}

func (controller *IsuStrategisControllerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Delete Isu Strategis",
			Data:   err.Error(),
		})
	}
	err = controller.IsuStrategisService.Delete(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Delete Isu Strategis",
			Data:   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Delete Isu Strategis",
	})
}

func (controller *IsuStrategisControllerImpl) FindById(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Find Isu Strategis",
			Data:   err.Error(),
		})
	}
	isuStrategisResponse, err := controller.IsuStrategisService.FindById(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Find Isu Strategis",
			Data:   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Find Isu Strategis",
		Data:   isuStrategisResponse,
	})
}

func (controller *IsuStrategisControllerImpl) FindAll(c echo.Context) error {
	kodeOpd := c.Param("kode_opd")
	tahunAwal := c.Param("tahun_awal")
	tahunAkhir := c.Param("tahun_akhir")
	isuStrategisResponse, err := controller.IsuStrategisService.FindAll(c.Request().Context(), kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Find Isu Strategis",
			Data:   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Find Isu Strategis",
		Data:   isuStrategisResponse,
	})
}
