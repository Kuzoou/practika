package db

import (
	"database/sql"
	"log/slog"
	"os"

	_ "modernc.org/sqlite" // драйвер SQLite
)

// Переменная, в которой будет храниться подключение к базе
var DB *sql.DB

// Connect открывает соединение с базой
func Connect() error {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./store.db"
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Проверяем, что соединение реально работает
	err = DB.Ping()
	if err != nil {
		return err
	}

	slog.Info("успешно подключились к базе данных", "path", dbPath)
	return nil
}

// Close закрывает соединение (на всякий случай)
func Close() {
	if DB != nil {
		DB.Close()
	}
}