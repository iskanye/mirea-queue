CREATE TABLE IF NOT EXISTS users
(
    chat_id BIGINT PRIMARY KEY,
    username TEXT NOT NULL,
    group TEXT NOT NULL
);