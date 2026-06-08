-- +goose Up
INSERT INTO users (uuid, login, password_hash, created_at)
VALUES (
    gen_random_uuid(),
    'testuser',
    '$2a$10$IXA/JIGTPY38UZN/rK0LJuBtYe9C8.B.F7.mNiKFjouSA3ynzhPiG',
    NOW()
) ON CONFLICT (login) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE login = 'testuser';
