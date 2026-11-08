package main

import (
	"log/slog"
	"net/http"
	"os"

	"clothing-store/internal/db"
	"clothing-store/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		slog.Warn("не удалось загрузить .env, используем значения по умолчанию")
	}

	// Подключаемся к базе
	err = db.Connect()
	if err != nil {
		slog.Error("не удалось подключиться к базе", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Выполняем миграции
	err = db.Migrate()
	if err != nil {
		slog.Error("ошибка при выполнении миграций", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Тестовый эндпоинт
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Роли
	r.Post("/role", handlers.CreateRole)
	r.Get("/role", handlers.GetRoles)
	r.Put("/role/{id}", handlers.UpdateRole)
	r.Delete("/role/{id}", handlers.DeleteRole)

	// Товары
	r.Post("/product", handlers.CreateProduct)
	r.Get("/product", handlers.GetProducts)
	r.Get("/product/{id}", handlers.GetProduct)
	r.Put("/product/{id}", handlers.UpdateProduct)
	r.Delete("/product/{id}", handlers.DeleteProduct)

	// Пользователи
	r.Post("/user", handlers.CreateUser)
	r.Get("/user", handlers.GetUsers)
	r.Get("/user/role/{id}", handlers.GetUsersByRole)
	r.Put("/user/{id}", handlers.UpdateUser)
	r.Delete("/user/{id}", handlers.DeleteUser)

	// Корзина
	r.Post("/cart", handlers.CreateCart)
	r.Get("/cart/{user_id}", handlers.GetCart)
	r.Put("/cart/{id}", handlers.UpdateCart)
	r.Delete("/cart/{id}", handlers.DeleteCart)

	// Заказы
	r.Post("/order", handlers.CreateOrder)
	r.Get("/order", handlers.GetOrders)
	r.Get("/order/{id}", handlers.GetOrder)
	r.Put("/order/{id}", handlers.UpdateOrder)

	// Элементы заказа
	r.Post("/order_item", handlers.CreateOrderItem)
	r.Get("/order_item/{order_id}", handlers.GetOrderItems)

	// Фотографии
	r.Post("/photo", handlers.CreatePhoto)
	r.Get("/photo/{id}", handlers.GetPhoto)
	r.Delete("/photo/{id}", handlers.DeletePhoto)

	r.Post("/order_item", handlers.CreateOrderItem)
	r.Get("/order_item/{id}", handlers.GetOrderItems)

	slog.Info("сервер запускается", "port", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		slog.Error("ошибка запуска сервера", "error", err)
		os.Exit(1)
	}
}
