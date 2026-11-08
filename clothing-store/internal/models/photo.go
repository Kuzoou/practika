package models

// Photo — фотография товара
type Photo struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	ProductID string `json:"product_id"`
}
