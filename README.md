# BookJournal-Server

v1.0.0
Добавление книги в библиотеку:

- Поиск книги по бд(ISBN/Название)
- Скачивание Название Автора Обложку год с https://openlibrary.org/api/books?bibkeys=ISBN:9781250319180&format=json&jscmd=data

v1.0.1
Добавление книги в библиотеку:

- Ручное добавление книги в бд(если ее не находит нигде)

v1.0.2
Добавление книги в библиотеку:

- Поиск книги по фотографии

## Схема БД

### users

| поле          | тип                      | заметки                       |
| ------------- | ------------------------ | ----------------------------- |
| id            | bigserial PK             |                               |
| username      | text, unique, not null   |                               |
| email         | text, unique, not null   |                               |
| password_hash | text, not null           | хэш пароля (bcrypt/argon2)    |
| created_at    | timestamptz, default now |                               |

### books

Глобальный каталог всех известных книг (общий для всех пользователей).

| поле       | тип                     | заметки                                    |
| ---------- | ----------------------- | ------------------------------------------ |
| id         | bigserial PK            |                                            |
| isbn       | text, unique, nullable  | у книг из openlibrary ISBN может не быть   |
| title      | text, not null          |                                            |
| author     | text                    |                                            |
| year       | int                     |                                            |
| genre      | text                    |                                            |
| source     | text                    | `manual` / `openlibrary` — откуда пришла   |
| created_at | timestamptz             |                                            |

### user_library

Полка/журнал пользователя: связь пользователь↔книга + статус и контекст чтения.

| поле        | тип                 | заметки                                            |
| ----------- | ------------------- | -------------------------------------------------- |
| id          | bigserial PK        |                                                    |
| user_id     | bigint FK→users(id) | on delete cascade                                  |
| book_id     | bigint FK→books(id) |                                                    |
| status      | text                | `want_to_read` / `reading` / `read`                |
| why_reading | text                | почему читаю                                       |
| started_at  | timestamptz, null   | заполняется при переходе в `reading`               |
| finished_at | timestamptz, null   | при переходе в `read`                              |
| created_at  | timestamptz         | когда добавил к себе                                |

UNIQUE(user_id, book_id) — нельзя добавить одну книгу дважды.

### quotes

Цитаты пользователя (много на одну книгу в полке).

| поле       | тип                          | заметки           |
| ---------- | ---------------------------- | ----------------- |
| id         | bigserial PK                 |                   |
| library_id | bigint FK→user_library(id)   | on delete cascade |
| text       | text, not null               |                   |
| page       | int, nullable                |                   |
| created_at | timestamptz                  |                   |

### notes

Мысли пользователя во времени (дневник чтения).

| поле       | тип                          | заметки           |
| ---------- | ---------------------------- | ----------------- |
| id         | bigserial PK                 |                   |
| library_id | bigint FK→user_library(id)   | on delete cascade |
| text       | text, not null               |                   |
| created_at | timestamptz                  |                   |
