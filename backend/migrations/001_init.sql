-- SmartQueue: начальная схема

CREATE TABLE IF NOT EXISTS queues (
    id                   TEXT PRIMARY KEY,
    name                 TEXT        NOT NULL,
    description          TEXT        NOT NULL DEFAULT '',
    start_time           TEXT        NOT NULL DEFAULT '',
    max_participants     INT         NOT NULL DEFAULT 0,
    has_priority         BOOLEAN     NOT NULL DEFAULT FALSE,
    priority_count       INT         NOT NULL DEFAULT 0,
    initial_priority     INT         NOT NULL DEFAULT 0,
    anonymous_chat       BOOLEAN     NOT NULL DEFAULT FALSE,
    system_notifications BOOLEAN     NOT NULL DEFAULT FALSE,
    swap_positions       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_number       INT         NOT NULL DEFAULT 0,
    finished             BOOLEAN     NOT NULL DEFAULT FALSE
);

-- Администраторы очереди (один к многим)
CREATE TABLE IF NOT EXISTS queue_admins (
    queue_id TEXT NOT NULL REFERENCES queues(id) ON DELETE CASCADE,
    admin    TEXT NOT NULL,
    PRIMARY KEY (queue_id, admin)
);

-- Участники очереди
CREATE TABLE IF NOT EXISTS participants (
    id        TEXT        PRIMARY KEY,
    queue_id  TEXT        NOT NULL REFERENCES queues(id) ON DELETE CASCADE,
    name      TEXT        NOT NULL,
    priority  INT         NOT NULL DEFAULT 0,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индекс для быстрой выборки участников очереди по порядку
CREATE INDEX IF NOT EXISTS idx_participants_queue_joined
    ON participants(queue_id, joined_at ASC);

ALTER TABLE queues ADD COLUMN IF NOT EXISTS current_participant_id TEXT;
ALTER TABLE queues ADD COLUMN IF NOT EXISTS current_participant_name TEXT;