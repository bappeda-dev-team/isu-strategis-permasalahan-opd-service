package controller

import (
	"net/http"
	"permasalahanService/model/web"
	"permasalahanService/service"

	"github.com/labstack/echo/v4"
)

type PermasalahanControllerImpl struct {
	permasalahanService service.PermasalahanService
}

func NewPermasalahanControllerImpl(permasalahanService service.PermasalahanService) *PermasalahanControllerImpl {
	return &PermasalahanControllerImpl{permasalahanService: permasalahanService}
}

func (controller *PermasalahanControllerImpl) Create(c echo.Context) error {
	request := web.PermasalahanCreateRequest{}
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Create Permasalahan",
			Data:   err.Error(),
		})
	}

	permasalahanResponse, err := controller.permasalahanService.Create(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Create Permasalahan",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Create Permasalahan",
		Data:   permasalahanResponse,
	})
}

func (controller *PermasalahanControllerImpl) Update(c echo.Context) error {

	request := web.PermasalahanUpdateRequest{}
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Failed Update Permasalahan",
			Data:   err.Error(),
		})
	}

	permasalahanResponse, err := controller.permasalahanService.Update(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Update Permasalahan",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Update Permasalahan",
		Data:   permasalahanResponse,
	})
}

func (controller *PermasalahanControllerImpl) Delete(c echo.Context) error {
	id := c.Param("id")

	err := controller.permasalahanService.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Delete Permasalahan",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Delete Permasalahan",
		Data:   nil,
	})
}

func (controller *PermasalahanControllerImpl) FindById(c echo.Context) error {
	id := c.Param("id")

	permasalahanResponse, err := controller.permasalahanService.FindById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Find Permasalahan",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Find Permasalahan",
		Data:   permasalahanResponse,
	})
}

// FindAllPohonKinerja godoc
// @Summary Get All Pohon Kinerja with Permasalahan
// @Description Get all pohon kinerja and merge with permasalahan data
// @Tags GET Pohon Kinerja
// @Accept json
// @Produce json
// @Param kode_opd path string true "Kode OPD"
// @Param tahun path string true "Tahun"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /pohon_kinerja_opd/findall/{kode_opd}/{tahun} [get]
func (controller *PermasalahanControllerImpl) FindAllPohonKinerja(c echo.Context) error {
	kodeOpd := c.Param("kode_opd")
	tahun := c.Param("tahun")

	result, err := controller.permasalahanService.FindAllPohonKinerja(c.Request().Context(), kodeOpd, tahun)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed Get All Pohon Kinerja",
			Data:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: "Success Get All Pohon Kinerja",
		Data:   result,
	})
}
