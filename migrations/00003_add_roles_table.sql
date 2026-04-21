-- +goose Up
CREATE TABLE roles (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL,
    name        TEXT NOT NULL,
    CONSTRAINT fk_roles_users 
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_roles_user_id ON roles(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_roles_user_id;
DROP TABLE IF EXISTS roles;