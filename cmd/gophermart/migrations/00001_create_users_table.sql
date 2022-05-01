-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT not null,
    password TEXT not null,
    unique (login)
);

-- +goose Down
DROP TABLE IF EXISTS users;