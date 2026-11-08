package models

// Order — заказ пользователя
type Order struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	DateTime string `json:"datetime"`
	Status   string `json:"status"`
}
