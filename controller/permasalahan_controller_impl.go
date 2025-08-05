package controller

import (
	"net/http"
	"permasalahanService/model/web"
	"permasalahanService/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PermasalahanControllerImpl struct {
	permasalahanService service.PermasalahanService
}

func NewPermasalahanControllerImpl(permasalahanService service.PermasalahanService) *PermasalahanControllerImpl {
	return &PermasalahanControllerImpl{permasalahanService: permasalahanService}
}

// Create godoc
// @Summary Create Permasalahan
// @Description Create a new permasalahan
// @Tags Permasalahan Service
// @Accept json
// @Produce json
// @Param permasalahan body web.PermasalahanCreateRequest true "Create Permasalahan"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /permasalahan [post]
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

// Update godoc
// @Summary Update Permasalahan
// @Description Update a permasalahan
// @Tags Permasalahan Service
// @Accept json
// @Produce json
// @Param id path string true "Permasalahan ID"
// @Param permasalahan body web.PermasalahanUpdateRequest true "Update Permasalahan"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /permasalahan/{id} [put]
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

	// Set ID dari parameter setelah bind
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Invalid ID",
			Data:   err.Error(),
		})
	}
	request.Id = id

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

// Delete godoc
// @Summary Delete Permasalahan
// @Description Delete an existing permasalahan
// @Tags Permasalahan Service
// @Accept json
// @Produce json
// @Param id path string true "Permasalahan ID"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /permasalahan/{id} [delete]
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

// FindById godoc
// @Summary FindById Permasalahan
// @Description FindById an existing permasalahan
// @Tags Permasalahan Service
// @Param id path int true "Permasalahan ID"
// @Accept json
// @Produce json
// @Success 200 {object} web.WebResponse{data=web.PermasalahanResponsesbyId}
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /permasalahan/{id} [get]
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

// FindAll godoc
// @Summary FindAll Permasalahan
// @Description FindAll an existing permasalahan
// @Tags Permasalahan Service
// @Param kode_opd path string true "Kode OPD"
// @Param tahun path string true "Tahun"
// @Accept json
// @Produce json
// @Success 200 {object} web.WebResponse{data=web.ChildResponse}
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Router /permasalahan/{kode_opd}/{tahun} [get]
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
