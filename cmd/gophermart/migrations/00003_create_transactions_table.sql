-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tr_type') THEN CREATE TYPE tr_type AS ENUM ('ADD', 'SUB');
    END IF;
END$$;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    created_at timestamp (0) with time zone DEFAULT NOW(),
    operation_type tr_type,
    delta integer,
    order_id INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TYPE IF EXISTS tr_type;
-- +goose StatementEnd