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

// Create godoc
// @Summary Create Isu Strategis
// @Description Create a new isu strategis
// @Tags Isu Strategis Service
// @Accept json
// @Produce json
// @Param isu_strategis body web.IsuStrategisCreateRequest true "Create Isu Strategis"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /isu_strategis [post]
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

// Update godoc
// @Summary Update Isu Strategis
// @Description Update a isu strategis
// @Tags Isu Strategis Service
// @Accept json
// @Produce json
// @Param id path string true "Isu Strategis ID"
// @Param isu_strategis body web.IsuStrategisUpdateRequest true "Update Isu Strategis"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /isu_strategis/{id} [put]
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

// Delete godoc
// @Summary Delete Isu Strategis
// @Description Delete a isu strategis
// @Tags Isu Strategis Service
// @Accept json
// @Produce json
// @Param id path string true "Isu Strategis ID"
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

// FindById godoc
// @Summary FindById Isu Strategis
// @Description FindById a isu strategis
// @Tags Isu Strategis Service
// @Accept json
// @Produce json
// @Param id path string true "Isu Strategis ID"
// @Success 200 {object} web.WebResponse{data=web.IsuStrategisResponse}
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /isu_strategis/{id} [get]
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

// FindAll godoc
// @Summary FindAll Isu Strategis
// @Description FindAll a isu strategis
// @Tags Isu Strategis Service
// @Accept json
// @Produce json
// @Param kode_opd path string true "Kode OPD"
// @Param tahun_awal path string true "Tahun Awal"
// @Param tahun_akhir path string true "Tahun Akhir"
// @Success 200 {object} web.WebResponse{data=web.IsuStrategisResponse}
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /isu_strategis/{kode_opd}/{tahun_awal}/{tahun_akhir} [get]
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
