-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL,
    session_id TEXT not null,
    created_at INTEGER NOT NULL,
    lifetime INTEGER NOT NULL,
    login TEXT not null
);


-- +goose Down
DROP TABLE IF EXISTS sessions;