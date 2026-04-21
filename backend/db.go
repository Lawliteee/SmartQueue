package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewDB открывает соединение с PostgreSQL.
// Строка подключения берётся из переменной окружения DATABASE_URL.
// Пример: postgres://user:password@localhost:5432/smartqueue?sslmode=disable
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

// RunMigrations применяет начальную схему из SQL-файла.
// При повторном запуске безопасно — использует IF NOT EXISTS.
func RunMigrations(db *sql.DB) {
	migration, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл миграции: %v", err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		log.Fatalf("Ошибка выполнения миграции: %v", err)
	}

	log.Println("Миграции применены")
}
