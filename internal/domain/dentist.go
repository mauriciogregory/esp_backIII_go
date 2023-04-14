package domain

type Dentist struct {
	Id            int    `json:"id"`
	Surname       string `json:"surname" binding:"required"`
	Name          string `json:"name" binding:"required"`
	LicenseNumber string `json:"license_number" binding:"required"`
}
