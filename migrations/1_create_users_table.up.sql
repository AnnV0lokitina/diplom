CREATE TABLE IF NOT EXISTS users (
    login TEXT not null,
    password TEXT not null,
    active_session_id TEXT not null,
    unique (login)
);