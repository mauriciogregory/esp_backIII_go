package domain

type Patient struct {
	Id             int    `json:"id"`
	Surname        string `json:"surname" binding:"required"`
	Name           string `json:"name" binding:"required"`
	IdentityNumber string `json:"identity_number" binding:"required"`
	CreatedAt      string `json:"created_at" binding:"required"`
}
