package main

import (
	"database/sql"
	"log"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Save сохраняет новую очередь вместе с администраторами.
func (s *Store) Save(queue *Queue) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO queues (
			id, name, description, start_time, max_participants,
			has_priority, priority_count, initial_priority,
			anonymous_chat, system_notifications, swap_positions,
			created_at, current_number, finished
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		queue.ID, queue.Name, queue.Description, queue.StartTime, queue.MaxParticipants,
		queue.HasPriority, queue.PriorityCount, queue.InitialPriority,
		queue.AnonymousChat, queue.SystemNotifications, queue.SwapPositions,
		queue.CreatedAt, queue.CurrentNumber, queue.Finished,
	)
	if err != nil {
		return err
	}

	for _, admin := range queue.Admins {
		if _, err = tx.Exec(
			`INSERT INTO queue_admins (queue_id, admin) VALUES ($1, $2)`,
			queue.ID, admin,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Get возвращает очередь со всеми участниками и администраторами.
func (s *Store) Get(id string) (*Queue, bool) {
	queue := &Queue{}

	err := s.db.QueryRow(`
		SELECT id, name, description, start_time, max_participants,
		       has_priority, priority_count, initial_priority,
		       anonymous_chat, system_notifications, swap_positions,
		       created_at, current_number, finished
		FROM queues WHERE id = $1`, id,
	).Scan(
		&queue.ID, &queue.Name, &queue.Description, &queue.StartTime, &queue.MaxParticipants,
		&queue.HasPriority, &queue.PriorityCount, &queue.InitialPriority,
		&queue.AnonymousChat, &queue.SystemNotifications, &queue.SwapPositions,
		&queue.CreatedAt, &queue.CurrentNumber, &queue.Finished,
	)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		log.Printf("Ошибка получения очереди %s: %v", id, err)
		return nil, false
	}

	// Администраторы
	adminRows, err := s.db.Query(`SELECT admin FROM queue_admins WHERE queue_id = $1`, id)
	if err == nil {
		defer adminRows.Close()
		for adminRows.Next() {
			var admin string
			if adminRows.Scan(&admin) == nil {
				queue.Admins = append(queue.Admins, admin)
			}
		}
	}

	// Участники (отсортированы по времени вступления — FIFO)
	pRows, err := s.db.Query(`
		SELECT id, name, priority FROM participants
		WHERE queue_id = $1 ORDER BY joined_at ASC`, id,
	)
	if err == nil {
		defer pRows.Close()
		for pRows.Next() {
			var p Participant
			if pRows.Scan(&p.ID, &p.Name, &p.Priority) == nil {
				queue.Participants = append(queue.Participants, p)
			}
		}
	}

	if queue.Participants == nil {
		queue.Participants = []Participant{}
	}

	return queue, true
}

// AddParticipant добавляет участника в очередь.
func (s *Store) AddParticipant(queueID string, p Participant) {
	_, err := s.db.Exec(`
		INSERT INTO participants (id, queue_id, name, priority, joined_at)
		VALUES ($1, $2, $3, $4, $5)`,
		p.ID, queueID, p.Name, p.Priority, time.Now(),
	)
	if err != nil {
		log.Printf("Ошибка добавления участника в очередь %s: %v", queueID, err)
	}
}

// ShiftParticipant удаляет первого участника (FIFO) и увеличивает счётчик.
func (s *Store) ShiftParticipant(queueID string) {
	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM participants WHERE id = (
			SELECT id FROM participants
			WHERE queue_id = $1
			ORDER BY joined_at ASC LIMIT 1
		)`, queueID,
	)
	if err != nil {
		log.Printf("Ошибка ShiftParticipant для %s: %v", queueID, err)
		return
	}

	tx.Exec(`UPDATE queues SET current_number = current_number + 1 WHERE id = $1`, queueID)
	tx.Commit()
}

// FinishQueue помечает очередь как завершённую.
func (s *Store) FinishQueue(queueID string) {
	if _, err := s.db.Exec(`UPDATE queues SET finished = TRUE WHERE id = $1`, queueID); err != nil {
		log.Printf("Ошибка FinishQueue для %s: %v", queueID, err)
	}
}
