package handlers

import (
	"encoding/json"
	"net/http"

	"clothing-store/internal/db"
	"clothing-store/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateCart — добавить товар в корзину
func CreateCart(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart

	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	cart.ID = uuid.New().String()
	if cart.Count < 1 {
		cart.Count = 1
	}

	_, err = db.DB.Exec(`
		INSERT INTO cart (id, user_id, product_id, count)
		VALUES (?, ?, ?, ?)`,
		cart.ID, cart.UserID, cart.ProductID, cart.Count,
	)
	if err != nil {
		http.Error(w, "ошибка при добавлении в корзину", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// GetCart — получить корзину пользователя
func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	rows, err := db.DB.Query(`
		SELECT id, user_id, product_id, count FROM cart WHERE user_id = ?`, userID)
	if err != nil {
		http.Error(w, "ошибка при получении корзины", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Cart
	for rows.Next() {
		var item models.Cart
		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Count)
		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateCart — обновить количество товара в корзине
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cart models.Cart
	err := json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`UPDATE cart SET count = ? WHERE id = ?`, cart.Count, id)
	if err != nil {
		http.Error(w, "ошибка при обновлении корзины", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "элемент корзины не найден", http.StatusNotFound)
		return
	}

	cart.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// DeleteCart — удалить товар из корзины
func DeleteCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.DB.Exec("DELETE FROM cart WHERE id = ?", id)
	if err != nil {
		http.Error(w, "ошибка при удалении из корзины", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "элемент корзины не найден", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
