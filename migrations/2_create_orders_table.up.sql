CREATE TABLE IF NOT EXISTS orders (
    order_number SERIAL PRIMARY KEY,
    login TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    current_status INTEGER
);