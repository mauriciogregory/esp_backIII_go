package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mauriciogregory/esp_backIII_go/internal/appointment"
	"github.com/mauriciogregory/esp_backIII_go/internal/domain"
	"github.com/mauriciogregory/esp_backIII_go/pkg/web"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type appointmentHandler struct {
	s appointment.Service
}

func NewAppointmentHandler(s appointment.Service) *appointmentHandler {
	return &appointmentHandler{
		s: s,
	}
}

func (h *appointmentHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := h.s.GetAll()
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		if response == nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", "was not found appointments registered")
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *appointmentHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}
		response, err := h.s.GetByID(id)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *appointmentHandler) GetAllByIdentityNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("identity_number")
		response, err := h.s.GetAllByIdentityNumber(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *appointmentHandler) GetAllByLicenseNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("license_number")
		response, err := h.s.GetAllByLicenseNumber(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}


func (h *appointmentHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var appointment domain.Appointment
		err := ctx.ShouldBindJSON(&appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid appointment data, please verify field(s): "+err.Error())
			return
		}

		isValid, err := isEmptyAppointment(&appointment)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		response, err := h.s.Create(appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

func (h *appointmentHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid id provided")
			return
		}

		var appointment domain.Appointment
		err = ctx.ShouldBindJSON(&appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", "invalid appointment data, verify the fields and try again")
			return
		}

		isValid, err := isEmptyAppointment(&appointment)
		if !isValid {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		response, err := h.s.Update(id, appointment)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}

// metodo patch
func (h *appointmentHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Description     string `json:"description,omitempty"`
		DateAndTime     string `json:"date_and_time,omitempty"`
		DentistLicense  string `json:"dentist_license,omitempty"`
		PatientIdentity string `json:"patient_identity,omitempty"`
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
		update := domain.Appointment{
			Description:     r.Description,
			DateAndTime:     r.DateAndTime,
			DentistLicense:  r.DentistLicense,
			PatientIdentity: r.PatientIdentity,
		}
		if update.DateAndTime != "" {
			if !validateDateTime(update.DateAndTime) {
				web.BadResponse(ctx, http.StatusBadRequest, "error", "please the appointment must be in format: 30/01/2023 23:59")
				return
			}
		}
		response, err := h.s.Update(id, update)
		if err != nil {
			web.BadResponse(ctx, http.StatusNotFound, "error", err.Error())
			return
		}
		web.ResponseOK(ctx, http.StatusOK, response)
	}
}


func (h *appointmentHandler) Delete() gin.HandlerFunc {
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
		web.DeleteResponse(ctx, http.StatusOK, "appointment removed")
	}
}

// Aux functions bellow->

func isEmptyAppointment(appointment *domain.Appointment) (bool, error) {
	dateTimeParsed, err := time.Parse("10/10/2023 20:00", appointment.DateAndTime)
	if err != nil {
		return false, err
	}
	switch {
	case appointment.Description == "" || appointment.DentistLicense == "" || appointment.DateAndTime == "" || appointment.PatientIdentity == "":
		return false, errors.New("fields can't be empty")
	case !validateDateTime(appointment.DateAndTime):
		return false, errors.New("please the appointment must be in format: 30/01/2023 23:59")
	case dateTimeParsed.Local().Add(time.Hour * 3).Before(time.Now().Add(time.Hour)):
		return false, errors.New("the appointment must be in +1 hour from now")
	}
	return true, nil
}

func validateDateTime(dateTime string) bool {
	datesInit := strings.Split(dateTime, " ")
	if len(datesInit) != 2 {
		log.Printf("invalid time, must be in format: 30/01/2023 23:59")
		return false
	}
	breakDate := strings.Split(datesInit[0], "/")
	if len(breakDate) != 3 {
		log.Println("invalid time, must be in format: 30/01/2023 23:59 or 30/01/2023 23:59:59")
		return false
	}
	breakTime := strings.Split(datesInit[1], ":")
	var listDate []int
	var listTime []int

	for _, date := range breakDate {
		number, err := strconv.Atoi(date)
		if err != nil {
			return false
		}
		listDate = append(listDate, number)
	}
	condition := (listDate[0] < 1 || listDate[0] > 31) && (listDate[1] < 1 || listDate[1] > 12) && (listDate[2] < 1 || listDate[2] > 9999)
	if condition {
		log.Println("invalid time, must be between: 1 and 31/12/2023 23:59")
		return false
	}

	for _, t := range breakTime {
		clock, err := strconv.Atoi(t)
		if err != nil {
			log.Println("invalid time, must be in format 23:59 (hours and minutes)")
			return false
		}
		listTime = append(listTime, clock)
	}

	if len(listTime) == 2 {
		condition = (listTime[0] < 0 || listTime[0] > 23) && (listTime[1] < 0 || listTime[1] > 59)
		if condition {
			log.Println("invalid time, must be between: 00:00 and 23:59")
			return false
		}
	}

	if len(listTime) == 3 {
		condition = (listTime[0] < 0 || listTime[0] > 23) && (listTime[1] < 0 || listTime[1] > 59)
		if condition {
			log.Println("invalid time, must be between: 00:00 and 23:59")
			return false
		}
	}
	return true
}
