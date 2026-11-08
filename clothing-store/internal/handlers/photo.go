package handlers

import (
	"clothing-store/internal/db"
	"clothing-store/internal/models"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const photoDir = "uploads/photos"
const maxPhotoSize = 5 << 20 // 5 МБ

// CreatePhoto — загрузить фотографию на сервер
func CreatePhoto(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер загружаемого файла
	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize)

	err := r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		http.Error(w, "ошибка при загрузке файла или превышен размер 5 МБ", http.StatusBadRequest)
		return
	}

	defer r.MultipartForm.RemoveAll()

	// Получаем файл из запроса
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "фотография не найдена", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем ID товара
	productID := r.FormValue("product_id")
	if productID == "" {
		http.Error(w, "не указан товар", http.StatusBadRequest)
		return
	}

	// Читаем содержимое файла
	data, err := io.ReadAll(io.LimitReader(file, maxPhotoSize+1))
	if err != nil {
		http.Error(w, "ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 || len(data) > maxPhotoSize {
		http.Error(w, "файл пустой или превышает размер 5 МБ", http.StatusBadRequest)
		return
	}

	// Проверяем формат изображения
	sample := data
	if len(sample) > 512 {
		sample = sample[:512]
	}

	contentType := http.DetectContentType(sample)

	var extension string

	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/webp":
		extension = ".webp"
	case "image/gif":
		extension = ".gif"
	default:
		http.Error(w, "поддерживаются только JPG, PNG, WEBP и GIF", http.StatusBadRequest)
		return
	}

	// Создаём папку для фотографий
	err = os.MkdirAll(photoDir, 0755)
	if err != nil {
		http.Error(w, "ошибка при создании папки для фотографий", http.StatusInternalServerError)
		return
	}

	// Создаём модель фотографии
	photo := models.Photo{
		ID:        uuid.New().String(),
		Filename:  uuid.New().String() + extension,
		ProductID: productID,
	}

	// Сохраняем файл на сервере
	filePath := filepath.Join(photoDir, photo.Filename)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		http.Error(w, "ошибка при сохранении фотографии", http.StatusInternalServerError)
		return
	}

	// Сохраняем информацию о фотографии в базе данных
	_, err = db.DB.Exec(`
		INSERT INTO photos (id, filename, product_id)
		VALUES (?, ?, ?)`,
		photo.ID, photo.Filename, photo.ProductID,
	)

	if err != nil {
		os.Remove(filePath)
		http.Error(w, "ошибка при добавлении фотографии в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

// GetPhoto — получить фотографию по ID
func GetPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var photo models.Photo

	err := db.DB.QueryRow(`
		SELECT id, filename, product_id
		FROM photos WHERE id = ?`, id,
	).Scan(&photo.ID, &photo.Filename, &photo.ProductID)

	if err == sql.ErrNoRows {
		http.Error(w, "фотография не найдена", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "ошибка при получении фотографии", http.StatusInternalServerError)
		return
	}

	// Получаем путь к файлу
	filePath := filepath.Join(photoDir, filepath.Base(photo.Filename))

	// Отправляем фотографию клиенту
	http.ServeFile(w, r, filePath)
}

// DeletePhoto — удалить фотографию по ID
func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var filename string

	err := db.DB.QueryRow(`
		SELECT filename FROM photos WHERE id = ?`, id,
	).Scan(&filename)

	if err == sql.ErrNoRows {
		http.Error(w, "фотография не найдена", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "ошибка при получении фотографии", http.StatusInternalServerError)
		return
	}

	// Удаляем запись из базы данных
	_, err = db.DB.Exec("DELETE FROM photos WHERE id = ?", id)
	if err != nil {
		http.Error(w, "ошибка при удалении фотографии из базы данных", http.StatusInternalServerError)
		return
	}

	// Удаляем сам файл с сервера
	filePath := filepath.Join(photoDir, filepath.Base(filename))

	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "запись удалена, но не удалось удалить файл", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
