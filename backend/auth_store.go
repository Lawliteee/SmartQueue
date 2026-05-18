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

// QueueRef — краткая информация об очереди для списков
type QueueRef struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// Сохраняет связь пользователь → участник очереди
func (s *Store) SaveUserParticipant(userID, queueID, participantID string) error {
    _, err := s.db.Exec(`
        INSERT INTO user_participant_map (user_id, queue_id, participant_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, queue_id) DO UPDATE SET participant_id = EXCLUDED.participant_id`,
        userID, queueID, participantID,
    )
    return err
}

// Возвращает активные очереди пользователя: как админа и как участника
func (s *Store) GetUserQueues(userID string) (adminQueues []QueueRef, participantQueues []QueueRef, err error) {
    adminQueues = []QueueRef{}
    participantQueues = []QueueRef{}

    // Очереди где пользователь — админ (незавершённые)
    rows, err := s.db.Query(`
        SELECT q.id, q.name
        FROM queues q
        JOIN queue_admins qa ON qa.queue_id = q.id
        WHERE qa.admin = $1 AND q.finished = FALSE
        ORDER BY q.created_at DESC`, userID,
    )
    if err != nil {
        return
    }
    defer rows.Close()
    for rows.Next() {
        var ref QueueRef
        if rows.Scan(&ref.ID, &ref.Name) == nil {
            adminQueues = append(adminQueues, ref)
        }
    }

    // Очереди где пользователь — участник (незавершённые)
    pRows, err := s.db.Query(`
        SELECT q.id, q.name
        FROM queues q
        JOIN user_participant_map upm ON upm.queue_id = q.id
        WHERE upm.user_id = $1 AND q.finished = FALSE
        ORDER BY q.created_at DESC`, userID,
    )
    if err != nil {
        return
    }
    defer pRows.Close()
    for pRows.Next() {
        var ref QueueRef
        if pRows.Scan(&ref.ID, &ref.Name) == nil {
            participantQueues = append(participantQueues, ref)
        }
    }

    return
}
