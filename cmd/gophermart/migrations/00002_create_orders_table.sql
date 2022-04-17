-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL,
    num TEXT,
    login TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status INTEGER NOT NULL DEFAULT 0,
    unique (num)
);

-- +goose Down
DROP TABLE IF EXISTS orders;