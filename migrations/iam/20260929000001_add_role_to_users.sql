-- +goose Up
ALTER TABLE users
    ADD COLUMN role VARCHAR(32) NOT NULL DEFAULT 'client'
        CONSTRAINT users_role_check CHECK (role IN ('client', 'manager'));

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS role;
