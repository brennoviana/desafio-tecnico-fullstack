-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    cpf TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES topics(id),
    open_at BIGINT NOT NULL,
    close_at BIGINT NOT NULL
);

CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES topics(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    choice TEXT NOT NULL,
    UNIQUE(topic_id, user_id)
); 

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS users; 