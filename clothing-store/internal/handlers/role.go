package handlers

import (
	"encoding/json"
	"net/http"

	"clothing-store/internal/db"
	"clothing-store/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateRole — создать новую роль
func CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role

	// Читаем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	// Генерируем ID
	role.ID = uuid.New().String()

	// Сохраняем в базу
	_, err = db.DB.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", role.ID, role.Name)
	if err != nil {
		http.Error(w, "ошибка при создании роли", http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную роль
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// GetRoles — получить все роли
func GetRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name FROM roles")
	if err != nil {
		http.Error(w, "ошибка при получении ролей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}
		roles = append(roles, role)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// UpdateRole — обновить роль
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var role models.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("UPDATE roles SET name = ? WHERE id = ?", role.Name, id)
	if err != nil {
		http.Error(w, "ошибка при обновлении роли", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "роль не найдена", http.StatusNotFound)
		return
	}

	role.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// DeleteRole — удалить роль
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.DB.Exec("DELETE FROM roles WHERE id = ?", id)
	if err != nil {
		http.Error(w, "ошибка при удалении роли", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "роль не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
