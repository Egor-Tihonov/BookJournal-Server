-- +goose Up
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE books (
    id         BIGSERIAL PRIMARY KEY,
    isbn       TEXT UNIQUE,
    title      TEXT NOT NULL,
    author     TEXT,
    year       INT,
    genre      TEXT,
    source     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_shelf (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    book_id     BIGINT NOT NULL REFERENCES books (id),
    status      TEXT NOT NULL DEFAULT 'want_to_read'
                CHECK (status IN ('want_to_read', 'reading', 'read')),
    why_reading TEXT,
    started_at  TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, book_id)
);

CREATE TABLE quotes (
    id         BIGSERIAL PRIMARY KEY,
    shelf_id   BIGINT NOT NULL REFERENCES user_shelf (id) ON DELETE CASCADE,
    text       TEXT NOT NULL,
    page       INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE notes (
    id         BIGSERIAL PRIMARY KEY,
    shelf_id   BIGINT NOT NULL REFERENCES user_shelf (id) ON DELETE CASCADE,
    text       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE notes;
DROP TABLE quotes;
DROP TABLE user_shelf;
DROP TABLE books;
DROP TABLE users;
