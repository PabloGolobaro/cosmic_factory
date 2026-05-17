-- +goose Up
UPDATE orders SET user_uuid = gen_random_uuid() WHERE user_uuid IS NULL;

-- +goose Down
UPDATE orders SET user_uuid = NULL;
