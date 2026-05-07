package main

import (
	"database/sql"
	"log"
	"sync"
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
	
	var cpID, cpName sql.NullString
	s.db.QueryRow(`
		SELECT current_participant_id, current_participant_name 
		FROM queues WHERE id = $1`, id,
	).Scan(&cpID, &cpName)
	if cpID.Valid && cpID.String != "" {
		queue.CurrentParticipant = &Participant{ID: cpID.String, Name: cpName.String}
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

func (s *Store) RemoveParticipant(queueID, participantID string) {
    s.db.Exec(`DELETE FROM participants WHERE id = $1 AND queue_id = $2`,
        participantID, queueID)
}

// ShiftParticipant удаляет первого участника (FIFO) и увеличивает счётчик.
func (s *Store) ShiftParticipant(queueID string) {
	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	// 1. Берём первого участника + время входа
	var pID, pName string
	var joinedAt time.Time

	err = tx.QueryRow(`
		SELECT id, name, joined_at
		FROM participants
		WHERE queue_id = $1
		ORDER BY joined_at ASC
		LIMIT 1
	`, queueID).Scan(&pID, &pName, &joinedAt)

	if err != nil {
		return
	}

	// 2. Считаем РЕАЛЬНОЕ время ожидания (в минутах)
	now := time.Now()
	waitMinutes := int(now.Sub(joinedAt).Minutes())

	// 3. Считаем сколько людей было перед ним
	var count int
	tx.QueryRow(`
		SELECT COUNT(*) FROM participants WHERE queue_id = $1
	`, queueID).Scan(&count)

	// 4. Удаляем из очереди
	tx.Exec(`
		DELETE FROM participants WHERE id = $1
	`, pID)

	// 5. Обновляем очередь
	tx.Exec(`
		UPDATE queues
		SET current_number = current_number + 1,
			current_participant_id = $2,
			current_participant_name = $3,
			last_served_at = NOW()
		WHERE id = $1
	`, queueID, pID, pName)

	// 6. Сохраняем реальное наблюдение
	tx.Exec(`
		INSERT INTO queue_wait_stats (queue_id, wait_time, participants_count)
		VALUES ($1, $2, $3)
	`, queueID, waitMinutes, count)

	tx.Commit()
}

// FinishQueue помечает очередь как завершённую.
func (s *Store) FinishQueue(queueID string) {
	if _, err := s.db.Exec(`UPDATE queues SET finished = TRUE WHERE id = $1`, queueID); err != nil {
		log.Printf("Ошибка FinishQueue для %s: %v", queueID, err)
	}
}


// SwapOffer — ожидающий запрос на обмен
type SwapOffer struct {
    ID       string
    QueueID  string
    FromID   string
    FromName string
    ToID     string
}

// SwapStore — хранит активные предложения обмена в памяти
type SwapStore struct {
    mu     sync.Mutex
    offers map[string]*SwapOffer // swapID → offer
}

func NewSwapStore() *SwapStore {
    return &SwapStore{offers: make(map[string]*SwapOffer)}
}

func (ss *SwapStore) Add(offer *SwapOffer) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    ss.offers[offer.ID] = offer
}

func (ss *SwapStore) Take(swapID string) (*SwapOffer, bool) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    o, ok := ss.offers[swapID]
    if ok { delete(ss.offers, swapID) }
    return o, ok
}

func (s *Store) SwapParticipants(queueID, idA, idB string) error {
    tx, err := s.db.Begin()
    if err != nil { return err }
    defer tx.Rollback()

    // Меняем joined_at местами — порядок определяется им
    var timeA, timeB time.Time
    tx.QueryRow(`SELECT joined_at FROM participants WHERE id = $1`, idA).Scan(&timeA)
    tx.QueryRow(`SELECT joined_at FROM participants WHERE id = $1`, idB).Scan(&timeB)

    tx.Exec(`UPDATE participants SET joined_at = $1 WHERE id = $2`, timeB, idA)
    tx.Exec(`UPDATE participants SET joined_at = $1 WHERE id = $2`, timeA, idB)

    return tx.Commit()
}