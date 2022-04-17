-- +goose Up
CREATE TABLE IF NOT EXISTS balance (
                                       id int NOT NULL PRIMARY KEY,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       operation_type integer,
                                       delta integer,
                                       order_id integer,
);

-- +goose Down
DROP TABLE IF EXISTS balance;