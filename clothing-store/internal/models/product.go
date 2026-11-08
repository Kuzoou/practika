package models

// Product — товар
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Gender      string  `json:"gender"`
	Size        string  `json:"size"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Season      string  `json:"season"`
	Price       float64 `json:"price"`
}
