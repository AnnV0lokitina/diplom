-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    login TEXT not null,
    password TEXT not null,
    active_session_id TEXT not null,
    unique (login)
);

-- +goose Down
DROP TABLE IF EXISTS users;