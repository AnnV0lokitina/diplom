-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER NOT NULL PRIMARY KEY,
    num TEXT,
    login TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status INTEGER NOT NULL DEFAULT 0,
    unique (num)
);

-- +goose Down
DROP TABLE IF EXISTS orders;