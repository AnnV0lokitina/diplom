CREATE TABLE IF NOT EXISTS balance (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    operation_type integer,
    balance_changed integer,
    order_id integer,
    reason integer
);