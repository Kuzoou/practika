package models

// Cart — элемент корзины
type Cart struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Count     int    `json:"count"`
}
