package domain

type Appointment struct {
	Id              int    `json:"id"`
	Description     string `json:"description" binding:"required"`
	DateAndTime     string `json:"date_and_time" binding:"required"`
	DentistLicense  string `json:"dentist_license" binding:"required"`
	PatientIdentity string `json:"patient_identity" binding:"required"`
}
