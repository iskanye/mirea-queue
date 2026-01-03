CREATE TABLE IF NOT EXISTS student_groups
(
    group_id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users
(
    chat_id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    group_id INTEGER REFERENCES student_groups (group_id),
    queue_access BOOLEAN NOT NULL
);
