-- +goose Up
CREATE TABLE sessions (
    id          VARCHAR(36) PRIMARY KEY,
    user_id     INTEGER NOT NULL, -- Изменено на INTEGER, так как users.id — это SERIAL (целое число)
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_sessions_users 
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE -- Рекомендуется: при удалении пользователя удалять и его сессии
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- +goose Down
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP TABLE IF EXISTS sessions;