package db

import (
	"log/slog"
	"os"
)

// Migrate выполняет все миграции
func Migrate() error {
	files := []string{
		"migrations/001_init.sql",
		"migrations/002_roles.sql",
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return err
		}

		slog.Info("выполнена миграция", "file", file)
	}

	slog.Info("все миграции успешно выполнены")
	return nil
}
