package models

// User — пользователь
type User struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	RoleID           string `json:"role_id"`
	RegistrationDate string `json:"registration_date"`
}
