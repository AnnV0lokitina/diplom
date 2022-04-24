-- +goose Up
CREATE TABLE IF NOT EXISTS balance (
    id SERIAL,
    created_at timestamp (0) with time zone DEFAULT NOW(),
    operation_type integer,
    delta integer,
    order_id integer
);

-- +goose Down
DROP TABLE IF EXISTS balance;