package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciogregory/esp_backIII_go/internal/dentist"
	"github.com/mauriciogregory/esp_backIII_go/internal/domain"
	"github.com/mauriciogregory/esp_backIII_go/pkg/web"
)

type dentistHandler struct {
	s dentist.Service
}

func NewDentistHandler(s dentist.Service) *dentistHandler {
	return &dentistHandler{
		s: s,
	}
}

func (h *dentistHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *dentistHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		response, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "dentist not found")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *dentistHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dentist domain.Dentist
		err := ctx.ShouldBindJSON(&dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid dentist")
			return
		}

		isValid, err := isEmptyDentist(&dentist)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Create(dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

func (h *dentistHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id")
			return
		}
		var dentist domain.Dentist
		err = ctx.ShouldBindJSON(&dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid dentist data")
		}

		isValid, err := isEmptyDentist(&dentist)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Update(id, dentist)
		if err != nil {
			web.BadResponse(ctx, http.StatusConflict, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *dentistHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Surname       string `json:"surname,omitempty"`
		Name          string `json:"name,omitempty"`
		LicenseNumber string `json:"license_number,omitempty"`
	}
	return func(ctx *gin.Context) {
		var r Request
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		if err := ctx.ShouldBindJSON(&r); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid request")
			return
		}
		update := domain.Dentist{
			Surname:       r.Surname,
			Name:          r.Name,
			LicenseNumber: r.LicenseNumber,
		}

		updated, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, updated)
	}
}

func (h *dentistHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.DeleteResponse(ctx, http.StatusOK, "dentist deleted")
	}
}

func isEmptyDentist(dentist *domain.Dentist) (bool, error) {
	switch {
	case dentist.Surname == "" || dentist.Name == "" || dentist.LicenseNumber == "":
		return false, errors.New("fields can't be empty")
	}
	return true, nil
}
