package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clothing-store/internal/db"
	"clothing-store/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateProduct — создать товар
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	product.ID = uuid.New().String()

	_, err = db.DB.Exec(`
		INSERT INTO products (id, name, type, gender, size, color, description, brand, season, price)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID, product.Name, product.Type, product.Gender, product.Size,
		product.Color, product.Description, product.Brand, product.Season, product.Price,
	)
	if err != nil {
		http.Error(w, "ошибка при создании товара", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProducts — список товаров с фильтрами и пагинацией
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Параметры пагинации
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Фильтры
	productType := r.URL.Query().Get("type")
	gender := r.URL.Query().Get("gender")
	size := r.URL.Query().Get("size")
	color := r.URL.Query().Get("color")
	brand := r.URL.Query().Get("brand")
	season := r.URL.Query().Get("season")
	sort := r.URL.Query().Get("sort") // price_asc или price_desc

	query := "SELECT id, name, type, gender, size, color, description, brand, season, price FROM products WHERE 1=1"
	args := []interface{}{}

	if productType != "" {
		query += " AND type = ?"
		args = append(args, productType)
	}
	if gender != "" {
		query += " AND gender = ?"
		args = append(args, gender)
	}
	if size != "" {
		query += " AND size = ?"
		args = append(args, size)
	}
	if color != "" {
		query += " AND color = ?"
		args = append(args, color)
	}
	if brand != "" {
		query += " AND brand = ?"
		args = append(args, brand)
	}
	if season != "" {
		query += " AND season = ?"
		args = append(args, season)
	}

	// Сортировка
	if sort == "price_asc" {
		query += " ORDER BY price ASC"
	} else if sort == "price_desc" {
		query += " ORDER BY price DESC"
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "ошибка при получении товаров", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Size, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Price)
		if err != nil {
			http.Error(w, "ошибка при чтении данных", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetProduct — получить один товар
func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p models.Product
	err := db.DB.QueryRow(`
		SELECT id, name, type, gender, size, color, description, brand, season, price 
		FROM products WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Size, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Price)

	if err != nil {
		http.Error(w, "товар не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdateProduct — обновить товар
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`
		UPDATE products SET name=?, type=?, gender=?, size=?, color=?, description=?, brand=?, season=?, price=?
		WHERE id=?`,
		product.Name, product.Type, product.Gender, product.Size, product.Color,
		product.Description, product.Brand, product.Season, product.Price, id,
	)
	if err != nil {
		http.Error(w, "ошибка при обновлении товара", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "товар не найден", http.StatusNotFound)
		return
	}

	product.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct — удалить товар
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "ошибка при удалении товара", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "товар не найден", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
