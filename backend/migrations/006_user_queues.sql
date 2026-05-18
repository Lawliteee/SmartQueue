CREATE TABLE IF NOT EXISTS user_participant_map (
    user_id        TEXT NOT NULL,
    queue_id       TEXT NOT NULL REFERENCES queues(id) ON DELETE CASCADE,
    participant_id TEXT NOT NULL,
    PRIMARY KEY (user_id, queue_id)
);