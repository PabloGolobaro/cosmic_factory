-- +goose Up
INSERT INTO users (uuid, login, password_hash, role, created_at)
VALUES (
    gen_random_uuid(),
    'testmanager',
    '$2a$10$IXA/JIGTPY38UZN/rK0LJuBtYe9C8.B.F7.mNiKFjouSA3ynzhPiG',
    'manager',
    NOW()
) ON CONFLICT (login) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE login = 'testmanager';
