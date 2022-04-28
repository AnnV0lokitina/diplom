-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL,
    num TEXT,
    user_id INTEGER NOT NULL,
    uploaded_at timestamp (0) with time zone DEFAULT NOW(),
    status INTEGER NOT NULL DEFAULT 0,
    unique (num)
);

-- +goose Down
DROP TABLE IF EXISTS orders;