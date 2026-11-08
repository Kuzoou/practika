package handlers

import (
	"clothing-store/internal/db"
	"clothing-store/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Допустимые статусы заказов
var orderStatuses = []string{
	"Принят",
	"В работе",
	"Собран",
	"Готов к выдаче",
	"Отменён",
}

// Проверка допустимого статуса
func isValidOrderStatus(status string) bool {
	for _, s := range orderStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// CreateOrder — создать заказ
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	if order.UserID == "" {
		http.Error(w, "не указан пользователь", http.StatusBadRequest)
		return
	}

	order.ID = uuid.New().String()
	order.DateTime = time.Now().Format(time.RFC3339)
	order.Status = "Принят"

	_, err = db.DB.Exec(`
		INSERT INTO orders (id, user_id, datetime, status)
		VALUES (?, ?, ?, ?)`,
		order.ID, order.UserID, order.DateTime, order.Status,
	)

	if err != nil {
		http.Error(w, "ошибка при создании заказа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrders — получить список заказов
func GetOrders(w http.ResponseWriter, r *http.Request) {
	// Пагинация
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Фильтры
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	status := r.URL.Query().Get("status")

	query := `
		SELECT id, user_id, datetime, status
		FROM orders
		WHERE 1=1`

	var args []interface{}

	// Фильтрация по начальной дате
	if from != "" {
		query += " AND datetime >= ?"
		args = append(args, from)
	}

	// Фильтрация по конечной дате
	if to != "" {
		query += " AND datetime <= ?"
		args = append(args, to)
	}

	// Фильтрация по статусу
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY datetime DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "ошибка при получении заказов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.DateTime,
			&order.Status,
		)

		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "ошибка при чтении заказов", http.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrder — получить конкретный заказ
func GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var order models.Order

	err := db.DB.QueryRow(`
		SELECT id, user_id, datetime, status
		FROM orders WHERE id = ?`, id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.DateTime,
		&order.Status,
	)

	if err != nil {
		http.Error(w, "заказ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// UpdateOrder — обновить статус заказа
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	if !isValidOrderStatus(order.Status) {
		http.Error(w, "неверный статус заказа", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
		UPDATE orders SET status = ? WHERE id = ?`,
		order.Status, id,
	)

	if err != nil {
		http.Error(w, "ошибка при обновлении заказа", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "заказ не найден", http.StatusNotFound)
		return
	}

	order.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
