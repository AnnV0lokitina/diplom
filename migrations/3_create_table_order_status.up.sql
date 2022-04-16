CREATE TABLE IF NOT EXISTS order_change_status (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    order_number INTEGER not null,
    status INTEGER
);