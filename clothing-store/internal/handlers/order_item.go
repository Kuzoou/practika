package handlers

import (
	"clothing-store/internal/db"
	"clothing-store/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateOrderItem — добавить товар в заказ
func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	var item models.OrderItem

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	if item.OrderID == "" || item.ProductID == "" {
		http.Error(w, "не указан заказ или товар", http.StatusBadRequest)
		return
	}

	if item.Price < 0 {
		http.Error(w, "цена не может быть отрицательной", http.StatusBadRequest)
		return
	}

	item.ID = uuid.New().String()

	_, err = db.DB.Exec(`
		INSERT INTO order_items (id, order_id, product_id, price)
		VALUES (?, ?, ?, ?)`,
		item.ID, item.OrderID, item.ProductID, item.Price,
	)

	if err != nil {
		http.Error(w, "ошибка при добавлении товара в заказ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// GetOrderItems — получить все товары заказа
func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	rows, err := db.DB.Query(`
		SELECT id, order_id, product_id, price
		FROM order_items
		WHERE order_id = ?`, orderID)

	if err != nil {
		http.Error(w, "ошибка при получении элементов заказа", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		var item models.OrderItem

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Price,
		)

		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "ошибка при чтении элементов заказа", http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = []models.OrderItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
