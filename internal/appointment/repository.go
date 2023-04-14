package appointment

import (
	"errors"
	"github.com/mauriciogregory/esp_backIII_go/internal/domain"
	"github.com/mauriciogregory/esp_backIII_go/pkg/store"
	"log"
	"time"
)

var table = store.AP

type Repository interface {
	GetAll() (interface{}, error)
	GetByID(entityId int) (interface{}, error)
	GetAllByIdentityNumber(identityNumber string) (interface{}, error)
	GetAllByLicenseNumber(licenseNumber string) (interface{}, error)
	Create(a domain.Appointment) (interface{}, error)
	Update(entityId int, a domain.Appointment) (interface{}, error)
	Delete(entityId int) error
}

type repository struct {
	store store.ApStore
}

func NewRepository(store store.ApStore) Repository {
	return &repository{store}
}

func (r *repository) GetAll() (interface{}, error) {
	return r.store.GetAll(table)
}

func (r *repository) GetByID(entityId int) (interface{}, error) {
	return r.store.GetByID(entityId, table)
}

func (r *repository) GetAllByIdentityNumber(identityNumber string) (interface{}, error) {
	return r.store.GetAllAppointmentsByPatientIdentify(identityNumber)
}

func (r *repository) GetAllByLicenseNumber(licenseNumber string) (interface{}, error) {
	return r.store.GetAllAppointmentsByDentistsLicense(licenseNumber)
}

func (r *repository) Create(a domain.Appointment) (interface{}, error) {
	if !r.isValidDate(a) {
		return nil, errors.New("some data is invalid")
	}
	if !r.isADateTimeAvailable(a.DateAndTime, a.DentistLicense, a.PatientIdentity) {
		return nil, errors.New("the date and time select aren't available for dentist or patient")
	}
	return r.store.Save(a, table)
}

func (r *repository) Update(entityId int, a domain.Appointment) (interface{}, error) {
	aInterface, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	appointments, ok := aInterface.([]domain.AppointmentDTO)
	if !ok {
		return nil, err
	}

	for _, appointment := range appointments {
		if appointment.Id == entityId {
			if !r.isValidDate(a) {
				return nil, errors.New("some data is invalid")
			}
			if !r.isADateTimeAvailable(a.DateAndTime, a.DentistLicense, a.PatientIdentity) {
				return nil, errors.New("the date and time select aren't available for dentist or patient")
			}
			return r.store.Update(entityId, a, table)
		}
	}
	return nil, errors.New("appointment not found")
}

func (r *repository) Delete(entityId int) error {
	return r.store.Delete(entityId, table)
}

func (r *repository) isValidDate(a domain.Appointment) bool {
	var appointments []domain.AppointmentDTO
	aDateTimeToValidate, err := time.Parse("10/10/2023 20:00", a.DateAndTime)
	if err != nil {
		log.Fatalln("error while trying to validate date and time provided from request body ->", err.Error())
		return false
	}

	appointmentsInterface, err := r.GetAll()
	if err != nil {
		log.Fatalln("error: ", err.Error())
		return false
	}
	appointments, ok := appointmentsInterface.([]domain.AppointmentDTO)
	if !ok {
		log.Fatalln("error parsing interface data fetched from db")
		return false
	}

	for _, appointment := range appointments {
		if err != nil {
			log.Println("error while trying to parse date and time when validating date_and_time provided")
			return false
		}
		if a.Id == appointment.Id {
			return aDateTimeToValidate.After(time.Now().Add(time.Hour))
		}
	}
	return aDateTimeToValidate.After(time.Now().Add(time.Hour))
}

func (r *repository) isADateTimeAvailable(dateTime string, d string, p string) bool {
	dateTimeToVerify, err := time.Parse("10/10/2023 20:00", dateTime)
	if err != nil {
		log.Println("error while trying to parse datetime provided at request body")
		return false
	}
	dateTimeEnd := dateTimeToVerify.Add(time.Hour * 1)
	appointmentsByDateTime, err := r.store.GetAllAppointmentsByDateTimeInterval(dateTimeToVerify.Local().Add(time.Hour*3).String(), dateTimeEnd.Local().Add(time.Hour*3).String())
	if err != nil {
		log.Println("an error occurred while trying to get appointments with same date to validation")
		return false
	}

	if len(appointmentsByDateTime) == 0 {
		return true
	}
	for _, appointment := range appointmentsByDateTime {
		appointmentDateTime, err := time.Parse("2006-01-02 15:04:05", appointment.DateAndTime)
		if err != nil {
			return false
		}

		if appointmentDateTime.Local().Add(time.Hour*3).String() == dateTimeToVerify.Local().Add(time.Hour*3).String() {
			if d != appointment.DentistLicense && p != appointment.PatientIdentity {
				return true
			}
			if d != appointment.DentistLicense {
				if p == appointment.PatientIdentity {
					return true
				}
			}
			if d == appointment.DentistLicense {
				if p != appointment.PatientIdentity {
					return true
				}
			}
		}
		return true
	}
	return false
}
