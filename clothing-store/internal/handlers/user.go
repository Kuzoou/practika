package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"clothing-store/internal/db"
	"clothing-store/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateUser — создать пользователя
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()
	user.RegistrationDate = time.Now().Format(time.RFC3339)

	_, err = db.DB.Exec(`
		INSERT INTO users (id, name, email, phone, role_id, registration_date)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.Phone, user.RoleID, user.RegistrationDate,
	)
	if err != nil {
		http.Error(w, "ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers — получить всех пользователей (с пагинацией)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := db.DB.Query(`
		SELECT id, name, email, phone, role_id, registration_date 
		FROM users LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		http.Error(w, "ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RegistrationDate)
		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUsersByRole — получить пользователей по id роли
func GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")

	rows, err := db.DB.Query(`
		SELECT id, name, email, phone, role_id, registration_date 
		FROM users WHERE role_id = ?`, roleID)
	if err != nil {
		http.Error(w, "ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RegistrationDate)
		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUser — обновить пользователя
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
		UPDATE users SET name=?, email=?, phone=?, role_id=? WHERE id=?`,
		user.Name, user.Email, user.Phone, user.RoleID, id,
	)
	if err != nil {
		http.Error(w, "ошибка при обновлении пользователя", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}

	user.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser — удалить пользователя
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
