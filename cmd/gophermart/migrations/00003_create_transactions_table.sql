-- +goose Up
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL,
    created_at timestamp (0) with time zone DEFAULT NOW(),
    operation_type integer,
    delta integer,
    order_id INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS balance;