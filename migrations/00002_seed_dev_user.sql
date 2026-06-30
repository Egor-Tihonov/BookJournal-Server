-- +goose Up
-- Dev-only: seed a user so the stubbed currentUserID()=1 has a real row to
-- satisfy the user_shelf.user_id foreign key. Remove once real auth lands.
INSERT INTO users (id, username, email, password_hash)
VALUES (1, 'dev', 'dev@bookjournal.local', 'dev-not-a-real-hash')
ON CONFLICT (id) DO NOTHING;

-- Keep the id sequence ahead of the explicit id=1 so app-created users don't collide.
SELECT setval(pg_get_serial_sequence('users', 'id'), GREATEST((SELECT max(id) FROM users), 1));

-- +goose Down
DELETE FROM users WHERE id = 1;
