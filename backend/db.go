package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Открывает соединение с PostgreSQL. Строка подключения берётся из переменной окружения DATABASE_URL.
func NewDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/smartqueue?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v\nПроверь DATABASE_URL или запусти docker-compose up -d", err)
	}

	log.Println("Подключение к БД установлено")
	return db
}

// Применяет все миграции по порядку
func RunMigrations(db *sql.DB) {
	migrations := []string{
		"migrations/001_init.sql",
		"migrations/002_users.sql",
		"migrations/003_rename_system_notifications.sql",
	}

	for _, path := range migrations {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Не удалось прочитать файл миграции %s: %v", path, err)
		}
		if _, err := db.Exec(string(data)); err != nil {
			log.Fatalf("Ошибка выполнения миграции %s: %v", path, err)
		}
		log.Printf("Миграция применена: %s", path)
	}
}
