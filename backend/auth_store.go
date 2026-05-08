package main

import (
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already registered")

// Сохраняет нового пользователя. Возвращает ErrEmailTaken, если email уже занят
func (s *Store) CreateUser(user *User) error {
	_, err := s.db.Exec(`
		INSERT INTO users (id, email, password_hash, display_name)
		VALUES ($1, $2, $3, $4)`,
		user.ID, user.Email, user.PasswordHash, user.DisplayName,
	)
	if err != nil {
		// pgx возвращает ошибку с кодом 23505 при нарушении UNIQUE
		if isUniqueViolation(err) {
			return ErrEmailTaken
		}
		return err
	}
	return nil
}

// Возвращает пользователя по email (Возвращает ErrUserNotFound, если не найден)
func (s *Store) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, display_name
		FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Проверяет, является ли ошибка нарушением уникального ограничения постгрес
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// pgx кодирует sqlstate в строке ошибки
	return contains(err.Error(), "23505") || contains(err.Error(), "unique")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
