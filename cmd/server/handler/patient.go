package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciogregory/esp_backIII_go/internal/domain"
	"github.com/mauriciogregory/esp_backIII_go/internal/patient"
	"github.com/mauriciogregory/esp_backIII_go/pkg/web"
)

type patientHandler struct {
	s patient.Service
}

func NewPatientHandler(s patient.Service) *patientHandler {
	return &patientHandler{
		s: s,
	}
}

func (h *patientHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		patients, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, patients)
	}
}

func (h *patientHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		patient, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "patient not found")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, patient)
	}
}

func (h *patientHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var patient domain.Patient
		err := ctx.ShouldBindJSON(&patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient")
			return
		}

		isValid, err := isEmptyPatient(&patient)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Create(patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		web.ResponseOK(ctx, http.StatusCreated, response)
	}
}

func (h *patientHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient id provided")
			return
		}
		var patient domain.Patient
		err = ctx.ShouldBindJSON(&patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid patient data")
		}

		isValid, err := isEmptyPatient(&patient)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		response, err := h.s.Update(id, patient)
		if err != nil {
			web.BadResponse(ctx, http.StatusConflict, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}


func (h *patientHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Surname        string `json:"surname,omitempty"`
		Name           string `json:"name,omitempty"`
		IdentityNumber string `json:"identity_number,omitempty"`
		CreatedAt      string `json:"created_at,omitempty"`
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
		update := domain.Patient{
			Surname:        r.Surname,
			Name:           r.Name,
			IdentityNumber: r.IdentityNumber,
			CreatedAt:      r.CreatedAt,
		}
		if update.CreatedAt != "" {
			if !validateDateTime(update.CreatedAt) {
				web.BadResponse(ctx, http.StatusBadRequest, "error", "please the patient created_at field must be in format: 30/01/2023 23:59 or 30/01/2023 23:59:59")
				return
			}
		}
		response, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *patientHandler) Delete() gin.HandlerFunc {
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
		web.DeleteResponse(ctx, http.StatusOK, "patient deleted")
	}
}

func isEmptyPatient(patient *domain.Patient) (bool, error) {
	switch {
	case patient.Surname == "" || patient.Name == "" || patient.CreatedAt == "" || patient.IdentityNumber == "":
		return false, errors.New("patient fields can't be empty")
	case !validateDateTime(patient.CreatedAt):
		return false, errors.New("please the patient created_at field must be in format: 30/01/2023 23:59 or 30/01/2023 23:59:59")
	}
	return true, nil
}
